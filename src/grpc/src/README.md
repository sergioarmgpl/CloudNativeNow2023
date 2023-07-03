export PATH="$PATH:$(go env GOPATH)/bin"
export GO111MODULE=on
go mod init grpc-football-match
brew install protobuf
brew install go

Install protoc-gen-go

$ go get google.golang.org/protobuf/cmd/protoc-gen-go
$ go install google.golang.org/protobuf/cmd/protoc-gen-go

Install protoc-gen-go-grpc

$ go get google.golang.org/grpc/cmd/protoc-gen-go-grpc 
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

Install grpc

$ go get google.golang.org/grpc
mkdir match;
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    match/match.proto 


protoc --go-grpc_out=require_unimplemented_servers=false:./match/ --go_out=./match/ match.proto

go mod tidy


https://go.dev/doc/tutorial/create-module