package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTypeValid(t *testing.T) {
	tests := []struct {
		inputName string
		expect    bool
	}{
		{inputName: "type", expect: true},
		{inputName: "type1", expect: false},
		{inputName: "type123&", expect: false},
		{inputName: ".type123", expect: false},
		{inputName: "typ-e123", expect: false},
		{inputName: "", expect: false},
	}

	for _, tt := range tests {
		testname := tt.inputName
		t.Run(testname, func(t *testing.T) {
			ok := isTypeValid(tt.inputName)
			assert.Equalf(t, ok, tt.expect, "want %t, got %t", ok, tt.expect)
		})
	}
}

func TestIsNameValid(t *testing.T) {
	type expect struct {
		code int
		ok   bool
	}
	tests := []struct {
		inputName string
		expect    expect
	}{
		{inputName: "type", expect: expect{0, true}},
		{inputName: "type123", expect: expect{0, true}},
		{inputName: "type123&", expect: expect{http.StatusBadRequest, false}},
		{inputName: "", expect: expect{http.StatusNotFound, false}},
	}

	for _, tt := range tests {
		testname := tt.inputName
		t.Run(testname, func(t *testing.T) {
			code, ok := isNameValid(tt.inputName)
			assert.Equalf(t, code, tt.expect.code, "want %d, got %d", code, tt.expect.code)
			assert.Equalf(t, ok, tt.expect.ok, "want %t, got %t", ok, tt.expect.ok)
		})
	}
}

func TestIsValueValid(t *testing.T) {
	tests := []struct {
		inputName string
		expect    bool
	}{
		{inputName: "0", expect: true},
		{inputName: "123", expect: true},
		{inputName: "123.0", expect: true},
		{inputName: "123.", expect: false},
		{inputName: "type", expect: false},
		{inputName: "", expect: false},
	}

	for _, tt := range tests {
		testname := tt.inputName
		t.Run(testname, func(t *testing.T) {
			ok := isValueValid(tt.inputName)
			assert.Equalf(t, ok, tt.expect, "want %t, got %t", tt.expect, ok)
		})
	}
}
