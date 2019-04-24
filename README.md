13185
222.85.230.14

go_grpc_microservice


```
go run pkg/cmd/server/main.go -grpc-port=9090 -db-host=222.85.230.14 -db-port=13185 -db-user=root -db-password=AllApp -db-schema=go_grpc_microservice

go run pkg/cmd/client_grpc/main.go -server=localhost:9090
```