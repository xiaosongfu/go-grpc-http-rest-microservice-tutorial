package v1

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.xiaosongfu.com/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
)

// https://github.com/amsokol/go-grpc-http-rest-microservice-tutorial/blob/part1/pkg/service/v1/todo-service_test.go
func TestCreate(t *testing.T) {

	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connectiong", err)
	}

	defer db.Close()

	toDoServer := NewToDoServiceServer(db)
	timeNow := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(timeNow)

	type args struct {
		ctx     context.Context
		request *v1.CreateRequest
	}

	tests := []struct {
		name    string
		s       v1.ToDoServiceServer
		args    args
		mock    func()
		want    *v1.CreateResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    toDoServer,
			args: args{
				ctx: ctx,
				request: &v1.CreateRequest{
					Api: "v1",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder:    reminder,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO ToDo").WithArgs("title", "description", timeNow).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &v1.CreateResponse{
				Api: "v1",
				Id:  1,
			},
		},
		{
			name: "Unsupported Api",
			s:    toDoServer,
			args: args{
				ctx: ctx,
				request: &v1.CreateRequest{
					Api: "v1000",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder: &timestamp.Timestamp{
							Seconds: 1,
							Nanos:   -1,
						},
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Invalid Reminder field format",
			s:    toDoServer,
			args: args{
				ctx: ctx,
				request: &v1.CreateRequest{
					Api: "v1",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder: &timestamp.Timestamp{
							Seconds: 1,
							Nanos:   -1,
						},
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "INSERT failed",
			s:    toDoServer,
			args: args{
				ctx: ctx,
				request: &v1.CreateRequest{
					Api: "v1",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder:    reminder,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO ToDo").WithArgs("title", "description", timeNow).WillReturnError(errors.New("INSERT failed"))
			},
			wantErr: true,
		},
		{
			name: "LatestInsertId failed",
			s:    toDoServer,
			args: args{
				ctx: ctx,
				request: &v1.CreateRequest{
					Api: "v1",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder:    reminder,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO ToDO").WithArgs("title", "description", timeNow).WillReturnResult(sqlmock.NewErrorResult(errors.New("LasterInsertId failed")))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := tt.s.Create(tt.args.ctx, tt.args.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("toDoServiceServer.Create() error =%v, wantErr=%v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDoServiceServer.Create() =%v, wantErr=%v", got, tt.wantErr)
			}
		})
	}
}
