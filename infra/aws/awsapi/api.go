package awsapi

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/iancoleman/strcase"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

type APIProps struct {
	DeploymentIdent *string
	HostedZone      awsroute53.IHostedZone             // optional: nil if no custom domain
	Certificate     awscertificatemanager.ICertificate // optional: nil if no custom domain
	MainTable       awsdynamodb.ITableV2
	MainTableName   *string
}

type Api interface {
	HttpApi() awsapigatewayv2.HttpApi
}

type api struct {
	logGroup awslogs.ILogGroup
	function awslambdago.GoFunction
	httpAPI  awsapigatewayv2.HttpApi
}

func NewApi(parent constructs.Construct, props APIProps) Api {
	scope, con := constructs.NewConstruct(parent, jsii.String("Api")), &api{}
	qual, stack := cdkutil.QualifierFromContext(scope), awscdk.Stack_Of(scope)

	con.logGroup = awslogs.NewLogGroup(scope, jsii.String("ApiLogGroup"), &awslogs.LogGroupProps{
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		Retention:     awslogs.RetentionDays_TWO_WEEKS,
	})

	adapterLayerArn := fmt.Sprintf("arn:aws:lambda:%s:753240598075:layer:LambdaAdapterLayerArm64:24",
		*stack.Region())

	con.function = awslambdago.NewGoFunction(scope, jsii.String("ApiFunction"),
		&awslambdago.GoFunctionProps{
			Entry:        jsii.String("../../../backend/lambda/httpapi"),
			LogGroup:     con.logGroup,
			Architecture: awslambda.Architecture_ARM_64(),
			Layers: &[]awslambda.ILayerVersion{
				awslambda.LayerVersion_FromLayerVersionArn(
					scope, jsii.String("LambdaAdapterLayer"), jsii.String(adapterLayerArn)),
			},
			Environment: &map[string]*string{
				"AWS_LAMBDA_EXEC_WRAPPER": jsii.String("/opt/bootstrap"),
				"AWS_LWA_PORT":            jsii.String("12001"),
				"MAIN_TABLE_NAME":         props.MainTableName,
			},
			Bundling: &awslambdago.BundlingOptions{},
		})

	lambdaIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("LambdaIntegration"),
		con.function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	con.httpAPI = awsapigatewayv2.NewHttpApi(scope, jsii.String("HttpApi"),
		&awsapigatewayv2.HttpApiProps{
			ApiName: jsii.Sprintf("%s-%s-httpapi",
				qual,
				strcase.ToKebab(*props.DeploymentIdent)),
			DefaultIntegration:         lambdaIntegration,
			DefaultAuthorizationScopes: jsii.Strings("main/admin"),
			CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
				AllowOrigins: jsii.Strings(
					"http://localhost:5173",
					"https://*",
				),
				AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{
					awsapigatewayv2.CorsHttpMethod_GET,
					awsapigatewayv2.CorsHttpMethod_POST,
					awsapigatewayv2.CorsHttpMethod_PUT,
					awsapigatewayv2.CorsHttpMethod_DELETE,
					awsapigatewayv2.CorsHttpMethod_OPTIONS,
				},
				AllowHeaders: jsii.Strings(
					"Content-Type",
					"Authorization",
					"Connect-Protocol-Version",
				),
				AllowCredentials: jsii.Bool(true),
				ExposeHeaders: jsii.Strings(
					"Grpc-Status",
					"Grpc-Message",
					"Grpc-Status-Details-Bin",
				),
				MaxAge: awscdk.Duration_Seconds(jsii.Number(7200)),
			},
		})

	con.httpAPI.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:                jsii.String("/{proxy+}"),
		Methods:             &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_OPTIONS},
		Integration:         lambdaIntegration,
		AuthorizationScopes: &[]*string{},
		Authorizer:          awsapigatewayv2.NewHttpNoneAuthorizer(),
	})

	props.MainTable.GrantReadWriteData(con.function)

	if props.HostedZone != nil && props.Certificate != nil {
		customDomainName := strcase.ToKebab(*props.DeploymentIdent) + "." + *props.HostedZone.ZoneName()

		domainName := awsapigatewayv2.NewDomainName(scope, jsii.String("DomainName"),
			&awsapigatewayv2.DomainNameProps{
				DomainName:  jsii.String(customDomainName),
				Certificate: props.Certificate,
			})

		awsapigatewayv2.NewApiMapping(scope, jsii.String("ApiMapping"),
			&awsapigatewayv2.ApiMappingProps{
				Api:        con.httpAPI,
				DomainName: domainName,
			})

		awsroute53.NewCfnRecordSet(scope, jsii.String("LatencyRecord"),
			&awsroute53.CfnRecordSetProps{
				Name:          jsii.String(customDomainName),
				Type:          jsii.String("A"),
				HostedZoneId:  props.HostedZone.HostedZoneId(),
				SetIdentifier: stack.Region(),
				Region:        stack.Region(),
				AliasTarget: &awsroute53.CfnRecordSet_AliasTargetProperty{
					DnsName:              domainName.RegionalDomainName(),
					HostedZoneId:         domainName.RegionalHostedZoneId(),
					EvaluateTargetHealth: jsii.Bool(true),
				},
			})

		awscdk.NewCfnOutput(scope, jsii.String("ApiEndpoint"), &awscdk.CfnOutputProps{
			Value:       jsii.String("https://" + customDomainName),
			Description: jsii.String("The API custom domain URL"),
		})
	} else {
		awscdk.NewCfnOutput(scope, jsii.String("ApiEndpoint"), &awscdk.CfnOutputProps{
			Value:       con.httpAPI.Url(),
			Description: jsii.String("The API Gateway URL"),
		})
	}

	return con
}

func (con api) HttpApi() awsapigatewayv2.HttpApi {
	return con.httpAPI
}
