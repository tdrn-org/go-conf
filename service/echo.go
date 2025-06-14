// Copyright (C) 2024-2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package service

import (
	"fmt"
	"reflect"

	"github.com/tdrn-org/go-conf"
)

// EchoService defines a simple service interface to echo messages in a [fmt.Print] style manner.
type EchoService interface {
	conf.Service
	Echo(a ...any)
}

type stdoutEchoService struct{}

func (*stdoutEchoService) Type() reflect.Type {
	return reflect.TypeFor[EchoService]()
}

func (*stdoutEchoService) Echo(a ...any) {
	_, _ = fmt.Print(a...)
}

var defaultEchoService *stdoutEchoService = &stdoutEchoService{}

// DefaultEchoService gets the default [EchoService] simply wrapping [fmt.Print].
func DefaultEchoService() EchoService {
	return defaultEchoService
}
