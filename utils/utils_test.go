////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package utils

import "testing"

func TestGetFullPath_happy(t *testing.T) {
	path := "~/test123"
	newPath := GetFullPath(path)
	if len(path) > len(newPath) {
		t.Errorf("GetFullPath: Expected to replace ~!")
	}
}

func TestGetFullPath_default(t *testing.T) {
	path := "/test123/cat/dog"
	newPath := GetFullPath(path)
	if len(path) != len(newPath) {
		t.Errorf("GetFullPath: Expected no replacement!")
	}
}
