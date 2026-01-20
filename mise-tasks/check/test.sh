#!/usr/bin/env bash
#MISE description="Test the codebase"
set -euo pipefail
cd backend
go test ./...

cd -
cd infra/aws
go test ./...
