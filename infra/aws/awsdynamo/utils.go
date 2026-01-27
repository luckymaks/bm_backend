package awsdynamo

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

type commonTableInterface interface {
	Table() awsdynamodb.ITableV2
	TableName() *string
}

func createScope(parent constructs.Construct, name string) constructs.Construct {
	return constructs.NewConstruct(parent, jsii.String(name))
}

func getBillingOnDemand() awsdynamodb.Billing {
	return awsdynamodb.Billing_OnDemand(&awsdynamodb.MaxThroughputProps{
		MaxReadRequestUnits:  jsii.Number(10),
		MaxWriteRequestUnits: jsii.Number(10),
	})
}

func buildReplicaConfigs(scope constructs.Construct) *[]*awsdynamodb.ReplicaTableProps {
	secondaryRegions := cdkutil.SecondaryRegions(scope)
	if len(secondaryRegions) == 0 {
		return nil
	}
	
	replicas := make([]*awsdynamodb.ReplicaTableProps, 0, len(secondaryRegions))
	for _, region := range secondaryRegions {
		replicas = append(replicas, &awsdynamodb.ReplicaTableProps{
			Region: jsii.String(region),
		})
	}
	
	return &replicas
}
