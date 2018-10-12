package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"testTask/microservice_comunication"
	"testTask/upload_server/model"
	"testTask/upload_server/parse"

	"google.golang.org/grpc"
)

const (
	HandlePath           = "/upload"
	UploadFileParameter  = "uploadFile"
	coroutineCount       = 10
	detectTypeByteLength = 512 // check file type, detectcontenttype only needs the first 512 bytes
)

type uploadHandler struct {
	personsServerBaseUrl string
	errorLogFilePath     string
}

func NewUploadHandler(personsServerBaseUrl, errorLogFilePath string) *uploadHandler {
	return &uploadHandler{personsServerBaseUrl, errorLogFilePath}
}

func (uH *uploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, "Only Post method is allowed", http.StatusMethodNotAllowed)
		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		fmt.Printf("Error ocured during receiving MultipartReader: %s", err)
		return
	}
	defer r.Body.Close()

	// Read multipart body until the "uploadFile" part
	uploadFile, err := obtainUploadFile(mr)
	if err != nil {
		renderError(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileBytes := make([]byte, detectTypeByteLength, detectTypeByteLength)
	uploadFile.Read(fileBytes)
	if !testFileType(fileBytes) {
		renderError(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// need to append already read bytes in order to receive entire file
	joinedReader := io.MultiReader(bytes.NewReader(fileBytes), uploadFile)
	recordsChan := make(chan *model.Record, coroutineCount)

	errorsChan := make(chan *model.ErrorRecord, coroutineCount)
	wg := sync.WaitGroup{}
	wg.Add(coroutineCount)

	go parse.ParseFile(joinedReader, recordsChan)
	for i := 0; i < coroutineCount; i++ {
		go streamRecords(uH.personsServerBaseUrl, recordsChan, errorsChan, &wg)
	}

	go logErrors(errorsChan, uH.errorLogFilePath)

	wg.Wait()
	close(errorsChan)
	w.Write([]byte("SUCCESS"))
}

func obtainUploadFile(mr *multipart.Reader) (io.Reader, error) {
	// Read multipart body until the "uploadFile" part
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		if part.FormName() == UploadFileParameter {
			return part, nil
		}
	}
	return nil, errors.New("uploading file not found")
}

func testFileType(fileBytes []byte) bool {
	fileType := http.DetectContentType(fileBytes)
	return fileType == "text/csv" || fileType == "application/csv" || fileType == "text/plain; charset=utf-8"
}

func streamRecords(url string, records chan *model.Record, errorChan chan *model.ErrorRecord, wg *sync.WaitGroup) {
	defer wg.Done()

	connection, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	defer connection.Close()

	client := contract.NewPersonsClient(connection)
	stream, err := client.AddRecord(context.Background())

	for record := range records {
		stream.Send(&contract.AddPersonRequest{
			Id:    int32(record.Id),
			Name:  record.Name,
			Email: record.Email,
			Phone: record.Phone,
		})

		response, err := stream.Recv()
		if err != nil {
			log.Fatalf("failed to parse: %s", err)
		}
		if !response.Processed {
			errorChan <- &model.ErrorRecord{
				Record: record,
				Error:  errors.New(response.ErrorMessage),
			}
		}
	}
	stream.CloseSend()
}

func logErrors(errorRecords chan *model.ErrorRecord, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error during logging errors: %s", err)
		return
	}
	defer file.Close()

	for errorRecord := range errorRecords {
		io.WriteString(file, errorRecord.String())
		io.WriteString(file, "\n")
	}
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
