#!/usr/bin/env sh
set -eu

repo_dir="$(pwd)"
cd "${repo_dir}"

run_mapgen() {
	if command -v mapgen >/dev/null 2>&1; then
		mapgen "$@"
		return
	fi

	go run github.com/kitti12911/lib-orm/v3/cmd/mapgen@v3.5.0 "$@"
}

rm -rf gen/grpc gen/database
rm -f internal/feature/user/mapper_generated.go
buf generate
run_mapgen fields -model-dir internal/database -root User -out gen/database/fieldmap_generated.go -package database
run_mapgen patch -config internal/feature/user/patchfields.yaml
run_mapgen proto -config protomapgen.yaml
