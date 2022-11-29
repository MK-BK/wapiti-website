package models

type Request struct {
	Address  map[string][]string `json:"address"`
	IsAttack bool                `json:"is_attack"`
	Modules  string              `json:"modules"`
	Timeout  TimeoutOption       `json:"timeout"`
}

type TimeoutOption struct {
	Protocol string `json:"protocol"`
	Attack   string `json:"attack"`
}
