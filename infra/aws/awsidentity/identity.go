package awsidentity

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/luckymaks/bm_backend/infra/aws/awsparams"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

const (
	identityNamespace = "identity"
	userPoolIDParam   = "user-pool-id"
)

type IdentityProps struct{}

type Identity interface {
	UserPool() awscognito.IUserPool
}

type identity struct {
	userPool awscognito.IUserPool
}

func New(scope constructs.Construct, _ IdentityProps) Identity {
	scope, con := constructs.NewConstruct(scope, jsii.String("Identity")), &identity{}
	qual := cdkutil.QualifierFromContext(scope)

	if cdkutil.IsPrimaryRegion(scope) {
		pool := awscognito.NewUserPool(scope, jsii.String("UserPool"), &awscognito.UserPoolProps{
			UserPoolName: jsii.Sprintf("%s-users", qual),
			SelfSignUpEnabled: jsii.Bool(false),
			SignInAliases: &awscognito.SignInAliases{
				Email: jsii.Bool(true),
			},
			AutoVerify: &awscognito.AutoVerifiedAttrs{
				Email: jsii.Bool(true),
			},
			StandardAttributes: &awscognito.StandardAttributes{
				Email: &awscognito.StandardAttribute{
					Required: jsii.Bool(true),
					Mutable:  jsii.Bool(true),
				},
			},
			PasswordPolicy: &awscognito.PasswordPolicy{
				MinLength:        jsii.Number(8),
				RequireLowercase: jsii.Bool(true),
				RequireUppercase: jsii.Bool(true),
				RequireDigits:    jsii.Bool(true),
				RequireSymbols:   jsii.Bool(true),
			},
			AccountRecovery: awscognito.AccountRecovery_EMAIL_ONLY,
			RemovalPolicy:   awscdk.RemovalPolicy_DESTROY,
		})

		pool.AddResourceServer(jsii.String("MainResourceServer"), &awscognito.UserPoolResourceServerOptions{
			Identifier: jsii.String("main"),
			Scopes: &[]awscognito.ResourceServerScope{
				awscognito.NewResourceServerScope(&awscognito.ResourceServerScopeProps{
					ScopeName:        jsii.String("admin"),
					ScopeDescription: jsii.String("Administrative access"),
				}),
			},
		})

		awsparams.Store(scope, "UserPoolIDParam", identityNamespace, userPoolIDParam, pool.UserPoolId())

		con.userPool = pool
	} else {
		userPoolID := awsparams.Lookup(scope, "LookupUserPoolID", identityNamespace, userPoolIDParam, "user-pool-id")
		con.userPool = awscognito.UserPool_FromUserPoolId(scope, jsii.String("UserPool"), userPoolID)
	}

	return con
}

func (i *identity) UserPool() awscognito.IUserPool {
	return i.userPool
}
