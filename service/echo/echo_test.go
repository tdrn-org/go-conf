//
// Copyright (C) 2024-2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package echo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-conf"
	"github.com/tdrn-org/go-conf/service/echo"
)

func TestDefaultEchoService(t *testing.T) {
	defaultEchoService := echo.DefaultEchoService()
	echoService := conf.LookupServiceOrDefault(defaultEchoService)
	require.Equal(t, defaultEchoService, echoService)
}

func TestEchoing(t *testing.T) {
	echo.Out("out message\n")
	echo.Err("err message\n")
}
