package awscertificate

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type certificate struct {
	certificate awscertificatemanager.ICertificate
}

type Certificate interface {
	WildcardCertificate() awscertificatemanager.ICertificate
}

type CertificateProps struct {
	HostedZone awsroute53.IHostedZone
}

func New(scope constructs.Construct, props CertificateProps) Certificate {
	scope, con := constructs.NewConstruct(scope, jsii.String("Certificate")), &certificate{}

	con.certificate = awscertificatemanager.NewCertificate(scope, jsii.String("WildcardCertificate"),
		&awscertificatemanager.CertificateProps{
			DomainName: jsii.String("*." + *props.HostedZone.ZoneName()),
			Validation: awscertificatemanager.CertificateValidation_FromDns(props.HostedZone),
		})

	return con
}

func (con certificate) WildcardCertificate() awscertificatemanager.ICertificate {
	return con.certificate
}
