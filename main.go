package main

import (
	"github.com/gin-contrib/cors"
)

func main() {
	dependencies := NewDependencies()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	dependencies.engine.Use(cors.New(config))
	_ = dependencies.Run()
}