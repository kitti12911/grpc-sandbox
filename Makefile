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
	env CGO_ENABLE=1 go test --race -v ./...

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

fix: 
	go fix ./...

# ____________________ Generate Command ____________________
gen: gen-proto

gen-proto:
	rm -rf gen/grpc
	buf generate https://github.com/kitti12911/proto-sandbox.git --path common/v1 --path user/v1