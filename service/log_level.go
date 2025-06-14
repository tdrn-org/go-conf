// Copyright (C) 2024-2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package service

import (
	"log/slog"
	"reflect"

	"github.com/tdrn-org/go-conf"
)

// LogLevelService is used to provide access to [slog.LevelVar] instance
// controlling the active log level.
type LogLevelService interface {
	conf.Service
	LevelVar() *slog.LevelVar
}

type logLevelService struct {
	level *slog.LevelVar
}

func (*logLevelService) Type() reflect.Type {
	return reflect.TypeFor[LogLevelService]()
}

func (s *logLevelService) LevelVar() *slog.LevelVar {
	return s.level
}

func init() {
	_ = conf.BindService[LogLevelService](&logLevelService{level: &slog.LevelVar{}})
}
