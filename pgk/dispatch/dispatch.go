package dispatch

import (
	"context"
	"log"
	"michaelracz/image-service/pgk/docker"
	"michaelracz/image-service/pgk/queue"
)

func Dispatch(ctx context.Context, queue queue.Queue, dockerClient docker.Client) {
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
