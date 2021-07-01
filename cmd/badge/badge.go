package main

import (
	"github/dirien/sonarqube-private-badges/pkg/config"
	"github/dirien/sonarqube-private-badges/pkg/server"
	"log"
)

func main() {
	config := config.LoadConfig()
	proxy := server.NewSonarQubeProxy(config)
	s := proxy.Server()
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
