package models

type Request struct {
	Config  map[string]Config `json:"config"`
	Modules string            `json:"modules"`
}

type Config struct {
	Ports    []string `json:"ports"`
	IsAttack bool     `json:"is_attack"`
}
