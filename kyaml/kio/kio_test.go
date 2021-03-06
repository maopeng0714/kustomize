// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package kio_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sigs.k8s.io/kustomize/kyaml/yaml"

	. "sigs.k8s.io/kustomize/kyaml/kio"
)

func TestPipe(t *testing.T) {
	p := Pipeline{
		Inputs:  []Reader{},
		Filters: []Filter{},
		Outputs: []Writer{},
	}

	err := p.Execute()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
}

type mockCallback struct {
	mock.Mock
}

func (c *mockCallback) Callback(op Filter) {
	c.Called(op)
}

func TestPipelineWithCallback(t *testing.T) {
	input := ResourceNodeSlice{yaml.MakeNullNode()}
	noopFilter1 := func(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
		return nodes, nil
	}
	noopFilter2 := func(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
		return nodes, nil
	}
	filters := []Filter{
		FilterFunc(noopFilter1),
		FilterFunc(noopFilter2),
	}
	p := Pipeline{
		Inputs:  []Reader{input},
		Filters: filters,
		Outputs: []Writer{},
	}

	callback := mockCallback{}
	// setup expectations. `Times` means the function is called no more than `times`.
	callback.On("Callback", mock.Anything).Times(len(filters))

	err := p.ExecuteWithCallback(callback.Callback)

	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	callback.AssertNumberOfCalls(t, "Callback", len(filters))

	// assert filters are called in the order they are defined.
	for i, filter := range filters {
		assert.Equal(
			t,
			reflect.ValueOf(callback.Calls[i].Arguments[0]).Pointer(),
			reflect.ValueOf(filter).Pointer(),
		)
	}
}
