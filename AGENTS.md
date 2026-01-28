
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
