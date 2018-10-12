package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

const recordCount = 100000000

func main() {
	file, err := os.Create("generatedFile.csv")
	if err != nil {
		log.Fatalf("error during create file: %s", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	data := []string{"id", "name", "email", "mobile_number"}
	if err := writer.Write(data); err != nil {
		fmt.Println(fmt.Errorf("error during create file: %s", err))
		return
	}

	for i := 0; i < recordCount; i++ {
		data := []string{
			strconv.Itoa(i),
			fmt.Sprintf("name%d", i),
			fmt.Sprintf("email%d@gmail.com", i),
			fmt.Sprintf(fmt.Sprintf("%011d", i)),
		}

		if err := writer.Write(data); err != nil {
			log.Println(fmt.Errorf("error during create file: %s", err))
			return
		}
	}
}
