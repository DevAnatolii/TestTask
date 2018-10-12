package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const chunkSize = 1024 * 1024

func main() {
	targetUrl := "http://localhost:8080/upload"
	filename := "generatedFile.csv"
	uploadLargeFile(targetUrl, filename, chunkSize)
}

func uploadLargeFile(uri, filePath string, chunkSize int) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error happened during opening file: %s", err)
	}
	defer file.Close()

	byteBuf := &bytes.Buffer{}
	mpWriter := multipart.NewWriter(byteBuf)
	mpWriter.CreateFormFile("uploadFile", file.Name())
	contentType := mpWriter.FormDataContentType()

	multipartChunk := make([]byte, byteBuf.Len())
	_, _ = byteBuf.Read(multipartChunk)

	//part: latest boundary
	//when multipart closed, latest boundary is added
	mpWriter.Close()
	lastBoundary := make([]byte, byteBuf.Len())
	_, err = byteBuf.Read(lastBoundary)
	if err != nil {
		log.Fatalf("Error happened during reading last boundary: %s", err)
	}
	//use pipe to pass request
	rd, wr := io.Pipe()
	defer rd.Close()

	go func() {
		defer wr.Close()

		//write multipart
		write(wr, multipartChunk)

		//write file
		buf := make([]byte, chunkSize)
		for {
			n, err := file.Read(buf)
			if err != nil {
				break
			}
			write(wr, buf[:n])
		}
		//write boundary
		write(wr, lastBoundary)
	}()

	request, _ := http.NewRequest("POST", uri, rd)
	request.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: time.Hour}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Recieved error: %s", err)
	} else {
		defer response.Body.Close()

		body := &bytes.Buffer{}
		_, _ = body.ReadFrom(response.Body)

		log.Printf("Status code = %d\n Header = %s, \n Body = %s",
			response.StatusCode, response.Header, body)
	}
}

func write(wr *io.PipeWriter, data []byte) {
	if _, err := wr.Write(data); err != nil {
		log.Fatalf("Error happened during reading last boundary: %s", err)
	}
}
