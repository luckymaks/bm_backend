package cdkutil

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func AllRegions(scope constructs.Construct) []string {
	return append([]string{PrimaryRegion(scope)}, SecondaryRegions(scope)...)
}

func IsPrimaryRegion(scope constructs.Construct) bool {
	return *awscdk.Stack_Of(scope).Region() == PrimaryRegion(scope)
}

func PrimaryRegion(scope constructs.Construct) string {
	return StringContext(scope, "kn-primary-region")
}

func SecondaryRegions(scope constructs.Construct) []string {
	val := scope.Node().TryGetContext(jsii.String("kn-secondary-regions"))
	if val == nil {
		return nil
	}

	slice, ok := val.([]any)
	if !ok {
		panic("invalid 'kn-secondary-regions', expected array")
	}

	regions := make([]string, 0, len(slice))
	for _, v := range slice {
		s, ok := v.(string)
		if !ok {
			panic("invalid 'kn-secondary-regions', expected array of strings")
		}
		regions = append(regions, s)
	}

	return regions
}
