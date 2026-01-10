//
// Copyright (C) 2024-2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package loglevel

import (
	"log/slog"
	"reflect"

	"github.com/tdrn-org/go-conf"
)

// LogLevelService is used to provide access to a single [slog.LevelVar]
// instance controlling the active log level.
type LogLevelService interface {
	conf.Service
	// LevelVar gets the [slog.LevelVar] instance used to control the active log level.
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

// LevelVar is shorthand for invoking LogLevelService.LevelVar.
func LevelVar() *slog.LevelVar {
	logLevel, _ := conf.LookupService[LogLevelService]()
	return logLevel.LevelVar()
}

func init() {
	_ = conf.BindService[LogLevelService](&logLevelService{level: &slog.LevelVar{}})
}
