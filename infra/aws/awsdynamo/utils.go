package awsdynamo

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
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
