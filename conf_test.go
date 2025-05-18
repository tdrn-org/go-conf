// Copyright (C) 2024-2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package conf_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-conf"
)

type TestTyped struct{}

func (t *TestTyped) Type() reflect.Type {
	return reflect.TypeFor[*TestTyped]()
}

func TestTypedResolve(t *testing.T) {
	typed := &TestTyped{}
	require.NotPanics(t, func() {
		conf.Resolve[conf.Typed](typed)
	})
	require.NotPanics(t, func() {
		conf.Resolve[*TestTyped](typed)
	})
	require.Panics(t, func() {
		conf.Resolve[conf.Configuration](typed)
	})
}

type TestConfig struct {
	ConfigurationType reflect.Type `yaml:"-"`
	StringValue       string       `yaml:"string"`
}

func (c *TestConfig) Type() reflect.Type {
	return c.ConfigurationType
}

type TestLookupUnknownConfigurationConfig struct {
	TestConfig
}

func (c *TestLookupUnknownConfigurationConfig) Bind() {
	conf.BindConfiguration(c)
}

func TestLookupUnknownConfiguration(t *testing.T) {
	unknownConfiguration, ok := conf.LookupConfiguration[*TestLookupUnknownConfigurationConfig]()
	require.Nil(t, unknownConfiguration)
	require.False(t, ok)
}

type TestBindAndLookupConfigurationConfig struct {
	TestConfig
}

func (c *TestBindAndLookupConfigurationConfig) Bind() {
	conf.BindConfiguration(c)
}

func TestBindAndLookupConfiguration(t *testing.T) {
	c := &TestBindAndLookupConfigurationConfig{
		TestConfig{
			ConfigurationType: reflect.TypeFor[*TestBindAndLookupConfigurationConfig](),
			StringValue:       "TestBindAndLookupConfigurationConfig",
		},
	}
	c.Bind()
	configuration, ok := conf.LookupConfiguration[*TestBindAndLookupConfigurationConfig]()
	require.NotNil(t, configuration)
	require.True(t, ok)
	require.Equal(t, c, configuration)
}

type TestBindBeforeBindToConfigurationConfig struct {
	TestConfig
}

func (c *TestBindBeforeBindToConfigurationConfig) Bind() {
	conf.BindConfiguration(c)
}

func TestBindBeforeBindToConfiguration(t *testing.T) {
	c := &TestBindBeforeBindToConfigurationConfig{
		TestConfig{
			ConfigurationType: reflect.TypeFor[*TestBindBeforeBindToConfigurationConfig](),
			StringValue:       "TestBindBeforeBindToConfigurationConfig",
		},
	}
	c.Bind()
	cApplied := false
	conf.BindToConfiguration(reflect.TypeFor[*TestBindBeforeBindToConfigurationConfig](), func(configuration conf.Configuration) {
		resolvedConfiguration := conf.Resolve[*TestBindBeforeBindToConfigurationConfig](configuration)
		require.Equal(t, c, resolvedConfiguration)
		cApplied = true
	})
	require.True(t, cApplied)
}

type TestBindAfterBindToConfigurationConfig struct {
	TestConfig
}

func (c *TestBindAfterBindToConfigurationConfig) Bind() {
	conf.BindConfiguration(c)
}

func TestBindAfterBindToConfiguration(t *testing.T) {
	c := &TestBindAfterBindToConfigurationConfig{
		TestConfig{
			ConfigurationType: reflect.TypeFor[*TestBindAfterBindToConfigurationConfig](),
			StringValue:       "TestBindAfterBindToConfigurationConfig",
		},
	}
	cApplied := false
	conf.BindToConfiguration(reflect.TypeFor[*TestBindAfterBindToConfigurationConfig](), func(configuration conf.Configuration) {
		resolvedConfiguration := conf.Resolve[*TestBindAfterBindToConfigurationConfig](configuration)
		require.Equal(t, c, resolvedConfiguration)
		cApplied = true
	})
	require.False(t, cApplied)
	c.Bind()
	require.True(t, cApplied)
}

type TestLookupUnknownServiceInterface interface {
	conf.Service
	InvokeTestLookupUnknownService()
}

func TestLookupUnknownService(t *testing.T) {
	service, ok := conf.LookupService[TestLookupUnknownServiceInterface]()
	require.Nil(t, service)
	require.False(t, ok)
}

type TestDoubleBindServiceInterface interface {
	conf.Service
	InvokeTestDoubleBindService()
}

type testDoubleBindService struct{}

func (*testDoubleBindService) Type() reflect.Type {
	return reflect.TypeFor[TestDoubleBindServiceInterface]()
}

func (*testDoubleBindService) InvokeTestDoubleBindService() { /* do nothing */ }

func TestDoubleBindService(t *testing.T) {
	err := conf.BindService(&testDoubleBindService{})
	require.NoError(t, err)
	err = conf.BindService(&testDoubleBindService{})
	require.ErrorIs(t, err, conf.ErrServiceAlreadyBound)
}

type TestBindAndLookupServiceInterface interface {
	conf.Service
	InvokeTestBindAndLookupService()
}

type testBindAndLookupService struct{}

func (*testBindAndLookupService) Type() reflect.Type {
	return reflect.TypeFor[TestBindAndLookupServiceInterface]()
}

func (*testBindAndLookupService) InvokeTestBindAndLookupService() { /* do nothing */ }

func TestBindAndLookupService(t *testing.T) {
	err := conf.BindService(&testBindAndLookupService{})
	require.NoError(t, err)
	service, ok := conf.LookupService[TestBindAndLookupServiceInterface]()
	require.NotNil(t, service)
	require.True(t, ok)
	service.InvokeTestBindAndLookupService()
}

type TestLookupServiceOrDefaultInterface interface {
	conf.Service
	InvokeTestLookupServiceOrDefault()
}

type testLookupServiceOrDefaultService struct{}

func (*testLookupServiceOrDefaultService) Type() reflect.Type {
	return reflect.TypeFor[TestLookupServiceOrDefaultInterface]()
}

func (*testLookupServiceOrDefaultService) InvokeTestLookupServiceOrDefault() { /* do nothing */ }

func TestLookupServiceOrDefault(t *testing.T) {
	defaultService := &testLookupServiceOrDefaultService{}
	service := conf.LookupServiceOrDefault(defaultService)
	require.Equal(t, defaultService, service)
}
