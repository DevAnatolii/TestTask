# Golang test task

There are 2 microservices.
First service uploads csv file and parse it (It doesn't matter how huge the uploading file is. I can be 200 KB or 5GB). It's launched on localhost on 8080 port.
These parameters can be changes in main.go file in root directory.
For uploading file, you need to use "/upload" endpoint, which accepts "uploadFile" parameter
Second service saves the records in database. It's launched on localhost on 8000 port.
These parameters can be changes in main.go file in root directory as well.
For uploading a single record to this service, you need to use "/persons" endpoint and put the record as an json text in body of request.
