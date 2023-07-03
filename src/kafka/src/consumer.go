package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
    "encoding/json"
	"context"
	//"flag"
	"log"

    "github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "usac.projects/match"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Match struct {
    Team1     string  `json:"team2"`
    Team2  string  `json:"team1"`
    Score string  `json:"score"`
    Phase  int32 `json:"phase"`
}

type MatchMongo struct {
    Team1     string  `bson:"team2,omitempty"`
    Team2  string  `bson:"team1,omitempty"`
    Score string  `bson:"score,omitempty"`
    Phase  int32 `bson:"phase,omitempty"`
}

/*var (
	addr = flag.String("addr", ":50051", "the address to connect to")
)*/

func sendMongo(team1 string, team2 string,score string,phase int32) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

    coll := client.Database("project").Collection("matches")
    newMatch := MatchMongo{Team1: team1, Team2: team2,Score: score, Phase:phase}
    result, err := coll.InsertOne(context.TODO(), newMatch)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", result)

}


func sendgRPC(team1 string, team2 string,score string,phase int32) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("GRPC_SERVER"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMatchClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//r, err = c.SendMessage(ctx,&pb.MatchRequest{Team1: *team1,Team2: *team2,Score: *score,Phase: phase})
	r, err := c.SendMessage(ctx,&pb.MatchRequest{Team1: team1,Team2: team2,Score: score,Phase: phase})
	if err != nil {
			log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Match Info: %s %s %s %i",r.GetTeam1(),r.GetTeam2(),r.GetScore(),r.GetPhase())

}



func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-0.my-cluster-kafka-brokers.kafka.svc:9092",
		"group.id": "mygroupid",
		"auto.offset.reset" : "earliest"})

    if err != nil {
        fmt.Printf("Failed to create consumer: %s", err)
        os.Exit(1)
    }

    topic := "matches"
    err = c.SubscribeTopics([]string{topic}, nil)
    // Set up a channel for handling Ctrl-C, etc
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    // Process messages
    run := true
    for run {
        select {
        case sig := <-sigchan:
            fmt.Printf("Caught signal %v: terminating\n", sig)
            run = false
        default:
            ev, err := c.ReadMessage(100 * time.Millisecond)
            if err != nil {
                // Errors are informational and automatically handled by the consumer
                continue
            }
            fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
                *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
            
            match := Match{}
            err2 := json.Unmarshal([]byte(string(ev.Value)), &match)
            if err2 != nil {
                fmt.Println(err2)
                return
            }
            //send to Mongo
            sendMongo(match.Team1,match.Team2,match.Score,match.Phase)
            //send to gRPC
            sendgRPC(match.Team1,match.Team2,match.Score,match.Phase)
        }
    }

    c.Close()

}