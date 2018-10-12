package upload_server

import (
	"log"
	"net/http"
	"testTask/upload_server/handler"
)

func Start(serverBaseUrl string, personsServerBaseUrl string, errorFileLoggingPath string) {
	serveMux := http.NewServeMux()
	serveMux.Handle(handler.HandlePath, handler.NewUploadHandler(personsServerBaseUrl, errorFileLoggingPath))
	log.Println("Start upload server on " + serverBaseUrl)
	log.Fatal(http.ListenAndServe(serverBaseUrl, serveMux))
}
