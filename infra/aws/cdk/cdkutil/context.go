package cdkutil

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func QualifierFromContext(scope constructs.Construct) string {
	qual := StringContext(scope, "kn-qualifier")
	if len(qual) > 10 { // https://github.com/aws/aws-cdk/pull/10121/files
		panic(fmt.Sprintf("CDK qualifier became too large (>10): '%s', adjust context.", qual))
	}

	return qual
}

func RegionAcronymIdentFromContext(scope constructs.Construct, region string) string {
	return StringContext(scope, "kn-region-ident-"+region)
}

func StringContext(scope constructs.Construct, key string) string {
	qual, ok := scope.Node().GetContext(jsii.String(key)).(string)
	if !ok {
		panic("invalid '" + key + "', is it set?")
	}

	return qual
}

func BaseDomainName(scope constructs.Construct) *string {
	return jsii.String(StringContext(scope, "kn-base-domain-name"))
}
