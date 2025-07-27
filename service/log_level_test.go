//
// Copyright (C) 2024-2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package service_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-conf"
	"github.com/tdrn-org/go-conf/service"
)

func TestLogLevelService(t *testing.T) {
	logLevel, ok := conf.LookupService[service.LogLevelService]()
	require.True(t, ok)
	require.NotNil(t, logLevel)
	// Check default
	require.Equal(t, slog.LevelInfo, logLevel.LevelVar().Level())
	// Check set
	logLevel.LevelVar().Set(slog.LevelDebug)
	require.Equal(t, slog.LevelDebug, logLevel.LevelVar().Level())
	// Check value after re-lookup
	logLevel, ok = conf.LookupService[service.LogLevelService]()
	require.True(t, ok)
	require.NotNil(t, logLevel)
	require.Equal(t, slog.LevelDebug, logLevel.LevelVar().Level())
}
