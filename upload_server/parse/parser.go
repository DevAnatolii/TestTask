package parse

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"testTask/upload_server/model"
)

func ParseFile(r io.Reader, records chan<- *model.Record) {
	reader := csv.NewReader(r)
	skipHeaders := true
	for {
		// read just one record
		rawRecord, err := reader.Read()

		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			break
		}

		//skip column headers, which are listed in first row
		if skipHeaders {
			skipHeaders = false
			continue
		}

		id, err := strconv.Atoi(rawRecord[0])
		if err != nil {
			close(records)
		}

		record := model.Record{
			Id:    id,
			Name:  rawRecord[1],
			Email: rawRecord[2],
			Phone: rawRecord[3],
		}

		records <- &record
	}
	close(records)
}
