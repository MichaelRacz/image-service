package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	router := setupRouter()

	body, contentType := createCorrectBody(t)
	req, err := http.NewRequest("POST", "/build", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, `{"status":"success"}`, rec.Body.String())
}

func TestBuildClientError(t *testing.T) {
	router := setupRouter()

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
	assert.Contains(t, rec.Body.String(), `"status":"failure"`)
}

// TODO: Add missing error test cases

func createCorrectBody(t *testing.T) (body io.Reader, contentType string) {
	return createBody(t, "Dockerfile")
}

func createInvalidBody(t *testing.T) (body io.Reader, contentType string) {
	return createBody(t, "InvalidFieldName")
}

func createBody(t *testing.T, fieldName string) (body io.Reader, contentType string) {
	var err error
	file, _ := os.Open("./test-dockerfiles/alpine-dockerfile")
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
