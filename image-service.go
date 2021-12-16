package main

import (
	"context"
	"log"
	"michaelracz/image-service/pkg/api"
	"michaelracz/image-service/pkg/dispatch"
	"michaelracz/image-service/pkg/docker"
	"michaelracz/image-service/pkg/queue"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const QUEUE_SIZE = 10
const SHUTDOWN_GRACE_PERIOD = 5 * time.Second

func main() {
	dispatcherCtx, dispatcherCancel := context.WithCancel(context.Background())
	defer dispatcherCancel()
	// NOTE: As a simplification, queue, dispatcher, and api run
	// in the same process. For production usage, a queue as a
	// service and splitting api and dispatching is a more robust
	// and scalable solution.
	q := queue.NewQueue(QUEUE_SIZE)
	startDispatcher(dispatcherCtx, q)
	srv := startApi(q)

	waitForQuit()

	log.Println("Shutdown image service ...")
	go shutdownApi(srv)
	dispatcherCancel()
	time.Sleep(SHUTDOWN_GRACE_PERIOD)
	log.Println("Exiting")
	os.Exit(0)
}

func startDispatcher(ctx context.Context, q queue.Dequeueer) {
	registryUserID := os.Getenv("REGISTRY_USER_ID")
	password := os.Getenv("REGISTRY_PASSWORD")
	registryUrl := os.Getenv("REGISTRY_URL")
	dc, err := docker.NewClient(registryUserID, password, registryUrl)
	if err != nil {
		log.Printf("ERROR: Cannot start image-service: %v\n", err)
		os.Exit(1)
	}

	go dispatch.Dispatch(ctx, q, dc)
}

func startApi(q queue.Enqueueer) *http.Server {
	router := api.SetupRouter(q)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("ERROR: Api error: %s\n", err)
		}
	}()

	return srv
}

func waitForQuit() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func shutdownApi(srv *http.Server) {
	apiCtx, apiCancel := context.WithTimeout(context.Background(), SHUTDOWN_GRACE_PERIOD)
	defer apiCancel()
	if err := srv.Shutdown(apiCtx); err != nil {
		log.Printf("ERROR: Cannot shut down image-service api: %v\n", err)
	}
}
