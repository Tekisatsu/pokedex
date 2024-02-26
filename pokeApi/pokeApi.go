package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"errors"
)

type CliContext struct {
	State *JsonConfig
	Cache *Cache
	Args []string
}
type JsonConfig struct {
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Location Location `json:"location"`
	Result   []Result `json:"results"`
	Encounter []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type Result struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type Cache struct {
	Data     map[string]cacheEntry
	Mutex    sync.Mutex
	interval time.Duration
}
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		Data:     make(map[string]cacheEntry),
		Mutex:    sync.Mutex{},
		interval: interval,
	}
	go cache.reaploop(cache.interval)
	return cache
}
func (c *Cache) add(key string, val []byte) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}
func (c *Cache) get(key string) ([]byte, bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	value, ok := c.Data[key]
	if !ok {
		return nil, false
	}
	return value.val, true
}
func (c *Cache) reaploop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		c.reapStaleEntries()
	}
}
func (c *Cache) reapStaleEntries() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for key, entry := range c.Data {
		if time.Since(entry.createdAt) > c.interval {
			delete(c.Data, key)
		}
	}
}
func Encounter(context *CliContext) error {
	if len(context.Args) == 0 {
		return errors.New("Area not provided")
	}
	areaName := context.Args[0]
	url := "https://pokeapi.co/api/v2/location-area/"+areaName+"/"
	fmt.Println(url)
	value, ok := context.Cache.get(url)
	if ok {
		errC := json.Unmarshal(value, context.State)
		if errC != nil {
			log.Fatal(errC)
		}
		fmt.Println(context.State)
	}
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
	context.Cache.add(url, body)
	var d JsonConfig
	errD := json.Unmarshal(body, &d)
	if errD != nil {
		log.Fatal(errD)
	}
	for _, item:= range d.Encounter {
		fmt.Println(item.Pokemon.Name)
	}
	return nil
}
func GetMapUrl(context *CliContext) error {
	var url string
	if context.State.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	} else {
		url = context.State.Next
	}
	CommandMap(url, context)
	return nil
}
func GetPrevMapUrl(context *CliContext) error {
	var url string
	if context.State.Previous == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	} else {
		url = context.State.Previous
	}
	CommandMap(url, context)
	return nil
}
func CommandMap(url string, context *CliContext) error {

	value, ok := context.Cache.get(url)
	if ok {
		errC := json.Unmarshal(value, context.State)
		if errC != nil {
			log.Fatal(errC)
		}
		fmt.Println(context.State)
	}
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
	context.Cache.add(url, body)
	var d JsonConfig
	errD := json.Unmarshal(body, &d)
	if errD != nil {
		log.Fatal(errD)
	}
	for _, item:= range d.Result {
		fmt.Println(item.Name)
	}
	context.State.Next, context.State.Previous = d.Next, d.Previous
	fmt.Printf("previous: %v\nnext: %v\n", d.Previous, d.Next)
	return nil
}
