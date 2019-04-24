package v1

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.xiaosongfu.com/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
)

const (
	// 由服务器指定的版本号
	apiVersion = "v1"
)

type toDoServiceServer struct {
	db *sql.DB
}

func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{db: db}
}

// checkAPI 检测客户端请求的 api 版本是否被服务器支持
func (t *toDoServiceServer) checkAPI(api string) error {
	// "" 版本号意味着使用现在的版本
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}

	return nil
}

// connect  从数据库连接池返回一个数据库连接
// connect returns SQL database connection from the pool
func (t *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	conn, err := t.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return conn, nil
}

//-----------

// Create 创建新 task
func (t *toDoServiceServer) Create(ctx context.Context, in *v1.CreateRequest) (*v1.CreateResponse, error) {
	// 检查客户端请求的 api 版本是否被支持
	if err := t.checkAPI(in.Api); err != nil {
		return nil, err
	}

	// 从数据库连接池中获取连接
	// get SQL connection from pool
	conn, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	reminder, err := ptypes.Timestamp(in.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format->"+err.Error())
	}

	res, err := conn.ExecContext(ctx, "INSERT INTO ToDo(`Title`, `Description`, `Reminder`) VALUES (?,?,?)",
		in.ToDo.Title, in.ToDo.Description, reminder)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into ToDo->"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created ToDo->"+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

// Read 读取 task
func (t *toDoServiceServer) Read(ctx context.Context, in *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := t.checkAPI(in.Api); err != nil {
		return nil, err
	}

	conn, err := t.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	rows, err := conn.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo WHERE `ID`=?", in.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo->"+err.Error())
	}

	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo->"+err.Error())
		}
		return nil, status.Errorf(codes.NotFound, "ToDo with ID='%d' is not found", in.Id)
	}

	var todo v1.ToDo
	var reminder time.Time
	if err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row->"+err.Error())
	}

	todo.Reminder, err = ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format->"+err.Error())
	}

	if rows.Next() {
		return nil, status.Errorf(codes.Unknown, "found multiple ToDo rows with ID='%d'", in.Id)
	}
	return &v1.ReadResponse{
		Api:  apiVersion,
		ToDo: &todo,
	}, nil
}

func (t *toDoServiceServer) Update(ctx context.Context, in *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := t.checkAPI(in.Api); err != nil {
		return nil, err
	}

	conn, err := t.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	reminder, err := ptypes.Timestamp(in.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format->"+err.Error())
	}

	res, err := conn.ExecContext(ctx, "UPDATE ToDo SET `Title`=?, `Description`=?, `Reminder`=? WHERE `ID`=?", in.ToDo.Title, in.ToDo.Description, reminder, in.ToDo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo->"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value->"+err.Error())
	}

	if rows == 0 {
		return nil, status.Errorf(codes.NotFound, "ToDo with ID='%d' is not found", in.ToDo.Id)
	}

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil

}

func (t *toDoServiceServer) Delete(ctx context.Context, in *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := t.checkAPI(in.Api); err != nil {
		return nil, err
	}

	conn, err := t.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	res, err := conn.ExecContext(ctx, "DELETE FROM ToDo WHERE `ID`=?", in.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete ToDo->"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value->"+err.Error())
	}

	if rows == 0 {
		return nil, status.Errorf(codes.NotFound, "ToDo with ID='%d' is not found", in.Id)
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: rows,
	}, nil

}

func (t *toDoServiceServer) ReadAll(ctx context.Context, in *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := t.checkAPI(in.Api); err != nil {
		return nil, err
	}

	conn, err := t.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	rows, err := conn.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo->"+err.Error())
	}

	defer rows.Close()

	var reminder time.Time
	var list []*v1.ToDo

	for rows.Next() {
		todo := new(v1.ToDo)

		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row->"+err.Error())
		}

		todo.Reminder, err = ptypes.TimestampProto(reminder)
		if err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has invalid format->"+err.Error())
		}

		list = append(list, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo->"+err.Error())
	}

	return &v1.ReadAllResponse{
		Api:   apiVersion,
		ToDos: list,
	}, nil
}
