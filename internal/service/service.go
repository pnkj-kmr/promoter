package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/pnkj-kmr/promoter/internal/settings"
	"go.uber.org/zap"
)

type _service struct {
	l       *zap.Logger
	app     settings.App
	timeout time.Duration
	_check  bool
}

func New(s string) Service {
	app, ok := settings.APPS[s]
	if !ok {
		return &_service{}
	}
	timeout := settings.C.Timeout * time.Second
	if timeout == 0 {
		timeout = settings.LongTimeout
	}
	return &_service{settings.L, app, timeout, false}
}

func (s *_service) GetID() int32 {
	return s.app.AppId
}

func (s *_service) GetPriority() int32 {
	return s.app.Priority
}

func (s *_service) GetPersist() int32 {
	return s.app.Persist
}

func (s *_service) Start() error {
	s.l.Warn("application starting...", zap.Any("app_id", s.app.AppId))
	if s.app.Start == "" {
		s.l.Debug("app_does_not_exist - cannot start the app")
		return fmt.Errorf("app_does_not_exist")
	}
	return s._execute(s.app.Start, os.Stdout)
}

func (s *_service) Stop() error {
	s.l.Warn("application stopping...", zap.Any("app_id", s.app.AppId))
	if s.app.Stop == "" {
		s.l.Debug("app_does_not_exist - cannot stop the app")
		return fmt.Errorf("app_does_not_exist")
	}
	return s._execute(s.app.Stop, os.Stdout)
}

func (s *_service) Check() (err error) {
	s.l.Debug("application status...", zap.Any("app_id", s.app.AppId))
	if s.app.Status == "" {
		s.l.Warn("app_does_not_exist - cannot check status of app")
		return fmt.Errorf("app_does_not_exist")
	}
	s._check = false
	err = s._execute(s.app.Status, s)
	if s._check {
		s.l.Debug("application status... app is running", zap.Any("app_id", s.app.AppId), zap.Error(err))
		err = nil
	} else {
		s.l.Debug("application status... app is NOT running", zap.Any("app_id", s.app.AppId), zap.Error(err))
		err = fmt.Errorf("not_running")
	}
	return
}

func (s *_service) Write(p []byte) (int, error) {
	matchCheck := string(p)
	if strings.Contains(matchCheck, s.app.StatusMatch) {
		s.l.Debug("<<<<<<<<<<< STATUS MATCH FOUND >>>>>>>>>>")
		s._check = true
	}
	s.l.Info("=== application status ===", zap.Any("app_id", s.app.AppId), zap.Any("status", s._check), zap.Any("output", string(p)))
	return len(p), nil
}

func (s *_service) _execute(str string, stdout io.Writer) (err error) {
	if runtime.GOOS == "windows" {
		s.l.Debug("======= can't Execute this on a windows machine =======")
		err = fmt.Errorf("unable_to_execute_windows_system")
	} else {
		cmdArgs := strings.Split(str, "\n")
		for _, val := range cmdArgs {
			var cmds []string
			for _, x := range strings.Split(val, " ") {
				if x != "" {
					cmds = append(cmds, x)
				}
			}
			if len(cmds) > 0 {
				s.l.Debug("commands ------", zap.Any("cmds", cmds))
				ctx, _ := context.WithTimeout(context.Background(), s.timeout)
				cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
				cmd.Stdout = stdout
				if _err := cmd.Run(); _err != nil {
					s.l.Error("commands ERROR ------", zap.Any("cmds", cmds), zap.Error(_err))
					err = _err
				}
			}
		}
	}
	return
}
