#!/usr/bin/env bash
#MISE description="Bootstrap CDK for the project"
#USAGE arg "<bootstrap_aws_profile>" help="The super-user AWS profile for bootstrapping CDK"
set -euo pipefail

# check access for the bootstrap profile
aws sts get-caller-identity --profile "${usage_bootstrap_aws_profile:?}"

cd infra/aws/cdk
cdk_context=$(cdk context --json)
qualifier=$(echo "$cdk_context" | jq -re '.["kn-qualifier"] // error("missing kn-qualifier context")')
secondary_regions=$(echo "$cdk_context" | jq -re '.["kn-secondary-regions"] // [] | join(",")')

# Determine the directory where the cloudformation template file is.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Deploy the deployer policy stack (--no-fail-on-empty-changeset handles already up-to-date stacks)
aws cloudformation deploy \
	--stack-name "${qualifier}-pre-bootstrap" \
	--template-file "$SCRIPT_DIR/pre-bootstrap.cfn.yaml" \
	--parameter-overrides "Qualifier=$qualifier" "SecondaryRegions=$secondary_regions" \
	--capabilities CAPABILITY_NAMED_IAM \
	--no-fail-on-empty-changeset \
	--profile "${usage_bootstrap_aws_profile:?}"

# Get outputs from the pre-bootstrap stack
get_output() {
	aws cloudformation describe-stacks \
		--stack-name "${qualifier}-pre-bootstrap" \
		--profile "${usage_bootstrap_aws_profile:?}" \
		--query "Stacks[0].Outputs[?OutputKey=='$1'].OutputValue" \
		--output text
}

execution_policy_arn=$(get_output ExecutionPolicyArn)
permissions_boundary_name=$(get_output PermissionBoundaryName)

# Verify that the CDK context has the correct permissions boundary configured. Else deploying CDK code will yield
# very cryptic errors when creating roles/users.
context_boundary=$(echo "$cdk_context" | jq -re '.["@aws-cdk/core:permissionsBoundary"].name // error("missing permission boundary context")')
if [[ "$context_boundary" != "$permissions_boundary_name" ]]; then
	echo "ERROR: CDK context @aws-cdk/core:permissionsBoundary.name must be set to '$permissions_boundary_name'" >&2
	echo "       Current value: '${context_boundary:-<not set>}'" >&2
	exit 1
fi

cdk bootstrap \
	--profile "${usage_bootstrap_aws_profile:?}" \
	--qualifier "$qualifier" \
	--toolkit-stack-name "${qualifier}Bootstrap" \
	--cloudformation-execution-policies "$execution_policy_arn" \
	--custom-permissions-boundary "$permissions_boundary_name"
