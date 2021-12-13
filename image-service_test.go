package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"michaelracz/image-service/pgk/queue"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const TEST_DOCKERFILE = "./test-dockerfiles/alpine-dockerfile"

func TestBuild(t *testing.T) {
	queue := queue.NewQueue(1)
	router := setupRouter(queue)

	rec := requestBuild(t, router)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, `{"status":"success"}`, rec.Body.String())
	expectedFile := fileAsString(t, TEST_DOCKERFILE)
	enqueuedFile := <-queue.GetChannel()
	assert.Equal(t, expectedFile, enqueuedFile)
}

func TestBuildClientError(t *testing.T) {
	router := setupRouter(queue.NewQueue(1))

	body, contentType := createInvalidBody(t)
	req, err := http.NewRequest("POST", "/build", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), `"status":"error"`)
}

func TestBuildQueueLimitExceeded(t *testing.T) {
	router := setupRouter(queue.NewQueue(1))

	requestBuild(t, router)
	rec := requestBuild(t, router)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), `"status":"error"`)
}

// TODO: Add missing error test cases

func requestBuild(t *testing.T, router *gin.Engine) *httptest.ResponseRecorder {
	body, contentType := createCorrectBody(t)
	req, err := http.NewRequest("POST", "/build", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func createCorrectBody(t *testing.T) (body io.Reader, contentType string) {
	return createBody(t, "Dockerfile")
}

func createInvalidBody(t *testing.T) (body io.Reader, contentType string) {
	return createBody(t, "InvalidFieldName")
}

func createBody(t *testing.T, fieldName string) (body io.Reader, contentType string) {
	var err error
	file, _ := os.Open(TEST_DOCKERFILE)
	defer file.Close()
	var buffer bytes.Buffer
	multiPartWriter := multipart.NewWriter(&buffer)
	var fileWriter io.Writer
	if fileWriter, err = multiPartWriter.CreateFormFile(fieldName, file.Name()); err != nil {
		t.Fatal(err)
	}
	if _, err = io.Copy(fileWriter, file); err != nil {
		t.Fatal(err)
	}
	multiPartWriter.Close()
	return &buffer, multiPartWriter.FormDataContentType()
}

func fileAsString(t *testing.T, filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}
