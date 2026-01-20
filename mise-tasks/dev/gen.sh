#!/usr/bin/env bash
#MISE description="Generate code across the project"
set -euo pipefail

cd backend
go generate ./...
buf generate
