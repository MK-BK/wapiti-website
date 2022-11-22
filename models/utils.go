package models

import "os/exec"

type RequestCmd struct {
	Port     string
	Cmd      *exec.Cmd
	FilePath string
}
