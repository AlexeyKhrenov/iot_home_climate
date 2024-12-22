package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func postClimateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp, err := strconv.Atoi(r.FormValue("temp"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	humidity, err := strconv.Atoi(r.FormValue("humidity"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timestamp := time.Now()

	err = writeData("./sample_climate.txt", timestamp, temp, humidity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getClimateHandler(w http.ResponseWriter, r *http.Request) {
	/*
		_, err := time.Parse(time.RFC3339, r.URL.Query().Get("start"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, err = time.Parse(time.RFC3339, r.URL.Query().Get("end"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	*/

	data, err := readData("./sample_climate.txt")
	fmt.Println(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func readData(filePath string) ([]Measurement, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file with data : %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read .csv file : %w", err)
	}

	var result []Measurement
	for i, record := range records {
		if len(record) == 0 {
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			return nil, fmt.Errorf("Failed to pares .csv file. Timestamp. Error : %w. Row : %d", err, i+1)
		}
		temp, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("Failed to pares .csv file. Temperature. Error : %w. Row : %d", err, i+1)
		}
		humidity, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("Failed to pares .csv file. Humidity. Error : %w. Row : %d", err, i+1)
		}

		result = append(result, Measurement{Timestamp: timestamp, Temp: temp, Humidity: humidity})
	}

	return result, nil
}

func writeData(filePath string, timestamp time.Time, temp int, humidity int) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open the file %w", err, filePath)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{timestamp.Format(time.RFC3339), strconv.Itoa(temp), strconv.Itoa(humidity)}
	fmt.Println(record)
	if err = writer.Write(record); err != nil {
		return fmt.Errorf("Failed to write measurement to file %w", err, filePath)
	}

	return nil
}

func main() {
	http.HandleFunc("GET /climate", getClimateHandler)
	http.HandleFunc("POST /climate", postClimateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Measurement struct {
	Timestamp time.Time
	Temp      int // in celcius
	Humidity  int // %%
}
