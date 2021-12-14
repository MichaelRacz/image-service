package main

import (
	"michaelracz/image-service/pgk/api"
	"michaelracz/image-service/pgk/queue"
)

func main() {
	router := api.SetupRouter(queue.NewQueue(10))
	router.Run(":8080")
}
