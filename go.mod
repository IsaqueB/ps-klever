module github.com/IsaqueB/ps-klever

go 1.16

require (
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/mongo-driver v1.8.1
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.2
	google.golang.org/genproto v0.0.0-20211223182754-3ac035c7e7cb
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0
)

// +heroku goVersion go1.17.1
