#!/usr/bin/env bash
# shellcheck disable=SC2034 # Variables are used by scripts that source this file
script_dir="$(dirname "${BASH_SOURCE[0]}")"
cd "$script_dir/../../infra/aws/cdk" || exit
cdk_context=$(cdk context --json)
qualifier=$(echo "$cdk_context" | jq -re '.["kn-qualifier"] // error("missing kn-qualifier context")')

# Check if we're running in CI (GitHub Actions sets CI=true)
if [[ "${CI:-}" == "true" ]]; then
	echo "Running in CI mode (using OIDC credentials)"
	profile_args=()

	# CI role has full deployer permissions
	deployer_groups="kndr-deployers"
	default_deployment=""
	echo "Deployer: CI (via OIDC role)"
else
	deployer_profile=$(echo "$cdk_context" | jq -r '.["kn-kindred-deployer-profile"] // empty')
	if [[ -z "$deployer_profile" ]]; then
		echo "Error: 'kn-kindred-deployer-profile' is not configured in CDK context." >&2
		echo "Please add it to ~/.cdk.json, e.g.:" >&2
		echo '  { "kn-kindred-deployer-profile": "your-aws-profile-name" }' >&2
		exit 1
	fi
	profile_args=(--profile="$deployer_profile")

	caller_arn=$(aws sts get-caller-identity --profile "$deployer_profile" --query "Arn" --output text)
	caller_username=$(echo "$caller_arn" | sed -n 's|.*:user/||p')

	if [[ -n "$caller_username" ]]; then
		deployer_groups=$(aws iam list-groups-for-user --user-name "$caller_username" --profile "$deployer_profile" --query "Groups[].GroupName" --output text)
		echo "Deployer: $caller_username (groups: $deployer_groups)"
		default_deployment="Dev${caller_username}"
	else
		deployer_groups=""
		default_deployment=""
		echo "Deployer: $caller_arn (not an IAM user, cannot determine groups)"
	fi
fi

# Resolve deployment identifier: use provided arg or default to Dev<username>
deployment_ident="${usage_deployment_ident:-$default_deployment}"
if [[ -z "$deployment_ident" ]]; then
	echo "Error: No deployment specified and could not determine default (not an IAM user)" >&2
	exit 1
fi

cdk_common_args=(
	"${profile_args[@]}"
	--qualifier "$qualifier"
	--toolkit-stack-name "${qualifier}Bootstrap"
	--context "kn-deployer-groups=$deployer_groups"
)
