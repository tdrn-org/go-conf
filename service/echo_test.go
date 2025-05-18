// Copyright (C) 2024-2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package service_test

import (
	"testing"

	"github.com/hdecarne-github/go-conf"
	"github.com/hdecarne-github/go-conf/service"
	"github.com/stretchr/testify/require"
)

func TestDefaultEchoService(t *testing.T) {
	defaultEchoService := service.DefaultEchoService()
	echoService := conf.LookupServiceOrDefault(defaultEchoService)
	require.Equal(t, defaultEchoService, echoService)
	echoService.Echo("TestDefaultEchoService\n")
}
