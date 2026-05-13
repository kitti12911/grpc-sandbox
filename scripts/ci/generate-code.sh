#!/usr/bin/env sh
set -eu

repo_dir="${CI_PROJECT_DIR:-$(pwd)}"
cd "${repo_dir}"

run_fieldmapgen() {
	if command -v fieldmapgen >/dev/null 2>&1; then
		fieldmapgen "$@"
		return
	fi

	go run github.com/kitti12911/lib-orm/v2/cmd/fieldmapgen@v2.7.0 "$@"
}

run_patchfieldgen() {
	if command -v patchfieldgen >/dev/null 2>&1; then
		patchfieldgen "$@"
		return
	fi

	go run github.com/kitti12911/lib-orm/v2/cmd/patchfieldgen@v2.7.0 "$@"
}

rm -rf gen/grpc gen/database
buf generate
run_fieldmapgen -model-dir internal/database -root User -out gen/database/fieldmap_generated.go -package database
run_patchfieldgen -file internal/feature/user/user.go -root CreateParams -out internal/feature/user/patch_generated.go -package user -fieldmap-import grpc-sandbox/gen/database -root-selector params.User -paths-selector params.Fields -bucket root:userFields:fieldmap.IsUserRootField -bucket profile:profileFields:fieldmap.IsUserProfileField -bucket profile.address:addressFields:fieldmap.IsUserAddressField -copy params.User.Profile:data.profile -copy params.User.Profile.Address:data.address:params.User.Profile
