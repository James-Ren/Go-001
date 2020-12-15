package main

import (
	"context"
	"log"
	"time"

	pb "github.com/James-Ren/Go-001/tree/main/Week04/api/article/v1"
	"github.com/James-Ren/Go-001/tree/main/Week04/internal/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("../../configs")
	err := viper.ReadInConfig()
	if err != nil {
		// load config error
		panic(err)
	}
	srv, cleanup, err := server.InitializeServer()
	defer cleanup()
	if err != nil {
		log.Printf("Init Server error:%v\n", err)
		return
	}
	// go testGrpc()

	log.Println("Start Server")
	if err = srv.Run(); err != nil {
		log.Printf("Run Server error:%v\n", err)
		return
	}
}

func testGrpc() {
	conn, err := grpc.Dial(viper.GetString("grpc.port"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewArticleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetArticle(ctx, &pb.ArticleRequest{Id: 2})
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	log.Printf("article: %d:%s", r.GetId(), r.GetTitle())
}
