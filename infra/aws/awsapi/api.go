package awsapi

import (
	"fmt"
	
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ApiProps struct {
	DeploymentIdent *string
}

type Api interface {
	Function() awslambda.IFunction
}

type api struct {
	function awslambda.IFunction
}

func NewApi(scope constructs.Construct, props ApiProps) Api {
	construct := constructs.NewConstruct(scope, jsii.String("Api"))
	
	functionName := jsii.String(fmt.Sprintf("kndr-%s-httpapi", *props.DeploymentIdent))
	
	logGroup := awslogs.NewLogGroup(construct, jsii.String("LogGroup"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String(fmt.Sprintf("/aws/lambda/%s", *functionName)),
		Retention:     awslogs.RetentionDays_TWO_WEEKS,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	
	layerArn := fmt.Sprintf("arn:aws:lambda:%s:753240598075:layer:LambdaAdapterLayerArm64:24",
		*awscdk.Stack_Of(construct).Region())
	
	webAdapterLayer := awslambda.LayerVersion_FromLayerVersionArn(
		construct,
		jsii.String("WebAdapterLayer"),
		jsii.String(layerArn),
	)
	
	fn := awscdklambdagoalpha.NewGoFunction(construct, jsii.String("Function"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: functionName,
		Entry:        jsii.String("../../../backend/lambda/httpapi"),
		Architecture: awslambda.Architecture_ARM_64(),
		MemorySize:   jsii.Number(512),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
		LogGroup:     logGroup,
		Layers:       &[]awslambda.ILayerVersion{webAdapterLayer},
		Environment: &map[string]*string{
			"AWS_LAMBDA_EXEC_WRAPPER": jsii.String("/opt/bootstrap"),
			"AWS_LWA_PORT":            jsii.String("12001"),
		},
	})
	
	return &api{
		function: fn,
	}
}

func (a *api) Function() awslambda.IFunction {
	return a.function
}
