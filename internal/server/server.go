package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)


var (
	NlpService *Nlp
	ElasticSearchService *ElasticSearch
)


type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) StartServer() error {
	router := gin.Default()
	router.GET("/faq", Faq)

	err := router.Run(fmt.Sprintf(":%d", s.port))
	return err
}
