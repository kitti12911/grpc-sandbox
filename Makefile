# ____________________ Go Command ____________________
air:
	air

tidy:
	go mod tidy

run:
	go run ./cmd/server/main.go

fmt:
	go fmt ./...

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
	go run github.com/kitti12911/lib-orm/v2/cmd/fieldmapgen@v2.0.0 -model-dir internal/database -root User -out gen/database/fieldmap_generated.go -package database
	go run ./cmd/patchfieldgen -file internal/feature/user/user.go -root CreateParams -out internal/feature/user/patch_generated.go -package user -fieldmap-import grpc-sandbox/gen/database

gen-proto:
	rm -rf gen/grpc
	buf generate https://github.com/kitti12911/proto-sandbox.git --path common/v1 --path user/v1
