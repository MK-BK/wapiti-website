package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"wapiti/models"
)

func (m *Manager) Wapiti(ctx context.Context, ip string, ports []string) ([]*models.KapitiResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			cmds := combineCmd(ip, ports, m.Request.Modules)
			results := make([]*models.KapitiResponse, 0)

			for _, cmd := range cmds {
				cmd.Cmd.Stderr = os.Stderr

				if err := cmd.Cmd.Run(); err != nil {
					log.Error(err)
					continue
				}

				b, err := ioutil.ReadFile(cmd.FilePath)
				if err != nil {
					log.Error(err)
					continue
				}

				var scanResult models.ScanResult
				if err := json.Unmarshal(b, &scanResult); err != nil {
					log.Error(err)
					continue
				}

				results = append(results, &models.KapitiResponse{
					Port:        cmd.Port,
					ScanResults: scanResult,
				})
			}

			return results, nil
		}
	}
}

func combineCmd(ip string, ports []string, modules string) []*models.RequestCmd {
	cmds := make([]*models.RequestCmd, 0)

	for _, port := range ports {
		url := fmt.Sprintf("https://%+v/", ip)
		if port != "" {
			url = fmt.Sprintf("https://%+v:%+v/", ip, port)
		}

		tmp, _ := ioutil.TempFile(os.TempDir(), "*.json")

		args := make([]string, 0)
		args = append(args, "-u", url)
		args = append(args, "-o", tmp.Name())
		args = append(args, "-f", "json")

		if modules != "" {
			args = append(args, "-m", modules)
		}

		cmds = append(cmds, &models.RequestCmd{
			Cmd:      exec.Command("wapiti", args...),
			FilePath: tmp.Name(),
			Port:     port,
		})

	}

	return cmds
}
