package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Name struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

func GetGender(name string) string {

	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return ""
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return ""
	}

	defer resp.Body.Close()

	var record Name

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	//fmt.Printf(">> Name = %s  Gender = %s \n", record.Name, record.Gender)
	return record.Gender
}
