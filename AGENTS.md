
# General Guidelines
- The project is a mono-repo that includes frontend,backend and infrastructure code in the same repository
- All backend development happens in ./backend
- All infrastructure development happens in ./infra

# Checking Guidelines
- After changes to the code, run the followinig command to check if it compiles: `mise r check:compiles`
- If you want to run unit tests, use the following command: `mise r check:test`
- Whenever code works, run all code (quality) checks with: `mise r 'check:*' and fix them`

# Model Implementation Guidelines (backend)
- Keep `backend/internal/model` structs transport-agnostic:
  - Create a struct for Input and Output types.
  - Create a struct for dynamoDB type.

# Model + RPC Implementation Guidelines (backend)
- Keep `backend/internal/model` transport-agnostic:
  - Do NOT return Connect errors from model code.
  - Return domain/sentinel errors (e.g. `model.ErrValidation`, `model.ErrNotFound`, `model.ErrAlreadyExists`) and wrap details.
- Use `backend/internal/rpc/modelconv` as the single place for:
  - proto request -> model input conversions `modelconv.FromProto[model.InputType](req)`
  - model output -> proto response conversions `modelconv.ToProto(out, &bstrv1.Response{})`
  - model error -> `connect.Error` mapping

# Snapshot Testing Guidelines (backend)
- Prefer snapshot tests where it improves reliability and reviewability:
  - Snapshot normalized outputs (avoid timestamps/random IDs directly), or
  - Snapshot normalized persistence state.
- Snapshots live under `snapshot/` in the package that owns the test.
- Update snapshots by running tests with `UPDATE_SNAPSHOTS=1`.

# AWS CDK Development
- When you want to see the diff for new infrastructure to be deployed, use the `mise r aws:diff` script without arguments
- When you want to deploy new infrastructure, use the `mise r aws:deploy` script without arguments
- To test if the CDK code compiles, make sure to first change directories: `cd infra/aws`
- Never deploy to production (Prod deployment) yourself
- Never deploy to staging (Stag deployment) yourself
- Never manually deploy or change a resource that is managed by the CDK, or Cloudformation.

# AWS Development
- Whenever you want to run AWS CLI command to debug code, use the "kndr-admin" profile.
- Whenever you want to find log groups for the deployment, run: `mise r aws:log-groups`

# Project Structure
```
bm_backend/
├── backend/                    # Backend application code
│   ├── internal/               # Internal packages (not exported)
│   │   └── rpc/                # RPC layer and model conversions
│   └── lambda/                 # Lambda function entry points
│       └── httpapi/            # HTTP API Lambda (Echo framework)
├── infra/                      # Infrastructure code
│   └── aws/                    # AWS CDK infrastructure
│       ├── cdk/                # CDK app entry point and utilities
│       │   └── cdkutil/        # CDK helper utilities (context, stack, regions)
│       ├── awsapi/             # API Gateway + Lambda construct
│       ├── awsdynamo/          # DynamoDB table construct
│       ├── awsdns/             # Route53 DNS construct
│       ├── awscertificate/     # ACM certificate construct
│       ├── awsidentity/        # Cognito identity construct
│       ├── awslambda/          # Lambda-related constructs
│       ├── awsparams/          # SSM parameters construct
│       ├── awss3/              # S3 bucket construct
│       ├── awssecret/          # Secrets Manager construct
│       ├── deployment.go       # Per-deployment resources (API, DynamoDB)
│       └── shared.go           # Shared resources (DNS, certs, identity)
└── mise-tasks/                 # Mise task runner scripts
    ├── check/                  # Code quality checks (compiles, lint, test)
    ├── dev/                    # Development tasks (fmt, gen)
    └── aws/                    # AWS deployment tasks (diff, deploy, destroy)
```

# Technology Stack
- **Language**: Go 1.25.3
- **HTTP Framework**: Echo v4 (`github.com/labstack/echo/v4`)
- **AWS SDK**: aws-sdk-go-v2 for DynamoDB and other AWS services
- **Infrastructure**: AWS CDK v2 (Go bindings)
- **Database**: DynamoDB with single-table design (pk/sk pattern with GSI)
- **Compute**: AWS Lambda with Lambda Web Adapter for Echo compatibility
- **API Gateway**: HTTP API (API Gateway v2)
- **Linting**: golangci-lint, buf (for protobuf), shellcheck
- **Formatting**: gofumpt, shfmt, yamlfmt, buf format, terraform fmt
- **Task Runner**: mise

# Development Commands
- `mise r dev:gen` - Generate code (go generate, buf generate)
- `mise r dev:fmt` - Format all code (runs dev:gen first)
- `mise r check:compiles` - Check if code compiles
- `mise r check:lint` - Run all linters
- `mise r check:test` - Run unit tests
- `mise r 'check:*'` - Run all checks

# AWS CDK Patterns
- CDK constructs follow interface + struct pattern (e.g., `type Api interface {...}` + `type api struct {...}`)
- Constructor naming: `NewXxx(parent constructs.Construct, props XxxProps) Xxx`
- Use `cdkutil.QualifierFromContext(scope)` for resource naming prefixes
- Use `cdkutil.IsPrimaryRegion(scope)` to handle primary vs secondary region logic
- Multi-region deployment with primary/secondary stack dependencies
- Deployment identifiers (Dev/Stag/Prod) used for resource naming

# DynamoDB Patterns
- Single-table design with `pk` (partition key) and `sk` (sort key)
- Key pattern: `TYPE#id` (e.g., `ITEM#123`)
- GSI1 available with `gsi1pk` and `gsi1sk` attributes
- Table is global with multi-region replication

# Lambda Patterns
- Uses Lambda Web Adapter to run Echo as a Lambda function
- Lambda listens on port 12001 (configured via `AWS_LWA_PORT`)
- Environment variables passed from CDK (e.g., `MAIN_TABLE_NAME`)
- ARM64 architecture for cost efficiency
