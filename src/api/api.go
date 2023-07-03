package main

import (
    "net/http"
	"fmt"
    "os"
    "strconv"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/gin-gonic/gin"
)

type match struct {
    Team1     string  `json:"team2"`
    Team2  string  `json:"team1"`
    Score string  `json:"score"`
    Phase  int32 `json:"phase"`
}

type result struct {
    Processed     string  `json:"processed"`
}

func main() {
    router := gin.Default()
	router.GET("/_health", getHealth)
    router.POST("/match", postMatch)
    router.Run(":3000")
}

func getHealth(c *gin.Context) {
	var   newResult result
	newResult.Processed = "done"
    c.IndentedJSON(http.StatusOK,newResult)
}

func postMatch(c *gin.Context) {
    var newMatch match
    if err := c.BindJSON(&newMatch); err != nil {
        return
    }
/*	fmt.Println("team1",newMatch.Team1)
	fmt.Println("team2",newMatch.Team2)
	fmt.Println("score",newMatch.Score)
	fmt.Println("phase",newMatch.Phase)
*/
//**************
    topic := "matches"
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "my-cluster-kafka-0.my-cluster-kafka-brokers.kafka.svc:9092"})
    if err != nil {
        fmt.Printf("Failed to create producer: %s", err)
        os.Exit(1)
    }

    p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Key: []byte("match"),
        Value: []byte("{\"team1\":\""+newMatch.Team1+"\",\"team2\":\""+newMatch.Team2+"\",\"score\":\""+newMatch.Score+"\",\"phase\":"+strconv.FormatInt(int64(newMatch.Phase), 10)+"}"),
    }, nil)

    // Wait for all messages to be delivered
    //p.Flush(15 * 1000)
    p.Flush(500)
    p.Close()
//**************


    c.IndentedJSON(http.StatusCreated, newMatch)
}
