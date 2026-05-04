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

```bash
make tidy
make fmt
make test
make cov
make gen
make run
```
