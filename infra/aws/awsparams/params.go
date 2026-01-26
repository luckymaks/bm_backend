package awsparams

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/aws-cdk-go/awscdk/v2/customresources"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

func ParameterName(scope constructs.Construct, namespace string, name string) *string {
	qual := cdkutil.QualifierFromContext(scope)
	return jsii.Sprintf("/%s/%s/%s", qual, namespace, name)
}

func Store(scope constructs.Construct, id string, namespace string, name string, value *string) {
	awsssm.NewStringParameter(scope, jsii.String(id),
		&awsssm.StringParameterProps{
			ParameterName: ParameterName(scope, namespace, name),
			StringValue:   value,
		})
}

func Lookup(scope constructs.Construct, id string, namespace string, name string, physicalID string) *string {
	lookup := customresources.NewAwsCustomResource(scope, jsii.String(id),
		&customresources.AwsCustomResourceProps{
			OnCreate: &customresources.AwsSdkCall{
				Service: jsii.String("SSM"),
				Action:  jsii.String("getParameter"),
				Parameters: map[string]any{
					"Name": ParameterName(scope, namespace, name),
				},
				Region:             jsii.String(cdkutil.PrimaryRegion(scope)),
				PhysicalResourceId: customresources.PhysicalResourceId_Of(jsii.String(physicalID)),
			},
			Policy: customresources.AwsCustomResourcePolicy_FromSdkCalls(&customresources.SdkCallsPolicyOptions{
				Resources: customresources.AwsCustomResourcePolicy_ANY_RESOURCE(),
			}),
		})
	return lookup.GetResponseField(jsii.String("Parameter.Value"))
}
