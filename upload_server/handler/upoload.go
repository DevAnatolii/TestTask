package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testTask/upload_server/model"
	"testTask/upload_server/parse"
)

const (
	HandlePath           = "/upload/"
	detectTypeByteLength = 512 // check file type, detectcontenttype only needs the first 512 bytes
)

type uploadHandler struct {
	personsServerBaseUrl string
}

func NewUploadHandler(personsServerBaseUrl string) *uploadHandler {
	return &uploadHandler{personsServerBaseUrl}
}

func (uH *uploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse and validate file and post parameters
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		renderError(w, "Invalid file", http.StatusBadRequest)
		return
	}

	defer file.Close()

	fileBytes := make([]byte, detectTypeByteLength, detectTypeByteLength)
	io.ReadFull(file, fileBytes)
	file.Seek(detectTypeByteLength, io.SeekStart)
	fileType := http.DetectContentType(fileBytes)
	switch fileType {
	case "application/csv":
		break
	case "text/csv":
		break
	default:
		renderError(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	recordsChan := make(chan model.Record)

	go parse.ParseFile(file, recordsChan)

	for record := range recordsChan {
		go uH.uploadRecordToPersonsServer(record)
	}

	w.Write([]byte("SUCCESS"))
}

func (uH *uploadHandler) uploadRecordToPersonsServer(record model.Record) {
	jsonStr, err := json.Marshal(record)
	if err != nil {
		fmt.Printf("Error during converting record into strign: %s", err)
	}
	req, err := http.NewRequest("POST", uH.personsServerBaseUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error during uploading record: %s", err)
	}
	defer resp.Body.Close()
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}
