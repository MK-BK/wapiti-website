package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"wapiti/models"
)

func (m *Manager) Wapiti(ip string, ports []string) error {
	if m.Response.Kapiti[ip] == nil {
		m.Response.Kapiti[ip] = make([]*models.KapitiResponse, 0)
	}

	requestCmds, _ := combineCmd(ip, ports, m.Request.Modules)

	for _, cmd := range requestCmds {
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

		m.Response.Kapiti[ip] = append(m.Response.Kapiti[ip], &models.KapitiResponse{
			Port:        cmd.Port,
			ScanResults: scanResult,
		})
	}

	return nil
}

func combineCmd(ip string, ports []string, modules string) ([]*models.RequestCmd, []string) {
	requestCmds := make([]*models.RequestCmd, 0)
	cmds := make([]string, 0)

	for _, port := range ports {
		url := fmt.Sprintf("http://%+v/", ip)
		if port != "" {
			url = fmt.Sprintf("http://%+v:%+v/", ip, port)
		}

		tmp, _ := ioutil.TempFile(os.TempDir(), "*.json")

		args := make([]string, 0)
		args = append(args, "-u", url)
		args = append(args, "-o", tmp.Name())
		args = append(args, "-f", "json")
		args = append(args, "--skip-crawl")

		if modules != "" {
			args = append(args, "-m", modules)
		}

		requestCmds = append(requestCmds, &models.RequestCmd{
			Cmd:      exec.Command("wapiti", args...),
			FilePath: tmp.Name(),
			Port:     port,
		})

		cmds = append(cmds, exec.Command("wapiti", args...).String())
	}

	return requestCmds, cmds
}
