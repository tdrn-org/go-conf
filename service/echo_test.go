// Copyright (C) 2024-2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package service_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-conf"
	"github.com/tdrn-org/go-conf/service"
)

func TestDefaultEchoService(t *testing.T) {
	defaultEchoService := service.DefaultEchoService()
	echoService := conf.LookupServiceOrDefault(defaultEchoService)
	require.Equal(t, defaultEchoService, echoService)
	echoService.Echo("TestDefaultEchoService\n")
}
