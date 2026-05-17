#!/usr/bin/env sh
set -eu

repo_dir="${CI_PROJECT_DIR:-$(pwd)}"
cd "${repo_dir}"

run_fieldmapgen() {
	if command -v fieldmapgen >/dev/null 2>&1; then
		fieldmapgen "$@"
		return
	fi

	go run github.com/kitti12911/lib-orm/v3/cmd/fieldmapgen@v3.0.1 "$@"
}

run_patchfieldgen() {
	if command -v patchfieldgen >/dev/null 2>&1; then
		patchfieldgen "$@"
		return
	fi

	go run github.com/kitti12911/lib-orm/v3/cmd/patchfieldgen@v3.0.1 "$@"
}

rm -rf gen/grpc gen/database
buf generate
run_fieldmapgen -model-dir internal/database -root User -out gen/database/fieldmap_generated.go -package database
run_patchfieldgen -config internal/feature/user/patchfields.yaml
