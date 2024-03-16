package main

import (
	"context"
	"errors"
	"os"
	"sandboxer/pkg/config"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"strconv"
)

type UninstallStageDelete struct {
	name string
	path string
}

var _ UninstallStage = &UninstallStageDelete{}

func NewUninstallStageDelete(name, path string) *UninstallStageDelete {
	return &UninstallStageDelete{
		name: name,
		path: path,
	}
}

func (u *UninstallStageDelete) Name() string {
	return u.name
}

func (u *UninstallStageDelete) Execute() error {
	return os.RemoveAll(u.path)
}

type UninstallStageStopProcess struct {
}

var _ UninstallStage = &UninstallStageStopProcess{}

func NewUninstallStageStopProcess() *UninstallStageStopProcess {
	return &UninstallStageStopProcess{}
}

func (*UninstallStageStopProcess) Name() string {
	return "Stop Program"
}

func (*UninstallStageStopProcess) Execute() error {
	pidFilePath, err := globals.PidFilePath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(pidFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.Errorf("Stop "+globals.AppName+": %s: %v", pidFilePath, err)
			return nil
		}
		return err
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		logging.Errorf("Stop"+globals.AppName+": %s: %v", string(data), err)
		return nil
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		logging.Errorf("Stop"+globals.AppName+": FindProcess(%d): %v", pid, err)
		return nil
	}
	if err := proc.Kill(); err != nil {
		logging.Errorf("Stop"+globals.AppName+": Kill %d: %v", pid, err)
		return nil
	}
	return nil
}

type UninstallStageUnregister struct {
	conf *config.DDAn
}

var _ UninstallStage = &UninstallStageUnregister{}

func NewUninstallStageUnregister(conf *config.DDAn) *UninstallStageUnregister {
	return &UninstallStageUnregister{
		conf: conf,
	}
}

func (*UninstallStageUnregister) Name() string {
	return "Analyzer Unregister"
}

func (u *UninstallStageUnregister) Execute() error {
	if err := u.conf.LoadClientUUID(); err != nil {
		logging.LogError(err)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	analyzer, err := u.conf.Analyzer()
	if err != nil {
		return err
	}
	return analyzer.Unregister(context.TODO())
}

type UninstallStageDone struct{}

var _ UninstallStage = &UninstallStageDone{}

func NewUninstallStageDone() *UninstallStageDone {
	return &UninstallStageDone{}
}

func (*UninstallStageDone) Name() string {
	return "Done"
}

func (*UninstallStageDone) Execute() error {
	return nil
}
