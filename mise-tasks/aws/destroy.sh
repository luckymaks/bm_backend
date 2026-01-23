#!/usr/bin/env bash
#USAGE arg "[deployment_ident]" help="Deployment identifier (defaults to Dev<YourUsername>)"
#MISE description="Compare deployed infra with what is configured"
set -euo pipefail
source "$(dirname "$0")/_cdk-common.sh"
cdk destroy "${cdk_common_args[@]}" "${qualifier}*${deployment_ident}"
