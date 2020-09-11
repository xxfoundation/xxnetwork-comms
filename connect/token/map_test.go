package token

import "testing"

// Unit test for NewMap
func TestNewMap(t *testing.T) {
	m := NewMap()
	if m.m == nil {
		t.Errorf("Failed to initialize map")
	}
}

// Unit test for Map.Generate()
func TestMap_Generate(t *testing.T) {
	m := NewMap()
	_ = m.Generate()
}

// Unit test for Map.Validate()
func TestMap_Validate(t *testing.T) {
	m := NewMap()

	token := m.Generate()

	valid := m.Validate(token)
	if !valid {
		t.Error("Failed to validate token")
	}
}
