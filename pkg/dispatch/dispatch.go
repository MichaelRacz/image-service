package dispatch

import (
	"context"
	"log"
	"michaelracz/image-service/pkg/docker"
	"michaelracz/image-service/pkg/queue"
)

// Dispatch takes Dockerfiles from a queue, builds and pushes the image
func Dispatch(ctx context.Context, queue queue.Dequeueer, dockerClient docker.Client) {
	for {
		df, ok := queue.Dequeue(ctx)
		if !ok {
			break
		}

		var err error
		tag := dockerClient.CreateTag()

		if err = dockerClient.BuildImage(ctx, df, tag); err == nil {
			err = dockerClient.PushImage(ctx, tag)
			// NOTE: Compensating the build of the orphaned image is needed here
		}

		if err != nil {
			log.Printf("ERROR: %v\n", err)
		}
	}
}
