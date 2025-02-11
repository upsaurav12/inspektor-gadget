// Copyright 2025 The Inspektor Gadget authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFields(t *testing.T) {
	tests := []struct {
		name     string
		field    []string
		expected struct {
			fields     []string
			useDefault bool
		}
		option Option
	}{
		{
			name:  "Testing WithFields with custom values",
			field: []string{"field"},
			expected: struct {
				fields     []string
				useDefault bool
			}{
				fields:     []string{"field"},
				useDefault: false,
			},
			option: func(formatter *Formatter) {
				formatter.fields = []string{"field"}
				formatter.useDefault = false
				formatter.array = false
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			formatter := &Formatter{}
			actual := WithFields(test.field)
			actual(formatter)
			test.option(formatter)
			assert.Equal(t, test.expected.fields, formatter.fields)
			assert.Equal(t, test.expected.useDefault, formatter.useDefault)
		})
	}
}

func TestWithShowAll(t *testing.T) {
	tests := []struct {
		name     string
		val      bool
		expected struct {
			showAll    bool
			useDefault bool
			array      bool
		}
		option Option
	}{
		{
			name: "Testing WithShowAll with custom values",
			val:  true,
			expected: struct {
				showAll    bool
				useDefault bool
				array      bool
			}{
				showAll:    true,
				useDefault: true,
				array:      false,
			},
			option: func(formatter *Formatter) {
				formatter.showAll = true
				formatter.useDefault = true
				formatter.array = false
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			formatter := &Formatter{}
			actual := WithShowAll(test.val)
			actual(formatter)
			test.option(formatter)
			assert.Equal(t, test.expected.showAll, formatter.showAll)
			assert.Equal(t, test.expected.useDefault, formatter.useDefault)
		})
	}
}

func TestWithPretty(t *testing.T) {
	tests := []struct {
		name     string
		val      bool
		indent   string
		expected struct {
			pretty bool
			indent string
		}
		option Option
	}{
		{
			name:   "Testing WithPretty with custom values",
			val:    true,
			indent: "example indent",
			expected: struct {
				pretty bool
				indent string
			}{
				pretty: true,
				indent: "example indent",
			},
			option: func(formatter *Formatter) {
				formatter.pretty = true
				formatter.indent = "example indent"
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			formatter := &Formatter{}
			actual := WithPretty(test.val, test.indent)
			actual(formatter)
			test.option(formatter)
			assert.Equal(t, test.expected.pretty, formatter.pretty)
			assert.Equal(t, test.expected.indent, formatter.indent)
		})
	}
}

func TestWithArray(t *testing.T) {
	tests := []struct {
		name     string
		val      bool
		expected struct {
			array bool
		}
		option Option
	}{
		{
			name: "Testing WithArray with custom values",
			val:  true,
			expected: struct {
				array bool
			}{
				array: true,
			},
			option: func(formatter *Formatter) {
				formatter.array = true
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			formatter := &Formatter{}
			actual := WithArray(test.val)
			actual(formatter)
			test.option(formatter)
			assert.Equal(t, test.expected.array, formatter.array)
		})
	}
}
