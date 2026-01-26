package cdkutil

import (
	"slices"
	"strings"
	
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// DeployerGroups returns the IAM groups of the current deployer, or nil if not set.
// This context is passed by _cdk-common.sh during deploy/diff operations.
// During bootstrap (run by admins), this context is not set and returns nil.
func DeployerGroups(scope constructs.Construct) []string {
	val := scope.Node().TryGetContext(jsii.String("kn-deployer-groups"))
	if val == nil {
		return nil
	}
	
	str, ok := val.(string)
	if !ok || str == "" {
		return nil
	}
	
	return strings.Fields(str)
}

func HasDeployerGroup(scope constructs.Construct, group string) bool {
	return slices.Contains(DeployerGroups(scope), group)
}

// AllowedDeployments returns the list of deployments the current deployer is allowed to deploy.
// If kn-deployer-groups context is not set (e.g., during bootstrap), no deployments are allowed
// since bootstrap only needs the CDK toolkit, not application stacks.
// Otherwise, Stag/Prod require membership in the 'bstr-deployers' group.
// Restricted deployments are filtered out, so attempting to deploy them will result in a
// "stack not found" error from CDK.
func AllowedDeployments(scope constructs.Construct) []string {
	all := Deployments(scope)
	groups := DeployerGroups(scope)
	
	// No group context provided (e.g., during bootstrap), skip deployment stacks
	if groups == nil {
		return nil
	}
	
	isFull := HasDeployerGroup(scope, "kndr-deployers")
	if isFull {
		return all
	}
	
	// Filter out Stag/Prod for non-full deployers
	allowed := make([]string, 0, len(all))
	for _, d := range all {
		if d == "Stag" || d == "Prod" {
			continue
		}
		allowed = append(allowed, d)
	}
	return allowed
}

func Deployments(scope constructs.Construct) []string {
	val := scope.Node().GetContext(jsii.String("kn-deployments"))
	if val == nil {
		panic("invalid 'kn-deployments', is it set?")
	}
	
	slice, ok := val.([]any)
	if !ok {
		panic("invalid 'kn-deployments', expected array")
	}
	
	regions := make([]string, 0, len(slice))
	for _, v := range slice {
		s, ok := v.(string)
		if !ok {
			panic("invalid 'kn-deployments', expected array of strings")
		}
		regions = append(regions, s)
	}
	
	return regions
}
