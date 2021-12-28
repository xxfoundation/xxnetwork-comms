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

func UpdateNotificationCSV(l *NotificationData, buf *bytes.Buffer) *bytes.Buffer {
	output := make([]string, 2)
	output = []string{base64.StdEncoding.EncodeToString(l.MessageHash),
		base64.StdEncoding.EncodeToString(l.IdentityFP)}

	w := csv.NewWriter(buf)
	if err := w.Write(output); err != nil {
		jww.FATAL.Printf("Failed to make notificationsCSV: %+v", err)
	}
	return buf
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
