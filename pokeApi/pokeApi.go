package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type jsonConfig struct {
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Location Location `json:"location"`
	Result   []Result `json:"results"`
}
type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type Result struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func CommandMap() error {
	res, err := http.Get("https://pokeapi.co/api/v2/location-area/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed. Status code: %d", res.StatusCode)
	}
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var d jsonConfig
	errD := json.Unmarshal(body, &d)
	if errD != nil {
		log.Fatal(errD)
	}
	for _, item := range d.Result {
		fmt.Println(item.Name)
	}
	fmt.Printf("previous: %v\nnext: %v\n", d.Previous, d.Next)
	return nil
}
