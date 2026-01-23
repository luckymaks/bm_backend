package aws

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
)

type SharedProps struct{}

type Shared interface {
	DNS() DNS
	Certificate() Certificate
	Identity() awscognito.UserPool
	CrewIdentity() awscognito.UserPool
}

type DNS interface {
	HostedZone() awsroute53.IHostedZone
}

type Certificate interface {
	WildcardCertificate() awscertificatemanager.ICertificate
}

type shared struct {
	dns          DNS
	certificate  Certificate
	identity     awscognito.UserPool
	crewIdentity awscognito.UserPool
}

func NewShared(stack awscdk.Stack, props SharedProps) Shared {
	return &shared{}
}

func (s *shared) DNS() DNS {
	return s.dns
}

func (s *shared) Certificate() Certificate {
	return s.certificate
}

func (s *shared) Identity() awscognito.UserPool {
	return s.identity
}

func (s *shared) CrewIdentity() awscognito.UserPool {
	return s.crewIdentity
}
