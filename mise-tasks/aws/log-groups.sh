#!/usr/bin/env bash
#USAGE arg "[deployment_ident]" help="Deployment identifier (defaults to Dev<YourUsername>)"
#MISE description="Returns all CloudWatch log groups for a deployment across all regions"
set -euo pipefail

usage_deployment_ident="${usage_deployment_ident:-}"
source "$(dirname "${BASH_SOURCE[0]}")/_cdk-common.sh"

primary_region=$(echo "$cdk_context" | jq -re '.["kn-primary-region"]')
secondary_regions=$(echo "$cdk_context" | jq -re '.["kn-secondary-regions"] // []')
all_regions=$(echo "$secondary_regions" | jq -re --arg pr "$primary_region" '. + [$pr] | .[]')

found=0
for region in $all_regions; do
	region_ident=$(echo "$cdk_context" | jq -re --arg r "$region" '.["kn-region-ident-" + $r]')
	stack_name="${qualifier}${region_ident}${deployment_ident}"

	# Filter out CDK internal constructs (LogicalResourceId starting with 'AWS')
	groups=$(aws cloudformation list-stack-resources \
		--profile "$deployer_profile" \
		--region "$region" \
		--stack-name "$stack_name" \
		--query "StackResourceSummaries[?ResourceType=='AWS::Logs::LogGroup' && !starts_with(LogicalResourceId, 'AWS')].PhysicalResourceId" \
		--output json 2>/dev/null || echo "[]")

	while IFS= read -r group; do
		if [[ -n "$group" ]]; then
			echo "${region}: ${group}"
			found=1
		fi
	done < <(echo "$groups" | jq -re '.[]')
done

if [[ $found -eq 0 ]]; then
	echo "No log groups found for deployment: $deployment_ident" >&2
fi
