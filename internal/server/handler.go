package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pb "server/internal/server/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Nlp struct {
	host string
	port int
}

func NewNlpService(host string, port int) *Nlp {
    return &Nlp{
		host: host,
		port: port,
	}
}

func Faq(c *gin.Context) {
	question := c.Query("question")
	if question == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'question' is required"})
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", NlpService.host, NlpService.port), grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()

    client := pb.NewKeywordExtractorClient(conn)
	answer, err := client.ExtractKeywords(context.Background(), &pb.KeywordRequest{
		Question: question,
	})
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	result, err := ElasticSearchService.GetFaq(answer)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong"})
		return
	}

	if result == nil {
		c.JSON(http.StatusOK, gin.H{"Info": "There are no similar questions in our database, the question has been sent to the manager"})
		return
	}

	c.JSON(http.StatusOK, result)
}