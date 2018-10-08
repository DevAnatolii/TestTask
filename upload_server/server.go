package upload_server

import (
	"log"
	"net/http"
	"testTask/upload_server/handler"
)

func Start(address string, personsServerBaseUrl string) {
	serveMux := http.NewServeMux()
	serveMux.Handle(handler.HandlePath, handler.NewUploadHandler(personsServerBaseUrl))
	log.Fatal(http.ListenAndServe(address, serveMux))
}
