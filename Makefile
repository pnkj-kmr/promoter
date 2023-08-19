clean_rpc:
	rm -rf medium/pb/*.go

gen: clean_rpc
	protoc --proto_path=medium/proto --go_out=. --go-grpc_out=. medium/proto/*.proto

run:
	go run main.go

test:
	go run tmp/test/main.go

h1:
	curl http://localhost:8000/debug/pprof/heap --output heap1.tar.gz

p1:
	curl  http://localhost:8000/debug/pprof/profile  --output profile1.tar.gz

h2:
	curl http://localhost:8001/debug/pprof/heap --output heap2.tar.gz

p2:
	curl  http://localhost:8001/debug/pprof/profile  --output profile2.tar.gz

tidy:
	go mod tidy
	
cover:
	go test -cover -race ./...

release:
	goreleaser release --snapshot --clean 

build:
	goreleaser build --single-target --snapshot --clean
