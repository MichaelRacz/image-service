package api

import (
	"errors"
	"io/ioutil"
	"log"
	"michaelracz/image-service/pkg/docker"
	"michaelracz/image-service/pkg/queue"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(queue queue.Enqueueer) *gin.Engine {
	router := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 256 << 10 // 256 KiB should be enough for a Dockerfile

	router.POST("/build", func(c *gin.Context) {
		dockerfile, err := readDockerfile(c)

		// NOTE: Some kind of validation needs to be done here
		if err != nil {
			handleWebError(c, err)
			return
		}

		if ok := queue.Enqueue(dockerfile); ok {
			c.JSON(http.StatusAccepted, map[string]string{"status": "success"})
		} else {
			err := &WebError{
				errors.New("cannot accept request, queue limit exceeded"),
				http.StatusServiceUnavailable,
			}

			handleWebError(c, err)
		}
	})

	return router
}

func readDockerfile(c *gin.Context) (docker.Dockerfile, *WebError) {
	var err error

	file, err := c.FormFile("Dockerfile")
	if err != nil {
		return "", &WebError{err, http.StatusBadRequest}
	}

	openedFile, err := file.Open()
	if err != nil {
		return "", &WebError{err, http.StatusInternalServerError}
	}

	fileBytes, err := ioutil.ReadAll(openedFile)
	if err != nil {
		return "", &WebError{err, http.StatusInternalServerError}
	}

	return docker.NewDockerfile(fileBytes), nil
}

type WebError struct {
	NestedErr  error
	HttpStatus int
}

func handleWebError(c *gin.Context, err *WebError) {
	log.Printf("ERROR: %v\n", err)

	c.JSON(err.HttpStatus, map[string]string{
		"status": "error",
		"error":  err.NestedErr.Error(),
	})
}
