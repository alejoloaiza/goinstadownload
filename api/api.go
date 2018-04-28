package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Name struct {
	Scale     float64 `json:"scale"`
	Gender    string  `json:"gender"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	ID        string  `json:"id"`
}

func GetGender(name string) string {

	url := fmt.Sprintf("https://api.namsor.com/onomastics/api/json/gender/%s", name)
	fmt.Println(url)
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
	/*bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)*/
	return record.Gender
}
