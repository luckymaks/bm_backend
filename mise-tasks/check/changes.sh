#!/usr/bin/env bash
#MISE description="Check generated code is checked-in"
#MISE depends=["dev:fmt", "dev:gen"]
set -euo pipefail

if [[ "${CI:-false}" == "true" ]]; then
	# Exclude platform-specific files (.terraform.lock.hcl differs between macOS and Linux)
	changes=$(git status --porcelain | grep -v '.terraform.lock.hcl' || true)
	if [[ -n "$changes" ]]; then
		echo "ERROR: Code is not up to date."
		echo "Run generating tasks locally and commit the changes."
		echo
		echo "$changes"
		exit 1
	fi
fi
