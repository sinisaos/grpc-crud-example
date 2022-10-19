package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sinisaos/grpc-crud-example/models"
	pb "github.com/sinisaos/grpc-crud-example/proto"
	"google.golang.org/grpc"
)

var DB *gorm.DB

type server struct{ pb.TodoServiceServer }

func (*server) ListTodos(ctx context.Context, r *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	todos := make([]*pb.Todo, 0)
	DB.Find(&todos)
	return &pb.ListTodosResponse{
		Todos: todos,
	}, nil
}

func (*server) ReadTodo(ctx context.Context, r *pb.ReadTodoRequest) (*pb.ReadTodoResponse, error) {
	id := r.GetId()
	todo := &pb.Todo{}
	DB.First(&todo, id)
	return &pb.ReadTodoResponse{
		Todo: todo,
	}, nil
}

func (*server) DeleteTodo(ctx context.Context, r *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	id := r.GetId()
	todo := &pb.Todo{}
	DB.Delete(&todo, id)
	return &pb.DeleteTodoResponse{
		Success: true,
	}, nil
}

func (*server) CreateTodo(ctx context.Context, r *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	name := r.GetName()
	completed := r.GetCompleted()
	todo := &pb.Todo{Name: name, Completed: completed}
	DB.Create(&todo)
	return &pb.CreateTodoResponse{
		Todo: todo,
	}, nil
}

func (*server) UpdateTodo(ctx context.Context, r *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	id := r.GetId()
	name := r.GetName()
	completed := r.GetCompleted()
	todo := &pb.Todo{}
	DB.Model(&todo).Updates(map[string]interface{}{"id": id, "name": name, "completed": completed})
	return &pb.UpdateTodoResponse{
		Todo: todo,
	}, nil
}

func main() {
	fmt.Println("Starting gRPC server")
	database, err := gorm.Open("sqlite3", "../gprc_todo.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&models.Todo{})

	DB = database

	l, e := net.Listen("tcp", ":50051")
	if e != nil {
		log.Fatalf("Failed to start listener %v", e)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &server{})

	if e := s.Serve(l); e != nil {
		log.Fatalf("failed to serve %v", e)
	}
}
