package main

// swagger:model Config
type Config struct {
	Id      string            `json:"id"`
	Entries map[string]string `json:"entries"`
}

// swagger:model Group
type Group struct {
	Id      string   `json:"id"`
	Configs []Config `json:"configs"`
}
