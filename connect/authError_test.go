package connect

import (
	"errors"
	"testing"
)

func TestAuthError(t *testing.T) {
	expectedAuthErrorStr := "Failed to authenticate id: soisoisoi"
	result:= AuthError("soisoisoi")
	if result==nil{
		t.Error("AuthError did not return an error object")
	}
	if result.Error() != expectedAuthErrorStr{
		t.Errorf("returned error not as expected: Expected: %s, recieved: %s",
			expectedAuthErrorStr, result.Error())
	}
}


func TestIsAuthError(t *testing.T) {
	isAuthError := errors.New("Failed to authenticate id: soisoisoi")

	if !IsAuthError(isAuthError){
		t.Errorf("IsAuthError returned that authError is not an authError")
	}

	notAuthError := errors.New("dont feed the plants")

	if IsAuthError(notAuthError){
		t.Errorf("IsAuthError returned that a non authError is an authError")
	}
}
