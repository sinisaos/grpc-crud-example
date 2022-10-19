package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/sinisaos/grpc-todo/proto"
	"github.com/sinisaos/grpc-todo/schemas"
	"google.golang.org/grpc"
)

func main() {

	r := gin.Default()

	r.GET("/", AllTodos)
	r.GET("/:id", SingleTodo)
	r.POST("/", CreateTodo)
	r.PATCH("/:id", UpdateTodo)
	r.DELETE("/:id", DeleteTodo)

	r.Run(":8080")
}

func AllTodos(c *gin.Context) {
	conn, e := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if e != nil {
		log.Fatalf("Failed to connect to server %v", e)
	}
	defer conn.Close()
	client := pb.NewTodoServiceClient(conn)
	todos, e := client.ListTodos(context.Background(), &pb.ListTodosRequest{})
	if e != nil {
		log.Fatalf("Failed to get all todo data: %v", e)
	}

	c.JSON(http.StatusOK, todos)
}

func SingleTodo(c *gin.Context) {
	conn, e := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if e != nil {
		log.Fatalf("Failed to connect to server %v", e)
	}
	defer conn.Close()
	client := pb.NewTodoServiceClient(conn)
	i, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	id := uint64(i)
	todo, e := client.ReadTodo(context.Background(), &pb.ReadTodoRequest{Id: id})
	if e != nil {
		log.Fatalf("Failed to get single todo data: %v", e)
	}

	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
	conn, e := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if e != nil {
		log.Fatalf("Failed to connect to server %v", e)
	}
	defer conn.Close()
	client := pb.NewTodoServiceClient(conn)
	i, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	id := uint64(i)
	todo, e := client.DeleteTodo(context.Background(), &pb.DeleteTodoRequest{Id: id})
	if e != nil {
		log.Fatalf("Failed to delete todo data: %v", e)
	}

	c.JSON(http.StatusOK, todo)
}

func CreateTodo(c *gin.Context) {
	conn, e := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if e != nil {
		log.Fatalf("Failed to connect to server %v", e)
	}
	defer conn.Close()
	client := pb.NewTodoServiceClient(conn)
	var input schemas.CreateTodoIn
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, e := client.CreateTodo(context.Background(), &pb.CreateTodoRequest{Name: input.Name, Completed: input.Completed})
	if e != nil {
		log.Fatalf("Failed to create todo data: %v", e)
	}

	c.JSON(http.StatusOK, todo)
}

func UpdateTodo(c *gin.Context) {
	conn, e := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if e != nil {
		log.Fatalf("Failed to connect to server %v", e)
	}
	defer conn.Close()
	client := pb.NewTodoServiceClient(conn)
	i, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	id := uint64(i)
	var input schemas.UpdateTodoIn
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, e := client.UpdateTodo(context.Background(), &pb.UpdateTodoRequest{Id: id, Name: input.Name, Completed: input.Completed})
	if e != nil {
		log.Fatalf("Failed to update todo data: %v", e)
	}

	c.JSON(http.StatusOK, todo)
}
