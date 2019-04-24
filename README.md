13185
222.85.230.14
go_grpc_microservice


```
go run pkg/cmd/server/main.go -grpc-port=9090 -db-host=222.85.230.14 -db-port=13185 -db-user=root -db-password=AllApp -db-schema=go_grpc_microservice

go run pkg/cmd/client_grpc/main.go -server=localhost:9090
```


```
grpc.gateway.protoc_gen_swagger.options.openapiv2_swagger

grpc.gateway.protoc_gen_swagger.options.openapiv2_schema

grpc.gateway.protoc_gen_swagger.options.openapiv2_tag

grpc.gateway.protoc_gen_swagger.options.openapiv2_operation

google.api.http
```


```
go run pkg/cmd/server/main.go -grpc-port=9090 -http-port=9091 -db-host=222.85.230.14 -db-port=13185 -db-user=root -db-password=AllApp -db-schema=go_grpc_microservice

go run pkg/cmd/client_rest/main.go -server=http://localhost:9091
```