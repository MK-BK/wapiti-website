package models

import "encoding/json"

const (
	ProtocolHttp = "http"
	ProtocolUdp  = "udp"
	ProtocolTcp  = "tcp"
)

type Request struct {
	Targets []Target `json:"targets"`
	Modules string   `json:"modules"`
	Format  string   `json:"format"`
}

type Target struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type ScanResult struct {
	Vulnerabilities json.RawMessage `json:"vulnerabilities"`
	Infos           json.RawMessage `json:"infos"`
}

type Response struct {
	ScanResults ScanResult
	Protocol    string `json:"protocol"`
	Errors      []error
}
