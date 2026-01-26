package awssecret

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

type secret struct {
	mainSecret awssecretsmanager.ISecret
}

type Secret interface {
	MainSecret() awssecretsmanager.ISecret
}

type SecretProps struct{}

func New(scope constructs.Construct, _ SecretProps) Secret {
	scope, con := constructs.NewConstruct(scope, jsii.String("Secret")), &secret{}
	qual := cdkutil.QualifierFromContext(scope)
	
	con.mainSecret = awssecretsmanager.Secret_FromSecretNameV2(scope,
		jsii.String("LookupMainSecret"),
		jsii.Sprintf("%s/main-secret", qual))
	
	return con
}

func (c secret) MainSecret() awssecretsmanager.ISecret {
	return c.mainSecret
}
