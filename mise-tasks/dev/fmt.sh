#!/usr/bin/env bash
#MISE description="Format all code in the project"
#MISE depends=["dev:gen"]
set -euo pipefail

# for the backend package
cd backend
go mod tidy
go tool gofumpt -w .

# across backend/infra.
cd -
shfmt -w mise-tasks/**/*.sh
yamlfmt .
buf format -w

# for the infra package
cd infra
go tool gofumpt -w .
