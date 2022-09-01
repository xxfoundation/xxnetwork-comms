////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package mixmessages

import (
	"bytes"
	"fmt"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/primitives/netTime"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func checkdigest(t *testing.T, gs signature.GenericRsaSignable) {
	r := reflect.ValueOf(gs).Elem()

	h, err := hash.NewCMixHash()
	if err != nil {
		t.Errorf("Error creating CMix hash: %s", err)
	}

	// Get the old digest
	oldDigest := gs.Digest([]byte{}, h)

	// Setup RNG to fill fields with
	now := netTime.Now()
	rand.Seed(now.Unix())

	// For every value in the passed in struct
	for i := 0; i < r.NumField(); i++ {
		h.Reset()

		// Get the value and type of the field
		valField := r.Field(i)
		typeField := r.Type().Field(i)

		if typeField.Name == "EccSignature" || typeField.Name == "Signature" || typeField.Name == "Errors" || typeField.Name == "state" || typeField.Name == "sizeCache" || typeField.Name == "unknownFields" || strings.Contains(typeField.Name, "XXX") {
			fmt.Printf("Skipping field.\n")
			continue
		}

		// Replace the value with something random
		switch valField.Interface().(type) {
		case []byte:
			randomVal := make([]byte, 4)
			rand.Read(randomVal)
			valField.SetBytes(randomVal)

		case uint32:
			valField.SetUint(uint64(rand.Uint32()))

		case uint64:
			valField.SetUint(rand.Uint64())

		case string:
			valField.SetString(RandStringRunes(4))

		case [][]uint8:
			arr := [][]uint8{
				{uint8(rand.Int()), uint8(rand.Int())},
				{uint8(rand.Int()), uint8(rand.Int())},
			}
			v := reflect.ValueOf(arr)
			valField.Set(v)

		case []uint64:
			arr := []uint64{rand.Uint64(), rand.Uint64(), rand.Uint64()}
			v := reflect.ValueOf(arr)
			valField.Set(v)

		case []*ClientError:
			randClientId := make([]byte, 33)
			rand.Read(randClientId)
			randClientId2 := make([]byte, 33)
			rand.Read(randClientId2)
			rea := []*ClientError{
				{
					ClientId: randClientId,
					Error:    RandStringRunes(4),
				},
				{
					ClientId: randClientId2,
					Error:    RandStringRunes(4),
				},
			}
			v := reflect.ValueOf(rea)
			valField.Set(v)

		default:
			t.Errorf("checkdigest doesn't know how to handle type %s\n", reflect.TypeOf(valField.Interface()))
		}

		fmt.Printf("| Field Name: %s,\n| Field Value: %v,\n| Field Type: %s\n", typeField.Name, valField.Interface(), typeField.Type)

		// Get the new signature
		newDigest := gs.Digest([]byte{}, h)

		// Compare them to make sure the signature doesn't match
		if bytes.Compare(oldDigest, newDigest) == 0 {
			t.Errorf("Digests matched\n")
			fmt.Printf("^^^ FAILED DIGEST CHECK ^^^\n\n\n")
		} else {
			fmt.Printf("Digests did not match, field passed check!\n\n\n")
		}

		oldDigest = newDigest
	}
}
