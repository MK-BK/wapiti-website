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

// 在规定的时间内完成，或者超时失败
func (m *Manager) Handler() {
	start := time.Now()
	log.Infof("Job: %+v start", m.Uid)

	ctx := context.Background()

	defer log.Infof("Job: %+v end, spend: %+v", m.Uid, time.Now().Sub(start))

	var wg sync.WaitGroup
	for ip, config := range m.Request.Config {
		wg.Add(1)
		go func(ctx context.Context, ip string, config models.Config) {
			defer wg.Done()
			if err := m.Nmap(ctx, ip, config.Ports); err != nil {
				log.Error(err)
			}

			if config.IsAttack {
				if err := m.Wapiti(ip, config.Ports); err != nil {
					log.Error(err)
				}
			}
		}(ctx, ip, config)
	}
	wg.Wait()
}
