# ____________________ Go Command ____________________
air:
	air

tidy:
	go mod tidy

run:
	go run ./cmd/server/main.go

lint: vet golangci-lint markdownlint

vet:
	go vet ./...

golangci-lint:
	golangci-lint run --timeout=5m

markdownlint:
	markdownlint-cli2

fmt:
	go fmt ./...

pretty:
	prettier --write "**/*.{md,markdown,yml,yaml,json,jsonc}"

format: fmt pretty

test:
	env CGO_ENABLED=1 go test --race -v ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

fix: 
	go fix ./...

# ____________________ Generate Command ____________________
gen: gen-proto gen-go

gen-go:
	rm -rf gen/database
	go run github.com/kitti12911/lib-orm/v2/cmd/fieldmapgen@v2.7.0 -model-dir internal/database -root User -out gen/database/fieldmap_generated.go -package database
	go run github.com/kitti12911/lib-orm/v2/cmd/patchfieldgen@v2.7.0 -file internal/feature/user/user.go -root CreateParams -out internal/feature/user/patch_generated.go -package user -fieldmap-import grpc-sandbox/gen/database -root-selector params.User -paths-selector params.Fields -bucket root:userFields:fieldmap.IsUserRootField -bucket profile:profileFields:fieldmap.IsUserProfileField -bucket profile.address:addressFields:fieldmap.IsUserAddressField -copy params.User.Profile:data.profile -copy params.User.Profile.Address:data.address:params.User.Profile

gen-proto:
	rm -rf gen/grpc
	buf generate
