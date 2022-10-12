////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package mixmessages

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
)

func MakeNotificationsCSV(l []*NotificationData) string {
	output := make([][]string, len(l))
	for i, n := range l {
		output[i] = []string{base64.StdEncoding.EncodeToString(n.MessageHash),
			base64.StdEncoding.EncodeToString(n.IdentityFP)}
	}

	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	if err := w.WriteAll(output); err != nil {
		jww.FATAL.Printf("Failed to make notificationsCSV: %+v", err)
	}
	return string(buf.Bytes())
}

func DecodeNotificationsCSV(data string) ([]*NotificationData, error) {
	r := csv.NewReader(strings.NewReader(data))
	read, err := r.ReadAll()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to decode notifications CSV")
	}

	l := make([]*NotificationData, len(read))
	for i, touple := range read {
		messageHash, err := base64.StdEncoding.DecodeString(touple[0])
		if err != nil {
			return nil, errors.WithMessage(err, "Failed decode an element")
		}
		identityFP, err := base64.StdEncoding.DecodeString(touple[1])
		if err != nil {
			return nil, errors.WithMessage(err, "Failed decode an element")
		}
		l[i] = &NotificationData{
			EphemeralID: 0,
			IdentityFP:  identityFP,
			MessageHash: messageHash,
		}
	}
	return l, nil
}
