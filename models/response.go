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
	StatusRunning  = "Running"
	StatusComplete = "Complete"
	StatusFailed   = "Failed"
)

func NewResponses() *Responses {
	return &Responses{
		Nmap:   make(map[string]map[string]string, 0),
		Kapiti: make(map[string][]*KapitiResponse, 0),
		Status: StatusRunning,
	}
}

type Responses struct {
	Nmap   map[string]map[string]string
	Kapiti map[string][]*KapitiResponse
	Status string
	Error  string
}

type NmapResponse map[int]string

type KapitiResponse struct {
	Port        string
	ScanResults ScanResult
}

type ScanResult struct {
	Vulnerabilities json.RawMessage `json:"vulnerabilities"`
}
