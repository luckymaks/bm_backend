package aws

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/luckymaks/bm_backend/infra/aws/awsapi"
)

type DeploymentProps struct {
	DeploymentIdent *string
	HostedZone      awsroute53.IHostedZone      // optional: nil if no custom domain
	Certificate     awscertificatemanager.ICertificate // optional: nil if no custom domain
	Identity        awscognito.UserPool
	CrewIdentity    awscognito.UserPool
}

type Deployment interface {
	Api() awsapi.Api
}

type deployment struct {
	api awsapi.Api
}

func NewDeployment(stack awscdk.Stack, props DeploymentProps) Deployment {
	api := awsapi.NewApi(stack, awsapi.APIProps{
		DeploymentIdent: props.DeploymentIdent,
		HostedZone:      props.HostedZone,
		Certificate:     props.Certificate,
	})
	
	return &deployment{
		api: api,
	}
}

func (d *deployment) Api() awsapi.Api {
	return d.api
}
