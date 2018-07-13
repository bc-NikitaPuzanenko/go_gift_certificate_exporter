package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Data struct
type data []struct {
	ID           int    `json:"id"`
	CustomerID   string `json:"customer_id"`
	OrderID      string `json:"order_id"`
	Code         string `json:"code"`
	ToName       string `json:"to_name"`
	ToEmail      string `json:"to_email"`
	FromName     string `json:"from_name"`
	FromEmail    string `json:"from_email"`
	Amount       string `json:"amount"`
	Balance      string `json:"balance"`
	Status       string `json:"status"`
	Template     string `json:"template"`
	Message      string `json:"message"`
	PurchaseDate string `json:"purchase_date"`
	ExpiryDate   string `json:"expiry_date"`
}

// Config file struct
type Config struct {
	Host        string `json:"host"`
	Credentials string `json:"credentials"`
}

// LoadConfiguration and parse JSON config file
func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func main() {
	config := LoadConfiguration("./config.json")
	keepLooping := true
	page := 1
	var allJSON data

	for keepLooping {
		fullURL := strings.Join([]string{"https://", config.Host, "/api/v2/gift_certificates?limit=250&page=", strconv.Itoa(page)}, "")
		client := &http.Client{}
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Accept", `application/json`)
		req.Header.Set("Content-Type", `application/json`)
		req.Header.Set("Content-Type", `application/json`)
		req.Header.Set("Authorization", config.Credentials)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request error: ", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Read error: ", err)
		}

		var tmpRecords data

		json.Unmarshal(body, &tmpRecords)

		allJSON = append(allJSON, tmpRecords...)

		if resp.StatusCode == 200 {

			info := strings.Join([]string{"HTTP: ", resp.Status, " | Page: ", strconv.Itoa(page)}, "")
			fmt.Println(info)
			page++
		}

		if resp.StatusCode == 204 {

			file, err := os.OpenFile(config.Host+".csv", os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal(err)
			}

			defer file.Close()

			writer := csv.NewWriter(file)

			var newHeaders = []string{"id", "customer id", "order id", "code", "to name", "to email", "from name", "from email", "amount", "balance", "status", "template", "message", "purchase date", "expiry date"}
			writer.Write(newHeaders)

			for _, value := range allJSON {
				var record []string
				record = append(record, strconv.Itoa(value.ID))
				record = append(record, value.CustomerID)
				record = append(record, value.OrderID)
				record = append(record, value.Code)
				record = append(record, value.ToName)
				record = append(record, value.ToEmail)
				record = append(record, value.FromName)
				record = append(record, value.FromEmail)
				record = append(record, value.Amount)
				record = append(record, value.Balance)
				record = append(record, value.Status)
				record = append(record, value.Template)
				record = append(record, value.Message)
				record = append(record, value.PurchaseDate)
				record = append(record, value.ExpiryDate)
				writer.Write(record)
			}

			defer writer.Flush()

			fmt.Println("Done: No more pages")

			keepLooping = false
		}
	}
}
