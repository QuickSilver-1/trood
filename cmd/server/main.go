package main

import (
	"fmt"
	"log"
	"os"
	"server/internal/server"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Failed to load env fail")
		return
	}

	nlpHost := os.Getenv("NLP_HOST")
	nlpPortStr := os.Getenv("NLP_PORT")

	nlpPort, err := strconv.Atoi(nlpPortStr)
	if err != nil {
		log.Fatal("Invalid NLP_PORT")
		return
	}

	nlp := server.NewNlpService(nlpHost, nlpPort)
	server.NlpService = nlp

	time.Sleep(time.Second*10)
	fmt.Println(111)
	elasticHost := os.Getenv("ELASTIC_HOST")
	elasticPortStr := os.Getenv("ELASTIC_PORT")

	elasticPort, err := strconv.Atoi(elasticPortStr)
	if err != nil {
		log.Fatal("Invalid ELASTIC_PORT")
		return
	}

	elasticSearch := server.NewElasticService(elasticHost, elasticPort)
	server.ElasticSearchService = elasticSearch
	err = elasticSearch.CreateIndex()
	if err != nil {
		log.Fatal(err)
		return
	}

	srv := server.NewServer(8080)
	err = srv.StartServer()

	if err != nil {
		log.Fatal(err)
	}
}