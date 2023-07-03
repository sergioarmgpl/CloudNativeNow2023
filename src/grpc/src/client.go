package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "usac.projects/match"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	team1 = flag.String("team1", "team1", "Team1 value")
	team2 = flag.String("team2", "team2", "Team2 value")
	score = flag.String("score", "0-0", "Match score")
)

func main() {
	var phase int32 = 1;

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMatchClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//r, err = c.SendMessage(ctx,&pb.MatchRequest{Team1: *team1,Team2: *team2,Score: *score,Phase: phase})
	r, err := c.SendMessage(ctx,&pb.MatchRequest{Team1: *team1,Team2: *team2,Score: *score,Phase: phase})
	if err != nil {
			log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Match Info: %s %s %s %i",r.GetTeam1(),r.GetTeam2(),r.GetScore(),r.GetPhase())
}