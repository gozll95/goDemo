# lorem-grpc
This is simple service module. Only for showing the micro service with gRPC protocol
The purpose for this service is only generating lorem ipsum paragraph and return the payload.

I am fully using all three functions from the golorem library.

## Required libraries

    go get github.com/go-kit/kit
    go get github.com/drhodes/golorem
    go get github.com/gorilla/mux

# pb
Protocol buffer module. The place to create proto files.
Download protoc from [here](https://github.com/google/protobuf/releases)
Then execute `go get -u github.com/golang/protobuf/{proto,protoc-gen-go}`

*Note: don't forget to add GOBIN on your PATH*

To generate protobuf file into go file:
`protoc lorem.proto --go_out=plugins=grpc:.`

### service.go
Business logic will be put here

### endpoint.go
Endpoint will be created here

### model.go
Encode and Decode json

### transport.go
Implement interface from protocol buffer

### execute

    cd $GOPATH

    #Running grpc server
    go run src/github.com/ru-rocker/gokit-playground/lorem-grpc/server/server_grpc_main.go

    #Running client
    go run src/github.com/ru-rocker/gokit-playground/lorem-grpc/client/cmd/client_grpc_main.go lorem sentence 10 20
