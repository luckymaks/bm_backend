package awsdns

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/luckymaks/bm_backend/infra/aws/awsparams"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

const paramsNamespace = "dns"

type dns struct {
	hostedZone awsroute53.IHostedZone
}

type DNS interface {
	HostedZone() awsroute53.IHostedZone
}

type DNSProps struct {
	ZoneDomainName *string
}

func New(scope constructs.Construct, props DNSProps) DNS {
	scope, con := constructs.NewConstruct(scope, jsii.String("DNS")), &dns{}

	if cdkutil.IsPrimaryRegion(scope) {
		con.hostedZone = awsroute53.NewHostedZone(scope, jsii.String("HostedZone"), &awsroute53.HostedZoneProps{
			ZoneName: props.ZoneDomainName,
		})

		awsparams.Store(scope, "HostedZoneIDParam", paramsNamespace, "hosted-zone-id", con.hostedZone.HostedZoneId())
	} else {
		hostedZoneID := awsparams.Lookup(scope,
			"LookupHostedZoneID", paramsNamespace, "hosted-zone-id", "hosted-zone-id-lookup")
		con.hostedZone = awsroute53.HostedZone_FromHostedZoneAttributes(scope, jsii.String("HostedZone"),
			&awsroute53.HostedZoneAttributes{
				HostedZoneId: hostedZoneID,
				ZoneName:     props.ZoneDomainName,
			})
	}

	return con
}

func (con dns) HostedZone() awsroute53.IHostedZone {
	return con.hostedZone
}
