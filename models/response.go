package models

import (
	"encoding/json"
)

const (
	ProtocolHttp = "http"
	ProtocolUdp  = "udp"
	ProtocolTcp  = "tcp"
)

const (
	StatusTimeout  = "Timeout"
	StatusRunning  = "Running"
	StatusComplete = "Complete"
	StatusFailed   = "Failed"
)

func NewResponses() *Responses {
	return &Responses{
		Nmap:   make(map[string][]*NmapResponse, 0),
		Kapiti: make(map[string][]*KapitiResponse, 0),
		Status: StatusRunning,
	}
}

type Responses struct {
	Nmap   map[string][]*NmapResponse
	Kapiti map[string][]*KapitiResponse
	Status string
	Error  string
}

type NmapResponse struct {
	Port     uint16
	Protocol string
}

type KapitiResponse struct {
	Port        string
	ScanResults ScanResult
}

type ScanResult struct {
	Vulnerabilities json.RawMessage `json:"vulnerabilities"`
	Infos           json.RawMessage `json:"infos"`
}
