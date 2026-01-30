package awsidentity

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/luckymaks/bm_backend/infra/aws/awsparams"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

const paramsNamespace = "identity"

type IdentityProps struct{}

type Identity interface {
	UserPool() awscognito.IUserPool
	UserPoolID() *string
	MainResourceServer() awscognito.IUserPoolResourceServer
	CognitoDomain() *string
	AdminClientID() *string
}

type identity struct {
	userPool           awscognito.IUserPool
	userPoolID         *string
	mainResourceServer awscognito.IUserPoolResourceServer
	cognitoDomain      *string
	adminClientID      *string
}

func New(parent constructs.Construct, _ IdentityProps) Identity {
	scope, con := constructs.NewConstruct(parent, jsii.String("Identity")), &identity{}
	qual := cdkutil.QualifierFromContext(scope)

	if !cdkutil.IsPrimaryRegion(scope) {
		con.userPoolID = awsparams.Lookup(
			scope, "LookupUserPoolID", paramsNamespace, "user-pool-id", "user-pool-id-lookup")
		con.adminClientID = awsparams.Lookup(
			scope, "LookupAdminClientID", paramsNamespace, "admin-client-id", "admin-client-id-lookup")

		return con
	}

	con.userPool = awscognito.NewUserPool(scope, jsii.String("UserPool"), &awscognito.UserPoolProps{
		UserPoolName:      jsii.Sprintf("%s-users", qual),
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

	con.userPoolID = con.userPool.UserPoolId()

	userPoolDomain := con.userPool.AddDomain(jsii.String("CognitoDomain"), &awscognito.UserPoolDomainOptions{
		CognitoDomain: &awscognito.CognitoDomainOptions{
			DomainPrefix: jsii.Sprintf("%s-auth", qual),
		},
		ManagedLoginVersion: awscognito.ManagedLoginVersion_NEWER_MANAGED_LOGIN,
	})

	con.cognitoDomain = userPoolDomain.BaseUrl(nil)

	adminScope := awscognito.NewResourceServerScope(&awscognito.ResourceServerScopeProps{
		ScopeName:        jsii.String("admin"),
		ScopeDescription: jsii.String("Administrative access"),
	})

	con.mainResourceServer = con.userPool.AddResourceServer(jsii.String("MainResourceServer"),
		&awscognito.UserPoolResourceServerOptions{
			Identifier: jsii.String("main"),
			Scopes:     &[]awscognito.ResourceServerScope{adminScope},
		})

	adminClient := con.userPool.AddClient(jsii.String("AdminClient"), &awscognito.UserPoolClientOptions{
		UserPoolClientName: jsii.Sprintf("%s-admin", qual),
		GenerateSecret:     jsii.Bool(true),
		OAuth: &awscognito.OAuthSettings{
			Flows: &awscognito.OAuthFlows{
				ClientCredentials: jsii.Bool(true),
			},
			Scopes: &[]awscognito.OAuthScope{
				awscognito.OAuthScope_ResourceServer(con.mainResourceServer, adminScope),
			},
		},
		AuthFlows: &awscognito.AuthFlow{
			UserSrp:           jsii.Bool(false),
			UserPassword:      jsii.Bool(false),
			AdminUserPassword: jsii.Bool(false),
			Custom:            jsii.Bool(false),
		},
		AccessTokenValidity: awscdk.Duration_Hours(jsii.Number(24)),
	})

	con.adminClientID = adminClient.UserPoolClientId()

	replicaRegions := make([]*awssecretsmanager.ReplicaRegion, 0)
	for _, region := range cdkutil.SecondaryRegions(scope) {
		replicaRegions = append(replicaRegions, &awssecretsmanager.ReplicaRegion{
			Region: jsii.String(region),
		})
	}

	awssecretsmanager.NewSecret(scope, jsii.String("AdminClientSecret"), &awssecretsmanager.SecretProps{
		SecretName:        jsii.Sprintf("%s/admin-client-secret", qual),
		SecretStringValue: adminClient.UserPoolClientSecret(),
		ReplicaRegions:    &replicaRegions,
	})

	awsparams.Store(scope, "UserPoolIDParam", paramsNamespace, "user-pool-id", con.userPoolID)
	awsparams.Store(scope, "UserPoolDomainParam", paramsNamespace, "user-pool-domain", jsii.Sprintf("%s-auth", qual))
	awsparams.Store(scope, "AdminClientIDParam", paramsNamespace, "admin-client-id", con.adminClientID)

	return con
}

func (i *identity) UserPool() awscognito.IUserPool {
	return i.userPool
}

func (i *identity) UserPoolID() *string {
	return i.userPoolID
}

func (i *identity) MainResourceServer() awscognito.IUserPoolResourceServer {
	return i.mainResourceServer
}

func (i *identity) CognitoDomain() *string {
	return i.cognitoDomain
}

func (i *identity) AdminClientID() *string {
	return i.adminClientID
}
