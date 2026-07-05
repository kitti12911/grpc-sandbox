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

reusable CI entrypoints live in `scripts/ci/` so GitHub Actions can call the
same commands with workflow-specific orchestration around them.

| command                                            | purpose                                              |
| -------------------------------------------------- | ---------------------------------------------------- |
| `./scripts/ci/generate-code.sh`                    | generate protobuf, field-map, PATCH, and mapper code |
| `./scripts/ci/go-lint.sh`                          | run `go vet` and `golangci-lint`                     |
| `./scripts/ci/go-test.sh`                          | run tests with filtered coverage                     |
| `./scripts/ci/markdownlint.sh`                     | run Markdown linting                                 |
| `./scripts/ci/security-scan.sh`                    | run `govulncheck` and Semgrep                        |
| `./scripts/ci/supply-chain-scan.sh`                | run Trivy and Gitleaks                               |
| `./scripts/ci/semantic-release-publish.sh`         | publish the semantic release                         |
| `./scripts/ci/fast-forward-prerelease-branches.sh` | fast-forward `uat` and `develop` after `main`        |
| `./scripts/ci/update-helm-image-values.sh`         | update homelab GitOps image values                   |

GitHub Actions uses `TOOLCHAIN_REGISTRY` and `TOOLCHAIN_IMAGE_NAMESPACE` to
resolve shared CI toolchain images, and `IMAGE_REGISTRY` plus `IMAGE_NAMESPACE`
to publish the application image. `DEPLOY_IMAGE_REGISTRY` and
`DEPLOY_IMAGE_NAMESPACE` only affect the homelab GitOps values update and can be
omitted outside that workflow.

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
- proto mapper generation into `internal/feature/user/mapper_generated.go`

The generators come from
[`github.com/kitti12911/lib-orm/v4`](https://github.com/kitti12911/lib-orm)
and are **zero-config** — they discover everything by naming convention (see
the lib-orm README for the conventions and `//mapgen:*` directives).

Generator notes:

- `mapgen fields` reads Bun models under `internal/database` and generates
  field/column maps in `gen/database`.
- `mapgen patch` finds `PatchParams` and generates the `patchData` struct plus
  `patchFields(params PatchParams)`; buckets and nil-guarded copies are derived
  from the payload struct's `field:"..."`-tagged shape.
- `mapgen filter` generates `applyFilter`/`applyOrderBy` plus the custom-filter
  registry (`//mapgen:filter col=<name>` functions back virtual columns).
- `mapgen map` reads both the params structs and the generated proto types,
  then emits mappers by field intersection — including the string↔enum
  bridges. The worker feature has no bun root model, so the generator skips it
  and its hand-written fallible mapper stays.

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
| `make fmt`         | Format Go code with `go fmt`                          |
| `make pretty`      | Format Markdown, YAML, JSON, and JSONC                |
| `make format`      | Run Go and document/config formatting                 |
| `make test`        | Run tests with the race detector                      |
| `make cov`         | Generate and open an HTML coverage report             |
| `make fix`         | Apply standard Go source rewrites with `go fix`       |
| `make gen`         | Generate protobuf clients, field maps, and PATCH code |
| `make gen-go`      | Generate database field maps and PATCH helper code    |
| `make gen-proto`   | Generate protobuf clients from `proto-sandbox`        |
