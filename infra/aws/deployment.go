package aws

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
)

type DeploymentProps struct {
	DeploymentIdent *string
	HostedZone      awsroute53.IHostedZone
	Certificate     awscertificatemanager.ICertificate
	Identity        awscognito.UserPool
	CrewIdentity    awscognito.UserPool
}

type Deployment interface{}

type deployment struct{}

func NewDeployment(stack awscdk.Stack, props DeploymentProps) Deployment {
	return &deployment{}
}
