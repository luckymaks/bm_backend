package aws

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/luckymaks/bm_backend/infra/aws/awsapi"
	"github.com/luckymaks/bm_backend/infra/aws/awsdynamo"
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
	dynamo := awsdynamo.NewDynamo(stack, awsdynamo.DynamoProps{
		DeploymentIdent: props.DeploymentIdent,
	})

	api := awsapi.NewApi(stack, awsapi.APIProps{
		DeploymentIdent: props.DeploymentIdent,
		HostedZone:      props.HostedZone,
		Certificate:     props.Certificate,
		MainTable:       dynamo.Table(),
		MainTableName:   dynamo.TableName(),
	})
	
	return &deployment{
		api: api,
	}
}

func (d *deployment) Api() awsapi.Api {
	return d.api
}
