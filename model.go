package main

// swagger:model Config
type Config struct {
	// Id of the config
	// in: string
	Id string `json:"id"`

	// Map of entries
	// in: map[string]string
	Entries map[string]string `json:"entries"`
}

// swagger:model Group
type Group struct {
	// Id of the group
	// in: string
	Id string `json:"id"`

	// List of configurations
	// in: []Config
	Configs []Config `json:"configs"`
}
