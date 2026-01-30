package awsdynamo

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/iancoleman/strcase"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

type DynamoProps struct {
	DeploymentIdent *string
}

type Dynamo = commonTableInterface

type dynamo struct {
	table     awsdynamodb.ITableV2
	tableName *string
}

func NewDynamo(parent constructs.Construct, props DynamoProps) Dynamo {
	scope, con := createScope(parent, "Dynamo"), &dynamo{}
	qual := cdkutil.QualifierFromContext(scope)

	con.tableName = jsii.Sprintf("%s-%s-main-table", qual, strcase.ToKebab(*props.DeploymentIdent))

	if cdkutil.IsPrimaryRegion(scope) {
		replicas := buildReplicaConfigs(scope)

		con.table = awsdynamodb.NewTableV2(scope, jsii.String("MainTable"), &awsdynamodb.TablePropsV2{
			TableName: con.tableName,
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("pk"),
				Type: awsdynamodb.AttributeType_STRING,
			},
			SortKey: &awsdynamodb.Attribute{
				Name: jsii.String("sk"),
				Type: awsdynamodb.AttributeType_STRING,
			},
			Billing:            getBillingOnDemand(),
			RemovalPolicy:      awscdk.RemovalPolicy_DESTROY,
			DeletionProtection: jsii.Bool(false),
			PointInTimeRecoverySpecification: &awsdynamodb.PointInTimeRecoverySpecification{
				PointInTimeRecoveryEnabled: jsii.Bool(true),
				RecoveryPeriodInDays:       jsii.Number(3),
			},
			GlobalSecondaryIndexes: &[]*awsdynamodb.GlobalSecondaryIndexPropsV2{
				{
					IndexName:    jsii.String("gsi1"),
					PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("gsi1pk"), Type: awsdynamodb.AttributeType_STRING},
					SortKey:      &awsdynamodb.Attribute{Name: jsii.String("gsi1sk"), Type: awsdynamodb.AttributeType_STRING},
				},
			},
			Replicas: replicas,
		})
	} else {
		con.table = awsdynamodb.TableV2_FromTableName(scope, jsii.String("MainTable"), con.tableName)
	}

	return con
}

func (d *dynamo) Table() awsdynamodb.ITableV2 {
	return d.table
}

func (d *dynamo) TableName() *string {
	return d.tableName
}
