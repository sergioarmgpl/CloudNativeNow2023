// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
    "os"
    "strconv"

	"google.golang.org/grpc"
	pb "usac.projects/match"
	"github.com/redis/go-redis/v9"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedMatchServer
}


var ctx = context.Background()
	
func sendRedis(team1 string, team2 string,score string,phase int32) {
	rdb := redis.NewClient(&redis.Options{
		Addr:    os.Getenv("REDIS_SERVER"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	err := rdb.HIncrBy(ctx,"teams:counters", team1+"-"+team2,1).Err()
	err2 := rdb.IncrBy(ctx,"messages",1).Err()
	err3 := rdb.HIncrBy(ctx,"teams:phases", strconv.FormatInt(int64(phase), 10),1).Err()

	if err != nil {
		panic(err)
	}

	if err2 != nil {
		panic(err)
	}

	if err3 != nil {
		panic(err)
	}
}


// SayHello implements helloworld.GreeterServer
func (s *server) SendMessage(ctx context.Context, in *pb.MatchRequest) (*pb.MatchReply, error) {
	log.Printf("Received: %v", in.GetTeam1())
	log.Printf("Received: %v", in.GetTeam2())
	log.Printf("Received: %v", in.GetScore())
	log.Printf("Received: %v", in.GetPhase())
	//Writing into Redis
	sendRedis(in.GetTeam1(),in.GetTeam2(),in.GetScore(),in.GetPhase())
	return &pb.MatchReply{Team1: in.GetTeam1(),Team2: in.GetTeam2(),Score: in.GetScore(),Phase: in.GetPhase()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMatchServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}