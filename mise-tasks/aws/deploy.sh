#!/usr/bin/env bash
#USAGE arg "[deployment_ident]" help="Deployment identifier (defaults to Dev<YourUsername>)"
#USAGE flag "--hotswap" help="Enable CDK hotswap deployment for faster iterations"
#MISE description="Compare deployed infra with what is configured"
set -euo pipefail
source "$(dirname "$0")/_cdk-common.sh"

hotswap_flag=""
if [[ "${usage_hotswap:-false}" == "true" ]]; then
	hotswap_flag="--hotswap"
fi

outputs_file=$(mktemp)
trap 'rm -f "$outputs_file"' EXIT

cdk deploy "${cdk_common_args[@]}" \
	--require-approval "never" --outputs-file "$outputs_file" $hotswap_flag "${qualifier}*Shared" "${qualifier}*${deployment_ident}"

echo ""
echo "=== Stack Outputs ==="
cat "$outputs_file"
