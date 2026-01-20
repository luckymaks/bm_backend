#!/usr/bin/env bash
#MISE description="Check if all code compiles"
set -euo pipefail

# Backend code.
cd backend
go build ./...

# AWS CDK infra.
cd -
cd infra
go build ./...
