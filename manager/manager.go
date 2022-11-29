package manager

import (
	"context"
	"sync"
	"time"
	"wapiti/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Manager struct {
	Uid      string
	Request  *models.Request
	Response *models.Responses
	Err      error
}

func NewManager(request *models.Request) *Manager {
	return &Manager{
		Uid:      uuid.New().String(),
		Request:  request,
		Response: models.NewResponses(),
	}
}

func (m *Manager) Handler(ctx context.Context) error {
	attackTimeout, err := time.ParseDuration(m.Request.Timeout.Attack)
	if err != nil {
		log.Error(err)
	}

	protocolTimeout, errr := time.ParseDuration(m.Request.Timeout.Protocol)
	if err != nil {
		log.Error(errr)
	}

	var wg sync.WaitGroup

	for ip, ports := range m.Request.Address {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, protocolTimeout)
			defer cancel()

			result, err := m.Nmap(ctx, ip, ports)
			if err != nil {
				log.Error(err)
				return
			}

			m.Response.Nmap[ip] = result
		}()

		if m.Request.IsAttack {
			wg.Add(1)

			go func() {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(ctx, attackTimeout)
				defer cancel()

				result, err := m.Wapiti(ctx, ip, ports)
				if err != nil {
					log.Error(err)
					return
				}

				m.Response.Kapiti[ip] = result
			}()
		}
	}

	wg.Wait()

	return nil
}
