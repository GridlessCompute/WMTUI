package config

type Config struct {
	Sites []Site `json:"sites"`
}

type Site struct {
	Name    string `json:"name"`
	IPRange string `json:"iprange"`
}
