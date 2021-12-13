package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 256 << 10 // 256 KiB should be enough for a Dockerfile
	router.POST("/build", func(c *gin.Context) {
		var err error
		file, err := c.FormFile("Dockerfile")
		if err != nil {
			handleError(c, err, http.StatusBadRequest)
			return
		}
		openedFile, err := file.Open()
		if err != nil {
			handleError(c, err, http.StatusInternalServerError)
			return
		}
		fileBytes, err := ioutil.ReadAll(openedFile)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError)
			return
		}
		log.Println(string(fileBytes))
		c.JSON(http.StatusAccepted, map[string]string{"status": "success"})
	})
	return router
}

func handleError(c *gin.Context, err error, httpStatus int) {
	log.Printf("ERROR: %v\n", err)
	c.JSON(httpStatus, map[string]string{
		"status": "failure",
		"error":  err.Error(),
	})
}
