#!/usr/bin/env bash
#MISE description="Lint the code base"
set -euo pipefail
find mise-tasks -name '*.sh' -type f -print0 | xargs -0 shellcheck

# Backend code.
cd backend
golangci-lint config verify
golangci-lint run ./...
buf lint

# AWS CDK infra.
cd -
cd infra
golangci-lint run ./...

# cloudformation.
cd -
cfn-lint mise-tasks/aws/pre-bootstrap.cfn.yaml
