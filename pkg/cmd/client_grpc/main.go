package main

import (
	"context"
	"flag"
	"github.com/golang/protobuf/ptypes"
	v1 "go.xiaosongfu.com/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	apiVerison = "v1"
)

func main() {
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := v1.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	prefixx := t.Format(time.RFC3339Nano)

	// Create
	createRequest := &v1.CreateRequest{
		Api: apiVerison,
		ToDo: &v1.ToDo{
			Title: "title (" + prefixx + ")",
			Description: "description (" + prefixx + ")",
			Reminder: reminder,
		},
	}
	createResponse, err := client.Create(ctx, createRequest)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", createResponse)


	// 保存 id
	id := createResponse.Id

	// Read
	readRequest := &v1.ReadRequest{
		Api: apiVerison,
		Id: id,
	}
	readResponse, err := client.Read(ctx, readRequest)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}

	log.Printf("Read result: <%+v>\n\n", readResponse)

	// Update
	updateRequest := &v1.UpdateRequest{
		Api: apiVerison,
		ToDo: &v1.ToDo{
			Id: readResponse.ToDo.Id,
			Title: readResponse.ToDo.Title,
			Description: readResponse.ToDo.Description + " + udpate",
			Reminder: readResponse.ToDo.Reminder,
		},
	}
	updateResponse, err := client.Update(ctx, updateRequest)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	log.Printf("Update result: <%+v>\n\n", updateResponse)


	// ReadAll
	readAllRequest := &v1.ReadAllRequest{
		Api: apiVerison,
	}
	readAllResponse, err := client.ReadAll(ctx, readAllRequest)
	if err != nil {
		log.Fatalf("ReadAll failed: %v", err)
	}
	log.Printf("ReadAll result: <%+v>\n\n", readAllResponse)


	// Delete
	deleteRequest := &v1.DeleteRequest{
		Api:apiVerison,
		Id: id,
	}
	deleteResponse, err := client.Delete(ctx, deleteRequest)
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	log.Printf("Delete result: <%+v>\n\n", deleteResponse)
}


