# grpc-sandbox

Go gRPC sandbox service for the homelab API experiments. The service currently
implements user CRUD against PostgreSQL with Bun, shared protobuf contracts from
[`proto-sandbox`](https://github.com/kitti12911/proto-sandbox), and shared
helpers from [`lib-util`](https://github.com/kitti12911/lib-util),
[`lib-orm`](https://github.com/kitti12911/lib-orm), and
[`lib-monitor`](https://github.com/kitti12911/lib-monitor).

## features

- gRPC `user.v1.UserService`
- get, list, create, update, patch, and soft delete users
- list pagination, filtering, and ordering
- PATCH support through `google.protobuf.FieldMask`
- PostgreSQL migrations and seed fixtures
- OpenTelemetry tracing, structured logs, and optional Pyroscope profiling
- gRPC error responses with `trace_id` for log correlation

## requirements

- go 1.26 or higher
- [buf](https://buf.build/) for protobuf generation
- PostgreSQL for local runtime

Optional:

- [prettier](https://prettier.io/) for Markdown, YAML, JSON, and JSONC formatting

## ci commands

reusable CI entrypoints live in `scripts/ci/` so GitHub Actions and GitLab CI
can call the same commands with provider-specific orchestration around them.

| command                                            | purpose                                       |
| -------------------------------------------------- | --------------------------------------------- |
| `./scripts/ci/generate-code.sh`                    | generate protobuf, field-map, and PATCH code  |
| `./scripts/ci/go-lint.sh`                          | run `go vet` and `golangci-lint`              |
| `./scripts/ci/go-test.sh`                          | run tests with filtered coverage              |
| `./scripts/ci/markdownlint.sh`                     | run Markdown linting                          |
| `./scripts/ci/security-scan.sh`                    | run `govulncheck` and Semgrep                 |
| `./scripts/ci/supply-chain-scan.sh`                | run Trivy and Gitleaks                        |
| `./scripts/ci/semantic-release-plan.sh`            | preview the next semantic release             |
| `./scripts/ci/semantic-release-publish.sh`         | publish the semantic release                  |
| `./scripts/ci/fast-forward-prerelease-branches.sh` | fast-forward `uat` and `develop` after `main` |
| `./scripts/ci/update-helm-image-values.sh`         | update homelab GitOps image values            |

GitHub Actions uses `TOOLCHAIN_REGISTRY` and `TOOLCHAIN_IMAGE_NAMESPACE` to
resolve shared CI toolchain images, and `IMAGE_REGISTRY` plus `IMAGE_NAMESPACE`
to publish the application image. GitLab uses full image references so the
private mirror can point at Harbor without changing these scripts:

| GitLab variable                   | Purpose                                     |
| --------------------------------- | ------------------------------------------- |
| `CI_IMAGE_TOOLCHAIN_IMAGE`        | Image for generation, Go lint/test, builds  |
| `CI_SECURITY_TOOLCHAIN_IMAGE`     | Image for `govulncheck` and Semgrep         |
| `CI_SUPPLY_CHAIN_TOOLCHAIN_IMAGE` | Image for Trivy and Gitleaks                |
| `CI_RELEASE_TOOLCHAIN_IMAGE`      | Image for Markdownlint and semantic-release |
| `CI_DOCKER_CLI_IMAGE`             | Docker CLI image for build/publish jobs     |
| `CI_DOCKER_DIND_IMAGE`            | Docker-in-Docker service image              |
| `CI_TRIVY_RUNNER_IMAGE`           | Optional Trivy runner image override        |
| `IMAGE_REGISTRY`                  | Target application image registry           |
| `IMAGE_NAMESPACE`                 | Target application image namespace          |
| `GITLAB_AMD64_RUNNER_TAG`         | Optional runner tag override                |
| `GL_TOKEN` or `GITLAB_TOKEN`      | GitLab semantic-release API/write token     |

| GitLab secret                        | Purpose                  |
| ------------------------------------ | ------------------------ |
| `IMAGE_REGISTRY_USERNAME`            | Target registry username |
| `IMAGE_REGISTRY_PASSWORD`            | Target registry password |
| `COSIGN_PRIVATE_KEY` or `COSIGN_KEY` | Image signing key        |

The `homelab-devops` values update in `.github/workflows/go-ci.yml` is
GitHub-specific homelab orchestration, not part of the portable script contract.
The prerelease branch fast-forward helper is also GitHub-specific because it
pushes through a GitHub App token.
GitLab deployments can use a different project, folder layout, or deployment
tool by calling the same `scripts/ci` build/release helpers and adding its own
deploy job. `DEPLOY_IMAGE_REGISTRY` and `DEPLOY_IMAGE_NAMESPACE` only affect the
homelab GitOps values update and can be omitted outside that workflow.

`GO_TEST_RACE=true` or `GO_TEST_CGO=true` requires a C compiler in the selected
toolchain image. `grpc-sandbox` sets `GO_TEST_RACE=false` in GitHub Actions
while using `image-toolchain` v1.1.0 because that image does not include one.

## project structure

```bash
grpc-sandbox/
├── cmd/
│   └── server/                 # gRPC server entrypoint
├── gen/
│   ├── database/               # generated Bun field maps
│   └── grpc/                   # generated protobuf clients
├── internal/
│   ├── apperror/               # application error wrapper
│   ├── config/                 # config structs
│   ├── database/               # models, migrations, seeders
│   ├── feature/
│   │   └── user/               # user handler, service, repository
│   └── server/                 # gRPC server and interceptors
├── buf.gen.yaml
├── config.example.yml
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

## configuration

Copy `config.example.yml` to `config.yml` and adjust local values:

```bash
cp config.example.yml config.yml
```

Important sections:

- `service`: service name, gRPC port, and shutdown timeout
- `logging`: slog level and trace id injection
- `tracing`: OTLP exporter settings
- `profiling`: Pyroscope settings
- `database`: PostgreSQL connection, migrations, seeders, and pool settings

## generate code

```bash
make gen
```

`make gen` runs:

- protobuf generation from
  [`github.com/kitti12911/proto-sandbox`](https://github.com/kitti12911/proto-sandbox)
  pinned in `buf.gen.yaml`
- Bun field-map generation into `gen/database`
- PATCH field-mask extraction generation into
  `internal/feature/user/patch_generated.go`

The generated field maps and patch extractor come from
[`github.com/kitti12911/lib-orm/v2`](https://github.com/kitti12911/lib-orm)
generator commands.

Generator notes:

- `fieldmapgen` reads Bun models under `internal/database` and generates field
  maps plus validator functions in `gen/database`.
- `patchfieldgen` reads `internal/feature/user/user.go` and generates
  `patchFields(params PatchParams)`.
- `-root-selector params.User` means patch values are read from `params.User`.
- `-paths-selector params.Fields` means field mask paths are read from
  `params.Fields`.
- `-bucket root:userFields:fieldmap.IsUserRootField` routes top-level paths
  such as `email` into `data.userFields`.
- `-bucket profile:profileFields:fieldmap.IsUserProfileField` routes paths
  such as `profile.first_name` into `data.profileFields`.
- `-bucket profile.address:addressFields:fieldmap.IsUserAddressField` routes
  paths such as `profile.address.city` into `data.addressFields`.
- `-copy params.User.Profile:data.profile` copies the full profile value when
  present, so PATCH can create a missing profile row before updating it.
- `-copy params.User.Profile.Address:data.address:params.User.Profile` copies
  address with a profile nil guard, so generated code does not dereference a
  nil profile.

In short, buckets create SQL update maps, while copies carry nested create data
for create-if-missing PATCH flows.

## run locally

```bash
make run
```

The server reads `config.yml`, initializes logging, tracing, profiling,
database migrations/seeders, and starts the gRPC server on `service.port`.

## user API

Implemented RPCs:

- `GetUser`
- `ListUsers`
- `CreateUser`
- `UpdateUser`
- `PatchUser`
- `DeleteUser`

`ListUsers` accepts common `Filter`, `OrderBy`, and `PaginationRequest`
messages. Filters and order fields are validated against generated database
field maps before SQL is built.

`PatchUser` uses a field mask to decide which values to update. The generated
patch extractor splits fields into user, profile, and address update buckets so
each table is patched separately.

## available commands

| Command            | Description                                           |
| ------------------ | ----------------------------------------------------- |
| `make air`         | Run the service with Air live reload                  |
| `make tidy`        | Run `go mod tidy`                                     |
| `make run`         | Start the gRPC server locally                         |
| `make lint`        | Run Go and Markdown linting                           |
| `make ci-generate` | Run CI code generation                                |
| `make ci-lint`     | Run CI Go linting                                     |
| `make ci-test`     | Run CI tests with filtered coverage                   |
| `make fmt`         | Format Go code with `go fmt`                          |
| `make pretty`      | Format Markdown, YAML, JSON, and JSONC                |
| `make format`      | Run Go and document/config formatting                 |
| `make test`        | Run tests with the race detector                      |
| `make cov`         | Generate and open an HTML coverage report             |
| `make fix`         | Apply standard Go source rewrites with `go fix`       |
| `make gen`         | Generate protobuf clients, field maps, and PATCH code |
| `make gen-go`      | Generate database field maps and PATCH helper code    |
| `make gen-proto`   | Generate protobuf clients from `proto-sandbox`        |
