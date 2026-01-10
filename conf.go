//
// Copyright (C) 2024-2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

// Package conf provides functions to bind and lookup configuration objects and service interfaces dynamically.
package conf

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Typed provides type information for dynamically bound instances.
type Typed interface {
	// Type gets the [reflect.Type] this instance represents.
	Type() reflect.Type
}

func ensureIsAssignableTo(r, l reflect.Type) {
	if !r.AssignableTo(l) {
		panic(fmt.Sprintf("type mismatch %s not assignable to %s", r, l))
	}
}

// Resolve casts the given [Typed] instance to the type it represents.
// This function panics if the cast is not valid.
func Resolve[T Typed](typed Typed) T {
	ensureIsAssignableTo(typed.Type(), reflect.TypeFor[T]())
	return typed.(T)
}

// Configuration represents a bindable configuration instances.
type Configuration interface {
	Typed
	// Bind binds this instance to its type via [BindConfiguration].
	Bind()
}

type configurationTableEntry struct {
	configuration Configuration
	applys        []func(Configuration)
}

var configurationTable map[reflect.Type]*configurationTableEntry = make(map[reflect.Type]*configurationTableEntry)
var configurationTableMutex sync.RWMutex = sync.RWMutex{}

// BindConfiguration binds a [Configuration] instance to its type.
// All previously via [ApplyConfiguration] registered apply functions are invoked during
// this bind operation.
func BindConfiguration[C Configuration](configuration C) {
	configurationTableMutex.Lock()
	defer configurationTableMutex.Unlock()
	configurationType := configuration.Type()
	ensureIsAssignableTo(reflect.TypeOf(configuration), configurationType)
	entry, ok := configurationTable[configurationType]
	if ok {
		entry.configuration = configuration
		for _, apply := range entry.applys {
			apply(configuration)
		}
	} else {
		entry = &configurationTableEntry{
			configuration: configuration,
			applys:        make([]func(Configuration), 0),
		}
		configurationTable[configurationType] = entry
	}
}

// BindToConfiguration binds the given apply function to the given type.
// The apply function is invoked everytime a [Configuration] instance of
// the given type is bound via [BindConfiguration].
func BindToConfiguration(configurationType reflect.Type, apply func(Configuration)) {
	configurationTableMutex.Lock()
	defer configurationTableMutex.Unlock()
	entry, ok := configurationTable[configurationType]
	if ok {
		entry.applys = append(entry.applys, apply)
		apply(entry.configuration)
	} else {
		entry = &configurationTableEntry{
			applys: []func(Configuration){apply},
		}
		configurationTable[configurationType] = entry
	}
}

// LookupConfiguration looks up the [Configuration] instance for the given type.
// The 2nd bool return value indicates whether the requested configuration type
// has already been bound.
func LookupConfiguration[C Configuration]() (C, bool) {
	configurationTableMutex.RLock()
	defer configurationTableMutex.RUnlock()
	entry, ok := configurationTable[reflect.TypeFor[C]()]
	if ok {
		return Resolve[C](entry.configuration), true
	}
	var none C
	return none, false
}

// LookupConfigurationOrDefault looks up the [Configuration] instance for the given type.
// If the requested configuration type has not yet been bound, the given default configuration
// is returned.
func LookupConfigurationOrDefault[C Configuration](defaultConfiguration C) C {
	configuration, ok := LookupConfiguration[C]()
	if ok {
		return configuration
	}
	return defaultConfiguration
}

// Service represents bindable service instances.
type Service interface {
	Typed
}

var serviceTable map[reflect.Type]Service = make(map[reflect.Type]Service)
var serviceTableMutex sync.RWMutex = sync.RWMutex{}

// ErrServiceAlreadyBound indicates a service is already bound.
var ErrServiceAlreadyBound = errors.New("service already bound")

// BindService binds a [Service] instance to its type.
// This functions returns [ErrServiceAlreadyBound] if the submitted service type has already been bound.
// This function panics if the [Service] instance's type does not match the service type.
func BindService[S Service](service S) error {
	serviceTableMutex.Lock()
	defer serviceTableMutex.Unlock()
	serviceType := service.Type()
	ensureIsAssignableTo(reflect.TypeOf(service), serviceType)
	_, ok := serviceTable[serviceType]
	if ok {
		return fmt.Errorf("%w: %s", ErrServiceAlreadyBound, serviceType)
	}
	serviceTable[serviceType] = service
	return nil
}

// LookupService looks up the [Service] instance for the given type.
// The 2nd bool return value indicates whether the requested service type
// has already been bound.
func LookupService[S Service]() (S, bool) {
	serviceTableMutex.RLock()
	defer serviceTableMutex.RUnlock()
	service, ok := serviceTable[reflect.TypeFor[S]()]
	if ok {
		return Resolve[S](service), true
	}
	var none S
	return none, false
}

// LookupServiceOrDefault looks up the [Service] instance for the given type.
// If the requested service type has not yet been bound, the given default service
// instance is returned.
func LookupServiceOrDefault[S Service](defaultService S) S {
	service, ok := LookupService[S]()
	if ok {
		return Resolve[S](service)
	}
	return defaultService
}
