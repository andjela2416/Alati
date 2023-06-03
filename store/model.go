package store

// swagger:model Config
type Config struct {
	// Id of the config
	// in: string
	Id string `json:"id"`

	// Map of entries
	// in: map[string]string
	Entries map[string]string `json:"entries"`

	// Labels of the config
	// in: string
	Labels string `json:"labels"`

	// Version of the config
	// in: string
	Version string `json:"version"`
}

// swagger:model Group
type Group struct {
	// Id of the group
	// in: string
	Id string `json:"id"`

	// List of configurations
	// in: []Config
	Configs []Config `json:"configs"`

	// Version of the group
	// in: string
	Version string `json:"version"`

	// Labels of the config
	// in: string
	Labels string `json:"labels"` //ne treba da ima labele
}
