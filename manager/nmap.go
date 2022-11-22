package manager

import (
	"context"
	"strings"
	"wapiti/models"

	"github.com/Ullaakut/nmap/v2"
)

func (m *Manager) Nmap(ctx context.Context, ip string, ports []string) error {
	if m.Response.Nmap[ip] == nil {
		m.Response.Nmap[ip] = make([]*models.NmapResponse, 0)
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
		return err
	}

	result, _, err := scanner.Run()
	if err != nil {
		return err
	}

	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			if port.Status() == nmap.Open {
				response := &models.NmapResponse{
					Port: port.ID,
				}
				if strings.HasPrefix(port.Service.Name, models.ProtocolHttp) {
					response.Protocol = models.ProtocolHttp
				} else {
					response.Protocol = port.Protocol
				}

				m.Response.Nmap[ip] = append(m.Response.Nmap[ip], response)
			}
		}
	}

	return nil
}
