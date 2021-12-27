create: 
	protoc --go_out=. proto/*.proto
	protoc --go-grpc_out=. proto/*.proto
	protoc -I . --grpc-gateway_out=. \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    proto/vote.proto

clean:
	rm proto/*go

# mkdir -p google/api
# curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
# curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto