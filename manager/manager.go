package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"wapiti/models"

	"github.com/Ullaakut/nmap/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Manager struct {
	Request  *models.Request
	Response []*models.Response
	ctx      context.Context
	Uuid     uuid.UUID
}

func NewManager(request *models.Request) *Manager {
	return &Manager{
		Uuid:     uuid.New(),
		Request:  request,
		Response: make([]*models.Response, 0),
		ctx:      context.Background(),
	}
}

func (m *Manager) Collector() {
	var wg sync.WaitGroup

	for _, target := range m.Request.Targets {
		wg.Add(1)

		go func(target models.Target) {
			defer wg.Done()

			response := newResponse()
			var childWg sync.WaitGroup

			childWg.Add(2)

			go m.Wapiti(&childWg, target, response)
			go m.Nmap(&childWg, target, response)

			childWg.Wait()
			m.Response = append(m.Response, response)
		}(target)
	}

	wg.Wait()
}

func (m *Manager) Wapiti(wg *sync.WaitGroup, target models.Target, response *models.Response) {
	defer wg.Done()

	if m.Request.Format == "" {
		m.Request.Format = "json"
	}

	url := fmt.Sprintf("http://%+v:%+v/", target.IP, target.Port)

	path := filepath.Join(os.TempDir(), uuid.New().String())

	args := make([]string, 0)
	args = append(args, "-u", url)
	args = append(args, "-o", path)
	args = append(args, "-f", m.Request.Format)

	if m.Request.Modules != "" {
		args = append(args, "-m", m.Request.Modules)
	}

	cmd := exec.CommandContext(m.ctx, "wapiti", args...)

	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Error(err)
		response.Errors = append(response.Errors, err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
		response.Errors = append(response.Errors, err)
	}

	if err := json.Unmarshal(b, &response.ScanResults); err != nil {
		response.Errors = append(response.Errors, err)
	}
}

func (m *Manager) Nmap(wg *sync.WaitGroup, target models.Target, response *models.Response) {
	defer wg.Done()

	scanner, err := nmap.NewScanner(
		nmap.WithContext(m.ctx),
		nmap.WithBinaryPath("/usr/bin/nmap"),
		nmap.WithTargets(target.IP),
		nmap.WithPorts(target.Port),
		nmap.WithUDPScan(),
		nmap.WithSYNScan(),
	)
	if err != nil {
		response.Errors = append(response.Errors, err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		response.Errors = append(response.Errors, err)
	}

	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			if port.Status() == nmap.Open {
				if strings.HasPrefix(port.Service.Name, models.ProtocolHttp) {
					response.Protocol = models.ProtocolHttp
				} else {
					response.Protocol = port.Protocol
				}
			}
		}
	}
}

func newResponse() *models.Response {
	return &models.Response{
		Errors:   make([]error, 0),
		Protocol: models.ProtocolUdp,
	}
}
