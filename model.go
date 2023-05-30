package main

type Config struct {
	Id      string            `json:"id"`
	Entries map[string]string `json:"entries"`
	Version string            `json:"version"`
}

type Group struct {
	Id      string   `json:"id"`
	Configs []Config `json:"configs"`
}
