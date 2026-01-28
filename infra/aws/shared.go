package aws

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/luckymaks/bm_backend/infra/aws/awscertificate"
	"github.com/luckymaks/bm_backend/infra/aws/awsdns"
	"github.com/luckymaks/bm_backend/infra/aws/awsidentity"
	"github.com/luckymaks/bm_backend/infra/aws/awssecret"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

type shared struct {
	dns         awsdns.DNS
	certificate awscertificate.Certificate
	secret      awssecret.Secret
	identity    awsidentity.Identity
}

type Shared interface {
	DNS() awsdns.DNS
	Certificate() awscertificate.Certificate
	Secret() awssecret.Secret
	Identity() awsidentity.Identity
}

type SharedProps struct {
	EnableCustomDomain bool
}

func NewShared(scope constructs.Construct, props SharedProps) Shared {
	con := &shared{}
	
	if props.EnableCustomDomain {
		con.dns = awsdns.New(scope, awsdns.DNSProps{
			ZoneDomainName: cdkutil.BaseDomainName(scope),
		})
		
		con.certificate = awscertificate.New(scope, awscertificate.CertificateProps{
			HostedZone: con.dns.HostedZone(),
		})
	}
	
	con.secret = awssecret.New(scope, awssecret.SecretProps{})
	con.identity = awsidentity.New(scope, awsidentity.IdentityProps{})
	
	return con
}

func (s *shared) DNS() awsdns.DNS {
	return s.dns
}

func (s *shared) Certificate() awscertificate.Certificate {
	return s.certificate
}

func (s *shared) Secret() awssecret.Secret {
	return s.secret
}

func (s *shared) Identity() awsidentity.Identity {
	return s.identity
}
