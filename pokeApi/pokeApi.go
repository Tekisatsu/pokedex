package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type JsonConfig struct {
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

func GetMapUrl(state *JsonConfig) error {
	var url string
	if state.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	} else {
		url = state.Next
	}
	CommandMap(url, state)
	return nil
}
func GetPrevMapUrl(state *JsonConfig) error {
	var url string
	if state.Previous == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	} else {
		url = state.Previous
	}
	CommandMap(url, state)
	return nil
}
func CommandMap(url string, state *JsonConfig) error {

	res, err := http.Get(url)
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
	var d JsonConfig
	errD := json.Unmarshal(body, &d)
	if errD != nil {
		log.Fatal(errD)
	}
	for _, item := range d.Result {
		fmt.Println(item.Name)
	}
	state.Next, state.Previous = d.Next, d.Previous
	fmt.Printf("previous: %v\nnext: %v\n", d.Previous, d.Next)
	return nil
}
