package manager

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"wapiti/models"

	"github.com/Ullaakut/nmap/v2"
)

func (m *Manager) Nmap(ctx context.Context, ip string, ports []string) (map[string]string, error) {
	nmapResult := make(map[string]string, 0)
	for _, port := range ports {
		nmapResult[port] = ""
	}

	scanner, err := nmap.NewScanner(
		nmap.WithContext(ctx),
		nmap.WithBinaryPath("/usr/bin/nmap"),
		nmap.WithTargets(ip),
		nmap.WithPorts(ports...),
		nmap.WithUDPScan(),
		nmap.WithSYNScan(),
	)
	if err != nil {
		return nmapResult, err
	}

	result, _, err := scanner.Run()
	if err != nil {
		return nmapResult, err
	}

	b, _ := json.Marshal(result.Hosts)
	log.Error("++++++++++++++++++ string(b):", string(b))

	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			if port.Status() == nmap.Open {
				key := strconv.Itoa(int(port.ID))
				if strings.HasPrefix(port.Service.Name, models.ProtocolHttp) {
					nmapResult[key] = models.ProtocolHttp
				} else {
					nmapResult[key] = port.Protocol
				}
			}
		}
	}

	return nmapResult, nil
}
