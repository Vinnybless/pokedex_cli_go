package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Vinnybless/pokedex_cli_go/internal/pokecache"
)

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)
	return words
}

type Config struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	return nil
}

func commandMap(c *Config) error {
	url := "https://pokeapi.co/api/v2/location-area"

	if c.Results != nil {
		url = c.Next
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var client http.Client

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(c)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for i := 0; i < len(c.Results); i++ {
		fmt.Println(c.Results[i].Name)
	}

	return nil
}

func commandMapb(c *Config) error {
	if c.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	url, ok := c.Previous.(string)
	if !ok {
		return fmt.Errorf("c.Previous not a string")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var client http.Client

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(c)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for i := 0; i < len(c.Results); i++ {
		fmt.Println(c.Results[i].Name)
	}

	return nil
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Displays 20 location names",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Displays previous 20 location names",
		callback:    commandMapb,
	},
}

func main() {
	var config Config
	pc := pokecache.NewCache(5 * time.Second) // COME BACK HERE

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		txt := scanner.Text()
		if txt == "" {
			fmt.Println("Unknown command")
			continue
		}
		text := cleanInput(txt)[0]

		t, ok := commands[text]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := t.callback(&config)
		if err != nil {
			fmt.Println(err)
		}
		if text == "help" {
			for _, v := range commands {
				str := fmt.Sprintf("%s: %s", v.name, v.description)
				fmt.Println(str)
			}
		}
	}
}
