#!/usr/bin/env sh
set -eu

repo_dir="$(pwd)"
cd "${repo_dir}"

run_mapgen() {
	if command -v mapgen >/dev/null 2>&1; then
		mapgen "$@"
		return
	fi

	go run github.com/kitti12911/lib-orm/v4/cmd/mapgen@v4.0.0 "$@"
}

rm -rf gen/grpc gen/database
buf generate
run_mapgen fields
run_mapgen patch
run_mapgen filter
run_mapgen map
