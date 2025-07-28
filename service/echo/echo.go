//
// Copyright (C) 2024-2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package echo

import (
	"fmt"
	"os"
	"reflect"

	"github.com/tdrn-org/go-conf"
)

// EchoService mimics echoing messages to Stdout and Stderr in a [fmt.Print] style manner.
type EchoService interface {
	conf.Service
	// Out echos a message to a Stdout like output using a [fmt.Print] style formating
	Out(a ...any)
	// Out echos a message to a Stderr like output using a [fmt.Print] style formating
	Err(a ...any)
}

type stdoutEchoService struct{}

func (*stdoutEchoService) Type() reflect.Type {
	return reflect.TypeFor[EchoService]()
}

func (*stdoutEchoService) Out(a ...any) {
	_, _ = fmt.Fprint(os.Stdout, a...)
}

func (*stdoutEchoService) Err(a ...any) {
	_, _ = fmt.Fprint(os.Stderr, a...)
}

var defaultEchoService *stdoutEchoService = &stdoutEchoService{}

// DefaultEchoService gets the default [EchoService] simply wrapping [fmt.Fprint].
func DefaultEchoService() EchoService {
	return defaultEchoService
}

// Out is a shorthand for invoking [EchoService.Out].
func Out(a ...any) {
	conf.LookupServiceOrDefault[EchoService](defaultEchoService).Out(a...)
}

// Err is a shorthand for invoking [EchoService.Err].
func Err(a ...any) {
	conf.LookupServiceOrDefault[EchoService](defaultEchoService).Err(a...)
}
