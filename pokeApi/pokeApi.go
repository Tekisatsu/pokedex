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
	Pokedex map[string][]byte
	Info Pokemon
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
type Pokemon struct {
	Name string `json:"name"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
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
func Inspect(context *CliContext) error {
	pokemonName := context.Args[0]
	if len(pokemonName) == 0 {
		return errors.New("no pokemon provided")
	}
	var d Pokemon
	val,ok := context.Pokedex[pokemonName]
	if !ok {
		return errors.New("pokemon not caught")
	}
	err := json.Unmarshal(val,&d)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\n",d.Name,d.Height,d.Weight)
	for _,s := range d.Stats {
		fmt.Printf("%v: %v\n",s.Stat.Name,s.BaseStat)
	}
	for _,t:= range d.Types {
		fmt.Println(t.Type.Name)
	}
	return nil
}
func Catch(context *CliContext) error {
	if len(context.Args) == 0 {
		return errors.New("No pokemon provided")
	}
	pokemonName := context.Args[0]
	url := "https://pokeapi.co/api/v2/pokemon/"+pokemonName+"/"
	_, ok := context.Cache.get(url)
	if ok {
		context.Pokedex[pokemonName] = context.Cache.Data[url].val
		fmt.Printf("%v cached pokemon caught\n",pokemonName)
		return nil
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
	context.Pokedex[pokemonName] = body
	fmt.Printf("%v pokemon caught\n",pokemonName)
	return nil
}
func Encounter(context *CliContext) error {
	if len(context.Args) == 0 {
		return errors.New("Area not provided")
	}
	areaName := context.Args[0]
	url := "https://pokeapi.co/api/v2/location-area/"+areaName+"/"
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
