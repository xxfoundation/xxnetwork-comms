////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package network

import (
	"gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/comms/testutils"
	"gitlab.com/elixxir/primitives/states"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/ndf"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestBannedNodePartialNDFRemoval(t *testing.T) {
	oldNDF, _ := NewSecuredNdf(testutils.NDF)
	newNDF, _ := ndf.Unmarshal([]byte(`{
	"Timestamp": "2019-06-04T20:48:48-07:00",
	"gateways": [
		{
			"Id": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1],
			"Address": "2.5.3.122",
			"Tls_certificate": "-----BEGIN CERTIFICATE-----\nMIIDgTCCAmmgAwIBAgIJAKLdZ8UigIAeMA0GCSqGSIb3DQEBBQUAMG8xCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQx\nGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEaMBgGA1UEAwwRZ2F0ZXdheSou\nY21peC5yaXAwHhcNMTkwMzA1MTgzNTU0WhcNMjkwMzAyMTgzNTU0WjBvMQswCQYD\nVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJQ2xhcmVtb250\nMRswGQYDVQQKDBJQcml2YXRlZ3JpdHkgQ29ycC4xGjAYBgNVBAMMEWdhdGV3YXkq\nLmNtaXgucmlwMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9+AaxwDP\nxHbhLmn4HoZu0oUM48Qufc6T5XEZTrpMrqJAouXk+61Jc0EFH96/sbj7VyvnXPRo\ngIENbk2Y84BkB9SkRMIXya/gh9dOEDSgnvj/yg24l3bdKFqBMKiFg00PYB30fU+A\nbe3OI/le0I+v++RwH2AV0BMq+T6PcAGjCC1Q1ZB0wP9/VqNMWq5lbK9wD46IQiSi\n+SgIQeE7HoiAZXrGO0Y7l9P3+VRoXjRQbqfn3ETNL9ZvQuarwAYC9Ix5MxUrS5ag\nOmfjc8bfkpYDFAXRXmdKNISJmtCebX2kDrpP8Bdasx7Fzsx59cEUHCl2aJOWXc7R\n5m3juOVL1HUxjQIDAQABoyAwHjAcBgNVHREEFTATghFnYXRld2F5Ki5jbWl4LnJp\ncDANBgkqhkiG9w0BAQUFAAOCAQEAMu3xoc2LW2UExAAIYYWEETggLNrlGonxteSu\njuJjOR+ik5SVLn0lEu22+z+FCA7gSk9FkWu+v9qnfOfm2Am+WKYWv3dJ5RypW/hD\nNXkOYxVJNYFxeShnHohNqq4eDKpdqSxEcuErFXJdLbZP1uNs4WIOKnThgzhkpuy7\ntZRosvOF1X5uL1frVJzHN5jASEDAa7hJNmQ24kh+ds/Ge39fGD8pK31CWhnIXeDo\nvKD7wivi/gSOBtcRWWLvU8SizZkS3hgTw0lSOf5geuzvasCEYlqrKFssj6cTzbCB\nxy3ra3WazRTNTW4TmkHlCUC9I3oWTTxw5iQxF/I2kQQnwR7L3w==\n-----END CERTIFICATE-----"
		},
		{
			"Id": [1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1],
			"Address": "2.2.58.38",
			"Tls_certificate": "-----BEGIN CERTIFICATE-----\nMIIDgTCCAmmgAwIBAgIJAKLdZ8UigIAeMA0GCSqGSIb3DQEBBQUAMG8xCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQx\nGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEaMBgGA1UEAwwRZ2F0ZXdheSou\nY21peC5yaXAwHhcNMTkwMzA1MTgzNTU0WhcNMjkwMzAyMTgzNTU0WjBvMQswCQYD\nVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJQ2xhcmVtb250\nMRswGQYDVQQKDBJQcml2YXRlZ3JpdHkgQ29ycC4xGjAYBgNVBAMMEWdhdGV3YXkq\nLmNtaXgucmlwMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9+AaxwDP\nxHbhLmn4HoZu0oUM48Qufc6T5XEZTrpMrqJAouXk+61Jc0EFH96/sbj7VyvnXPRo\ngIENbk2Y84BkB9SkRMIXya/gh9dOEDSgnvj/yg24l3bdKFqBMKiFg00PYB30fU+A\nbe3OI/le0I+v++RwH2AV0BMq+T6PcAGjCC1Q1ZB0wP9/VqNMWq5lbK9wD46IQiSi\n+SgIQeE7HoiAZXrGO0Y7l9P3+VRoXjRQbqfn3ETNL9ZvQuarwAYC9Ix5MxUrS5ag\nOmfjc8bfkpYDFAXRXmdKNISJmtCebX2kDrpP8Bdasx7Fzsx59cEUHCl2aJOWXc7R\n5m3juOVL1HUxjQIDAQABoyAwHjAcBgNVHREEFTATghFnYXRld2F5Ki5jbWl4LnJp\ncDANBgkqhkiG9w0BAQUFAAOCAQEAMu3xoc2LW2UExAAIYYWEETggLNrlGonxteSu\njuJjOR+ik5SVLn0lEu22+z+FCA7gSk9FkWu+v9qnfOfm2Am+WKYWv3dJ5RypW/hD\nNXkOYxVJNYFxeShnHohNqq4eDKpdqSxEcuErFXJdLbZP1uNs4WIOKnThgzhkpuy7\ntZRosvOF1X5uL1frVJzHN5jASEDAa7hJNmQ24kh+ds/Ge39fGD8pK31CWhnIXeDo\nvKD7wivi/gSOBtcRWWLvU8SizZkS3hgTw0lSOf5geuzvasCEYlqrKFssj6cTzbCB\nxy3ra3WazRTNTW4TmkHlCUC9I3oWTTxw5iQxF/I2kQQnwR7L3w==\n-----END CERTIFICATE-----"
		}
	],
	"nodes": [
		{
			"Id": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2],
			"Address": "18.237.147.105",
			"Tls_certificate": "-----BEGIN CERTIFICATE-----\nMIIDbDCCAlSgAwIBAgIJAOUNtZneIYECMA0GCSqGSIb3DQEBBQUAMGgxCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQx\nGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjETMBEGA1UEAwwKKi5jbWl4LnJp\ncDAeFw0xOTAzMDUxODM1NDNaFw0yOTAzMDIxODM1NDNaMGgxCzAJBgNVBAYTAlVT\nMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQxGzAZBgNV\nBAoMElByaXZhdGVncml0eSBDb3JwLjETMBEGA1UEAwwKKi5jbWl4LnJpcDCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAPP0WyVkfZA/CEd2DgKpcudn0oDh\nDwsjmx8LBDWsUgQzyLrFiVigfUmUefknUH3dTJjmiJtGqLsayCnWdqWLHPJYvFfs\nWYW0IGF93UG/4N5UAWO4okC3CYgKSi4ekpfw2zgZq0gmbzTnXcHF9gfmQ7jJUKSE\ntJPSNzXq+PZeJTC9zJAb4Lj8QzH18rDM8DaL2y1ns0Y2Hu0edBFn/OqavBJKb/uA\nm3AEjqeOhC7EQUjVamWlTBPt40+B/6aFJX5BYm2JFkRsGBIyBVL46MvC02MgzTT9\nbJIJfwqmBaTruwemNgzGu7Jk03hqqS1TUEvSI6/x8bVoba3orcKkf9HsDjECAwEA\nAaMZMBcwFQYDVR0RBA4wDIIKKi5jbWl4LnJpcDANBgkqhkiG9w0BAQUFAAOCAQEA\nneUocN4AbcQAC1+b3To8u5UGdaGxhcGyZBlAoenRVdjXK3lTjsMdMWb4QctgNfIf\nU/zuUn2mxTmF/ekP0gCCgtleZr9+DYKU5hlXk8K10uKxGD6EvoiXZzlfeUuotgp2\nqvI3ysOm/hvCfyEkqhfHtbxjV7j7v7eQFPbvNaXbLa0yr4C4vMK/Z09Ui9JrZ/Z4\ncyIkxfC6/rOqAirSdIp09EGiw7GM8guHyggE4IiZrDslT8V3xIl985cbCxSxeW1R\ntgH4rdEXuVe9+31oJhmXOE9ux2jCop9tEJMgWg7HStrJ5plPbb+HmjoX3nBO04E5\n6m52PyzMNV+2N21IPppKwA==\n-----END CERTIFICATE-----"
		},
		{
			"Id": [1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2],
			"Address": "52.11.136.238",
			"Tls_certificate": "-----BEGIN CERTIFICATE-----\nMIIDbDCCAlSgAwIBAgIJAOUNtZneIYECMA0GCSqGSIb3DQEBBQUAMGgxCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQx\nGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjETMBEGA1UEAwwKKi5jbWl4LnJp\ncDAeFw0xOTAzMDUxODM1NDNaFw0yOTAzMDIxODM1NDNaMGgxCzAJBgNVBAYTAlVT\nMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQxGzAZBgNV\nBAoMElByaXZhdGVncml0eSBDb3JwLjETMBEGA1UEAwwKKi5jbWl4LnJpcDCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAPP0WyVkfZA/CEd2DgKpcudn0oDh\nDwsjmx8LBDWsUgQzyLrFiVigfUmUefknUH3dTJjmiJtGqLsayCnWdqWLHPJYvFfs\nWYW0IGF93UG/4N5UAWO4okC3CYgKSi4ekpfw2zgZq0gmbzTnXcHF9gfmQ7jJUKSE\ntJPSNzXq+PZeJTC9zJAb4Lj8QzH18rDM8DaL2y1ns0Y2Hu0edBFn/OqavBJKb/uA\nm3AEjqeOhC7EQUjVamWlTBPt40+B/6aFJX5BYm2JFkRsGBIyBVL46MvC02MgzTT9\nbJIJfwqmBaTruwemNgzGu7Jk03hqqS1TUEvSI6/x8bVoba3orcKkf9HsDjECAwEA\nAaMZMBcwFQYDVR0RBA4wDIIKKi5jbWl4LnJpcDANBgkqhkiG9w0BAQUFAAOCAQEA\nneUocN4AbcQAC1+b3To8u5UGdaGxhcGyZBlAoenRVdjXK3lTjsMdMWb4QctgNfIf\nU/zuUn2mxTmF/ekP0gCCgtleZr9+DYKU5hlXk8K10uKxGD6EvoiXZzlfeUuotgp2\nqvI3ysOm/hvCfyEkqhfHtbxjV7j7v7eQFPbvNaXbLa0yr4C4vMK/Z09Ui9JrZ/Z4\ncyIkxfC6/rOqAirSdIp09EGiw7GM8guHyggE4IiZrDslT8V3xIl985cbCxSxeW1R\ntgH4rdEXuVe9+31oJhmXOE9ux2jCop9tEJMgWg7HStrJ5plPbb+HmjoX3nBO04E5\n6m52PyzMNV+2N21IPppKwA==\n-----END CERTIFICATE-----"
		}
	],
	"registration": {
		"Address": "92.42.125.61",
		"Tls_certificate": "-----BEGIN CERTIFICATE-----\nMIIDkDCCAnigAwIBAgIJAJnjosuSsP7gMA0GCSqGSIb3DQEBBQUAMHQxCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQx\nGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEfMB0GA1UEAwwWcmVnaXN0cmF0\naW9uKi5jbWl4LnJpcDAeFw0xOTAzMDUyMTQ5NTZaFw0yOTAzMDIyMTQ5NTZaMHQx\nCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFy\nZW1vbnQxGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEfMB0GA1UEAwwWcmVn\naXN0cmF0aW9uKi5jbWl4LnJpcDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC\nggEBAOQKvqjdh35o+MECBhCwopJzPlQNmq2iPbewRNtI02bUNK3kLQUbFlYdzNGZ\nS4GYXGc5O+jdi8Slx82r1kdjz5PPCNFBARIsOP/L8r3DGeW+yeJdgBZjm1s3ylka\nmt4Ajiq/bNjysS6L/WSOp+sVumDxtBEzO/UTU1O6QRnzUphLaiWENmErGvsH0CZV\nq38Ia58k/QjCAzpUcYi4j2l1fb07xqFcQD8H6SmUM297UyQosDrp8ukdIo31Koxr\n4XDnnNNsYStC26tzHMeKuJ2Wl+3YzsSyflfM2YEcKE31sqB9DS36UkJ8J84eLsHN\nImGg3WodFAviDB67+jXDbB30NkMCAwEAAaMlMCMwIQYDVR0RBBowGIIWcmVnaXN0\ncmF0aW9uKi5jbWl4LnJpcDANBgkqhkiG9w0BAQUFAAOCAQEAF9mNzk+g+o626Rll\nt3f3/1qIyYQrYJ0BjSWCKYEFMCgZ4JibAJjAvIajhVYERtltffM+YKcdE2kTpdzJ\n0YJuUnRfuv6sVnXlVVugUUnd4IOigmjbCdM32k170CYMm0aiwGxl4FrNa8ei7AIa\nx/s1n+sqWq3HeW5LXjnoVb+s3HeCWIuLfcgrurfye8FnNhy14HFzxVYYefIKm0XL\n+DPlcGGGm/PPYt3u4a2+rP3xaihc65dTa0u5tf/XPXtPxTDPFj2JeQDFxo7QRREb\nPD89CtYnwuP937CrkvCKrL0GkW1FViXKqZY9F5uhxrvLIpzhbNrs/EbtweY35XGL\nDCCMkg==\n-----END CERTIFICATE-----"
	},
	"notification": {
		"Address": "notification.default.cmix.rip",
		"Tls_certificate": "-----BEGIN CERTIFICATE-----\nMIIDkDCCAnigAwIBAgIJAJnjosuSsP7gMA0GCSqGSIb3DQEBBQUAMHQxCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQx\nGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEfMB0GA1UEAwwWcmVnaXN0cmF0\naW9uKi5jbWl4LnJpcDAeFw0xOTAzMDUyMTQ5NTZaFw0yOTAzMDIyMTQ5NTZaMHQx\nCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFy\nZW1vbnQxGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEfMB0GA1UEAwwWcmVn\naXN0cmF0aW9uKi5jbWl4LnJpcDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC\nggEBAOQKvqjdh35o+MECBhCwopJzPlQNmq2iPbewRNtI02bUNK3kLQUbFlYdzNGZ\nS4GYXGc5O+jdi8Slx82r1kdjz5PPCNFBARIsOP/L8r3DGeW+yeJdgBZjm1s3ylka\nmt4Ajiq/bNjysS6L/WSOp+sVumDxtBEzO/UTU1O6QRnzUphLaiWENmErGvsH0CZV\nq38Ia58k/QjCAzpUcYi4j2l1fb07xqFcQD8H6SmUM297UyQosDrp8ukdIo31Koxr\n4XDnnNNsYStC26tzHMeKuJ2Wl+3YzsSyflfM2YEcKE31sqB9DS36UkJ8J84eLsHN\nImGg3WodFAviDB67+jXDbB30NkMCAwEAAaMlMCMwIQYDVR0RBBowGIIWcmVnaXN0\ncmF0aW9uKi5jbWl4LnJpcDANBgkqhkiG9w0BAQUFAAOCAQEAF9mNzk+g+o626Rll\nt3f3/1qIyYQrYJ0BjSWCKYEFMCgZ4JibAJjAvIajhVYERtltffM+YKcdE2kTpdzJ\n0YJuUnRfuv6sVnXlVVugUUnd4IOigmjbCdM32k170CYMm0aiwGxl4FrNa8ei7AIa\nx/s1n+sqWq3HeW5LXjnoVb+s3HeCWIuLfcgrurfye8FnNhy14HFzxVYYefIKm0XL\n+DPlcGGGm/PPYt3u4a2+rP3xaihc65dTa0u5tf/XPXtPxTDPFj2JeQDFxo7QRREb\nPD89CtYnwuP937CrkvCKrL0GkW1FViXKqZY9F5uhxrvLIpzhbNrs/EbtweY35XGL\nDCCMkg==\n-----END CERTIFICATE-----"
	},
	"udb": {
		"Id": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3]
	},
	"E2e": {
		"Prime": "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AACAA68FFFFFFFFFFFFFFFF",
		"Small_prime": "7FFFFFFFFFFFFFFFE487ED5110B4611A62633145C06E0E68948127044533E63A0105DF531D89CD9128A5043CC71A026EF7CA8CD9E69D218D98158536F92F8A1BA7F09AB6B6A8E122F242DABB312F3F637A262174D31BF6B585FFAE5B7A035BF6F71C35FDAD44CFD2D74F9208BE258FF324943328F6722D9EE1003E5C50B1DF82CC6D241B0E2AE9CD348B1FD47E9267AFC1B2AE91EE51D6CB0E3179AB1042A95DCF6A9483B84B4B36B3861AA7255E4C0278BA3604650C10BE19482F23171B671DF1CF3B960C074301CD93C1D17603D147DAE2AEF837A62964EF15E5FB4AAC0B8C1CCAA4BE754AB5728AE9130C4C7D02880AB9472D455655347FFFFFFFFFFFFFFF",
		"Generator": "02"
	},
	"Cmix": {
		"Prime": "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AACAA68FFFFFFFFFFFFFFFF",
		"Small_prime": "7FFFFFFFFFFFFFFFE487ED5110B4611A62633145C06E0E68948127044533E63A0105DF531D89CD9128A5043CC71A026EF7CA8CD9E69D218D98158536F92F8A1BA7F09AB6B6A8E122F242DABB312F3F637A262174D31BF6B585FFAE5B7A035BF6F71C35FDAD44CFD2D74F9208BE258FF324943328F6722D9EE1003E5C50B1DF82CC6D241B0E2AE9CD348B1FD47E9267AFC1B2AE91EE51D6CB0E3179AB1042A95DCF6A9483B84B4B36B3861AA7255E4C0278BA3604650C10BE19482F23171B671DF1CF3B960C074301CD93C1D17603D147DAE2AEF837A62964EF15E5FB4AAC0B8C1CCAA4BE754AB5728AE9130C4C7D02880AB9472D455655347FFFFFFFFFFFFFFF",
		"Generator": "02"
	}
}`))

	rmNodes, err := getBannedNodes(oldNDF.Get().Nodes, newNDF.Nodes)
	if err != nil {
		t.Errorf("Failed to run getBannedNodes")
	}
	if len(rmNodes) != 1 {
		t.Errorf("Incorrect number of nodes removed: removed %v, expected 1", len(rmNodes))
	}

	onid, err := id.Unmarshal(oldNDF.Get().Nodes[2].ID)
	if err != nil {
		t.Errorf("Failed to unmarshal oldNDF node ID 2")
	}
	if rmNodes[0].Cmp(onid) != true {
		t.Errorf("Function removed the wrong node")
	}
}

// Happy path
func TestNewInstanceTesting(t *testing.T) {
	_, err := NewInstanceTesting(&connect.ProtoComms{}, testutils.NDF, testutils.NDF, nil, nil, t)
	if err != nil {
		t.Errorf("Unable to create test instance: %+v", err)
	}
}

// Error path: pass in a non testing argument into the constructor
func TestNewInstanceTesting_Error(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	_, err := NewInstanceTesting(&connect.ProtoComms{}, testutils.NDF, testutils.NDF, nil, nil, nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error case, should not be able to create instance when testing argument is nil")

}

//tests newInstance errors properly when there is no NDF
func TestNewInstance_NilNDFs(t *testing.T) {
	_, err := NewInstance(&connect.ProtoComms{}, nil, nil, nil, 0, false)
	if err == nil {
		t.Errorf("Creation of NewInstance without an ndf succeded")
	} else if !strings.Contains(err.Error(), "Cannot create a network "+
		"instance without an NDF") {
		t.Errorf("Creation of NewInstance without an ndf returned "+
			"the wrong error: %s", err.Error())
	}
}

func TestInstance_GetFullNdf(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)
	i := Instance{
		full: secured,
	}
	if i.GetFullNdf() == nil {
		t.Error("Failed to retrieve full ndf")
	}
}

func TestInstance_GetPartialNdf(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)
	i := Instance{
		partial: secured,
	}
	if i.GetPartialNdf() == nil {
		t.Error("Failed to retrieve partial ndf")
	}
}

func TestInstance_GetRound(t *testing.T) {
	i := Instance{
		roundData: ds.NewData(),
	}

	// Construct a mock round object
	ri := &mixmessages.RoundInfo{ID: uint64(1)}
	testutils.SignRoundInfoRsa(ri, t)

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	rnd := ds.NewRound(ri, pubKey, nil)

	_ = i.roundData.UpsertRound(rnd)
	r, err := i.GetRound(id.Round(1))
	if err != nil || r == nil {
		t.Errorf("Failed to retrieve round: %+v", err)
	}
}

func TestInstance_GetWrappedRound(t *testing.T) {
	i := Instance{
		roundData: ds.NewData(),
	}

	// Construct a mock round object
	ri := &mixmessages.RoundInfo{ID: uint64(1)}
	testutils.SignRoundInfoRsa(ri, t)

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}

	rnd := ds.NewRound(ri, pubKey, nil)

	_ = i.roundData.UpsertRound(rnd)
	retrieved, err := i.GetWrappedRound(id.Round(1))
	if err != nil || retrieved == nil {
		t.Errorf("Failed to retrieve round: %+v", err)
	}

	if !reflect.DeepEqual(rnd, retrieved) {
		t.Errorf("Retrieved value did not match expected!"+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", rnd, retrieved)
	}
}

func TestInstance_GetRoundUpdate(t *testing.T) {
	i := Instance{
		roundUpdates: ds.NewUpdates(),
	}

	ri := &mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)}
	if err := testutils.SignRoundInfoRsa(ri, t); err != nil {
		t.Fatalf("Cannot sign round info: %v", err)
	}
	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load test key: %v", err)
	}
	rnd := ds.NewRound(ri, pubKey, nil)

	_ = i.roundUpdates.AddRound(rnd)
	r, err := i.GetRoundUpdate(1)
	if err != nil || r == nil {
		t.Errorf("Failed to retrieve round update: %+v", err)
	}
}

func TestInstance_GetRoundUpdates(t *testing.T) {
	i := Instance{
		roundUpdates: ds.NewUpdates(),
	}
	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}

	roundInfoOne := &mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)}
	testutils.SignRoundInfoRsa(roundInfoOne, t)
	roundInfoTwo := &mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(2)}
	testutils.SignRoundInfoRsa(roundInfoTwo, t)
	roundOne := ds.NewRound(roundInfoOne, pubKey, nil)
	roundTwo := ds.NewRound(roundInfoTwo, pubKey, nil)

	_ = i.roundUpdates.AddRound(roundOne)
	_ = i.roundUpdates.AddRound(roundTwo)
	r := i.GetRoundUpdates(1)
	if r == nil {
		t.Errorf("Failed to retrieve round updates")
	}
}

func setupComm(t *testing.T) (*Instance, *mixmessages.NDF) {
	privKey, err := testutils.LoadPrivateKeyTesting(t)
	if err != nil {
		t.Errorf("Could not load key: %v", err)
		t.FailNow()
	}
	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
	if err != nil {
		t.Errorf("Could not generate rsa key: %s", err)
	}

	f := &mixmessages.NDF{}
	f.Ndf = []byte(testutils.ExampleJSON)
	baseNDF := testutils.NDF

	if err != nil {
		t.Errorf("Could not generate serialized ndf: %s", err)
	}

	err = signature.SignRsa(f, privKey)
	if err != nil {
		t.Fatalf("Failed to sign round info: %v", err)
	}
	testManager := connect.NewManagerTesting(t)
	pc := &connect.ProtoComms{
		Id:      id.NewIdFromString("User", id.User, t),
		Manager: testManager,
	}
	i, err := NewInstance(pc, baseNDF, baseNDF, nil, 0, false)
	if err != nil {
		t.Error(nil)
	}

	_, err = i.comm.AddHost(&id.Permissioning, "0.0.0.0:4200", pub, connect.GetDefaultHostParams())
	if err != nil {
		t.Errorf("Failed to add permissioning host: %+v", err)
	}
	return i, f
}

func TestInstance_RoundUpdate(t *testing.T) {
	msg := &mixmessages.RoundInfo{
		ID:        2,
		UpdateID:  4,
		State:     6,
		BatchSize: 8,
	}
	privKey, err := testutils.LoadPrivateKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load private key: %v", err)
		t.FailNow()
	}
	err = signature.SignRsa(msg, privKey)
	if err != nil {
		t.Fatalf("Failed to sign round info: %v", err)
	}
	testManager := connect.NewManagerTesting(t)
	pc := connect.ProtoComms{
		Manager: testManager,
	}
	i, err := NewInstance(&pc, testutils.NDF, testutils.NDF, nil, 0, false)
	pub := testkeys.LoadFromPath(testkeys.GetGatewayCertPath())

	_, err = i.RoundUpdate(msg)
	if err == nil {
		t.Error("Should have failed to get perm host")
	}

	_, err = i.comm.AddHost(&id.Permissioning, "0.0.0.0:4200", pub, connect.GetDefaultHostParams())
	if err != nil {
		t.Errorf("failed to add bad host: %+v", err)
	}
	_, err = i.RoundUpdate(msg)
	// Fixme
	/*	if err == nil {
		t.Error("Should have failed to verify")
	}*/

	i, _ = setupComm(t)

	_, err = i.RoundUpdate(msg)
	if err != nil {
		t.Errorf("Failed to update ndf: %+v", err)
	}
}

func TestInstance_UpdateFullNdf(t *testing.T) {
	i, f := setupComm(t)

	err := i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf("Failed to update ndf: %+v", err)
	}
}

func TestInstance_UpdateFullNdf_nil(t *testing.T) {
	i, f := setupComm(t)
	i.full = nil

	err := i.UpdateFullNdf(f)
	if err == nil {
		t.Errorf("Full NDF update succeded when it shouldnt")
	} else if !strings.Contains(err.Error(),
		"Cannot update the full ndf when it is nil") {
		t.Errorf("Full NDF update when nil failed incorrectly: %s",
			err.Error())
	}
}

func TestInstance_UpdatePartialNdf(t *testing.T) {
	i, f := setupComm(t)
	err := i.UpdatePartialNdf(f)
	if err != nil {
		t.Errorf("Failed to update ndf: %+v", err)
	}
}

func TestInstance_UpdatePartialNdf_nil(t *testing.T) {
	i, f := setupComm(t)
	i.partial = nil

	err := i.UpdatePartialNdf(f)
	if err == nil {
		t.Errorf("Partial NDF update succeded when it shouldnt")
	} else if !strings.Contains(err.Error(),
		"Cannot update the partial ndf when it is nil") {
		t.Errorf("Partial NDF update when nil failed incorrectly: %s",
			err.Error())
	}
}

func TestInstance_GetLastRoundID(t *testing.T) {
	i := Instance{
		roundData: ds.NewData(),
	}

	expectedLastRound := 23
	ri := &mixmessages.RoundInfo{ID: uint64(expectedLastRound)}
	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	rnd := ds.NewRound(ri, pubKey, nil)

	_ = i.roundData.UpsertRound(rnd)
	lastRound := i.GetLastRoundID()
	if id.Round(expectedLastRound-1) != lastRound {
		t.Errorf("GetLastRoundID did not return expected value."+
			"\n\tExpected: %d"+
			"\n\tRecieved: %d", expectedLastRound-1, lastRound)
	}
}

func TestInstance_GetLastUpdateID(t *testing.T) {
	i := Instance{
		roundUpdates: ds.NewUpdates(),
	}

	expectedUpdateId := 5
	ri := &mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(expectedUpdateId)}
	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	rnd := ds.NewRound(ri, pubKey, nil)

	_ = i.roundUpdates.AddRound(rnd)
	lastUpdateId := i.GetLastUpdateID()

	if lastUpdateId != expectedUpdateId {
		t.Errorf("Last update Id returned unexpected result."+
			"\n\tExpected: %d"+
			"\n\tReceived: %d", expectedUpdateId, lastUpdateId)
	}
}

func TestInstance_GetOldestRoundID(t *testing.T) {
	i := Instance{
		roundData: ds.NewData(),
	}

	expectedOldRoundId := id.Round(0)
	expectedOldRoundInfo := &mixmessages.RoundInfo{ID: uint64(expectedOldRoundId)}
	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	expectedOldRound := ds.NewRound(expectedOldRoundInfo, pubKey, nil)

	mockRoundInfo := &mixmessages.RoundInfo{ID: uint64(2)}
	mockRound := ds.NewRound(mockRoundInfo, pubKey, nil)

	_ = i.roundData.UpsertRound(expectedOldRound)
	_ = i.roundData.UpsertRound(mockRound)

	returned := i.GetOldestRoundID()
	if returned != expectedOldRoundId {
		t.Errorf("Failed to get oldest round from buffer."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedOldRoundId, returned)
	}
}

// Test which forces a full buffer, causing overwriting of old rounds
func TestInstance_GetOldestRoundID_ManyRounds(t *testing.T) {
	testInstance := Instance{
		roundData: ds.NewData(),
	}

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}

	// Ensure a circle back in the round buffer
	for i := 1; i <= ds.RoundInfoBufLen; i++ {
		ri := &mixmessages.RoundInfo{ID: uint64(i)}
		rnd := ds.NewRound(ri, pubKey, nil)
		_ = testInstance.roundData.UpsertRound(rnd)

	}

	// This will have oldest round as 0, until we reach RoundInfoBufLen, then
	// round 0 will be overwritten by the newest round,
	// moving the oldest round to round 1
	expected := id.Round(1)
	returned := testInstance.GetOldestRoundID()
	if returned != expected {
		t.Errorf("Failed to get oldest round from buffer."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", 1, returned)
	}
}

func TestInstance_UpdateGatewayConnections(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)
	testManager := connect.NewManagerTesting(t)
	pc := &connect.ProtoComms{
		Id:      id.NewIdFromString("User", id.User, t),
		Manager: testManager,
	}
	i := Instance{
		full:       secured,
		comm:       pc,
		ipOverride: ds.NewIpOverrideList(),
	}
	err := i.UpdateGatewayConnections()
	if err != nil {
		t.Errorf("Failed to update gateway connections from full: %+v", err)
	}

	i = Instance{
		partial:    secured,
		comm:       pc,
		ipOverride: ds.NewIpOverrideList(),
	}
	err = i.UpdateGatewayConnections()
	if err != nil {
		t.Errorf("Failed to update gateway connections from partial: %+v", err)
	}

	i = Instance{}
	err = i.UpdateGatewayConnections()
	if err == nil {
		t.Error("Should error when attempting update with no ndf")
	}
}

// Tests that UpdateGatewayConnections() returns an error when a Gateway ID
// collides with a hard coded ID.
func TestInstance_UpdateGatewayConnections_GatewayIdError(t *testing.T) {
	testDef := *testutils.NDF
	testDef.Nodes = []ndf.Node{{ID: id.TempGateway.Marshal()}}
	secured, _ := NewSecuredNdf(&testDef)
	testManager := connect.NewManagerTesting(t)
	pc := &connect.ProtoComms{
		Manager: testManager,
	}

	i := Instance{
		full:       secured,
		comm:       pc,
		ipOverride: ds.NewIpOverrideList(),
	}
	err := i.UpdateGatewayConnections()
	if err == nil {
		t.Errorf("UpdateGatewayConnections() failed to produce an error when " +
			"the Gateway ID collides with a hard coded ID.")
	}
}

func TestInstance_UpdateNodeConnections(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)
	testManager := connect.NewManagerTesting(t)
	pc := &connect.ProtoComms{
		Manager: testManager,
	}

	i := Instance{
		full:       secured,
		comm:       pc,
		ipOverride: ds.NewIpOverrideList(),
	}
	err := i.UpdateNodeConnections()
	if err != nil {
		t.Errorf("Failed to update node connections from full: %+v", err)
	}

	i = Instance{
		partial:    secured,
		comm:       pc,
		ipOverride: ds.NewIpOverrideList(),
	}
	err = i.UpdateNodeConnections()
	if err != nil {
		t.Errorf("Failed to update node connections from partial: %+v", err)
	}

	i = Instance{}
	err = i.UpdateNodeConnections()
	if err == nil {
		t.Error("Should error when attempting update with no ndf")
	}
}

// Tests that UpdateNodeConnections() returns an error when a Node ID collides
// with a hard coded ID.
func TestInstance_UpdateNodeConnections_NodeIdError(t *testing.T) {
	testDef := *testutils.NDF
	testDef.Nodes = []ndf.Node{{ID: id.Permissioning.Marshal()}}
	secured, _ := NewSecuredNdf(&testDef)
	testManager := connect.NewManagerTesting(t)
	pc := &connect.ProtoComms{
		Manager: testManager,
	}

	i := Instance{
		full:       secured,
		comm:       pc,
		ipOverride: ds.NewIpOverrideList(),
	}
	err := i.UpdateNodeConnections()
	if err == nil {
		t.Errorf("UpdateNodeConnections() failed to produce an error when the " +
			"Node ID collides with a hard coded ID.")
	}
}

// Happy path: Tests GetPermissioningAddress with the full ndf set, the partial ndf set
// and no ndf set
func TestInstance_GetPermissioningAddress(t *testing.T) {
	// Create populated ndf (secured) and empty ndf
	secured, _ := NewSecuredNdf(testutils.NDF)

	// Create an instance object, setting full to be populated
	// and partial to be empty
	fullNdfInstance := Instance{
		full: secured,
	}

	// Expected address gotten from testutils.NDF
	expectedAddress := "92.42.125.61"

	// GetPermissioningAddress from the instance and compare with the expected value
	receivedAddress := fullNdfInstance.GetPermissioningAddress()
	if expectedAddress != receivedAddress {
		t.Errorf("GetPermissioningAddress did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedAddress, receivedAddress)
	}

	// Create an instance object, setting partial to be populated
	// and full to be empty
	partialNdfInstance := Instance{
		partial: secured,
	}

	// GetPermissioningAddress from the instance and compare with the expected value
	receivedAddress = partialNdfInstance.GetPermissioningAddress()
	if expectedAddress != receivedAddress {
		t.Errorf("GetPermissioningAddress did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedAddress, receivedAddress)
	}

	// Create an instance object, setting no ndf
	noNdfInstance := Instance{}

	// GetPermissioningAddress, should be an empty string as no ndf's are set
	receivedAddress = noNdfInstance.GetPermissioningAddress()
	if receivedAddress != "" {
		t.Errorf("GetPermissioningAddress did not get expected value!"+
			"No ndf set, address should be an empty string. "+
			"\n\tReceived: %+v", receivedAddress)
	}

}

// Happy path
func TestInstance_GetCmixGroup(t *testing.T) {
	expectedGroup := ds.NewGroup()

	i := Instance{
		cmixGroup: expectedGroup,
	}

	receivedGroup := i.GetCmixGroup()

	if !reflect.DeepEqual(expectedGroup.Get(), receivedGroup) {
		t.Errorf("Getter didn't get expected value! "+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedGroup, receivedGroup)
	}

}

// Happy path
func TestInstance_GetE2EGroup(t *testing.T) {
	expectedGroup := ds.NewGroup()

	i := Instance{
		e2eGroup: expectedGroup,
	}

	receivedGroup := i.GetE2EGroup()

	if !reflect.DeepEqual(expectedGroup.Get(), receivedGroup) {
		t.Errorf("Getter didn't get expected value! "+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedGroup, receivedGroup)
	}
}

// Happy path: Tests GetPermissioningCert with the full ndf set, the partial ndf set
// and no ndf set
func TestInstance_GetPermissioningCert(t *testing.T) {

	// Create populated ndf (secured) and empty ndf
	secured, _ := NewSecuredNdf(testutils.NDF)
	// Create an instance object, setting full to be populated
	// and partial to be empty
	fullNdfInstance := Instance{
		full: secured,
	}

	// Expected cert gotten from testutils.NDF
	expectedCert := "-----BEGIN CERTIFICATE-----\nMIIDkDCCAnigAwIBAgIJAJnjosuSsP7gMA0GCSqGSIb3DQEBBQUAMHQxCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQx\nGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEfMB0GA1UEAwwWcmVnaXN0cmF0\naW9uKi5jbWl4LnJpcDAeFw0xOTAzMDUyMTQ5NTZaFw0yOTAzMDIyMTQ5NTZaMHQx\nCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFy\nZW1vbnQxGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEfMB0GA1UEAwwWcmVn\naXN0cmF0aW9uKi5jbWl4LnJpcDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC\nggEBAOQKvqjdh35o+MECBhCwopJzPlQNmq2iPbewRNtI02bUNK3kLQUbFlYdzNGZ\nS4GYXGc5O+jdi8Slx82r1kdjz5PPCNFBARIsOP/L8r3DGeW+yeJdgBZjm1s3ylka\nmt4Ajiq/bNjysS6L/WSOp+sVumDxtBEzO/UTU1O6QRnzUphLaiWENmErGvsH0CZV\nq38Ia58k/QjCAzpUcYi4j2l1fb07xqFcQD8H6SmUM297UyQosDrp8ukdIo31Koxr\n4XDnnNNsYStC26tzHMeKuJ2Wl+3YzsSyflfM2YEcKE31sqB9DS36UkJ8J84eLsHN\nImGg3WodFAviDB67+jXDbB30NkMCAwEAAaMlMCMwIQYDVR0RBBowGIIWcmVnaXN0\ncmF0aW9uKi5jbWl4LnJpcDANBgkqhkiG9w0BAQUFAAOCAQEAF9mNzk+g+o626Rll\nt3f3/1qIyYQrYJ0BjSWCKYEFMCgZ4JibAJjAvIajhVYERtltffM+YKcdE2kTpdzJ\n0YJuUnRfuv6sVnXlVVugUUnd4IOigmjbCdM32k170CYMm0aiwGxl4FrNa8ei7AIa\nx/s1n+sqWq3HeW5LXjnoVb+s3HeCWIuLfcgrurfye8FnNhy14HFzxVYYefIKm0XL\n+DPlcGGGm/PPYt3u4a2+rP3xaihc65dTa0u5tf/XPXtPxTDPFj2JeQDFxo7QRREb\nPD89CtYnwuP937CrkvCKrL0GkW1FViXKqZY9F5uhxrvLIpzhbNrs/EbtweY35XGL\nDCCMkg==\n-----END CERTIFICATE-----"

	// GetPermissioningCert from the instance and compare with the expected value
	receivedCert := fullNdfInstance.GetPermissioningCert()
	if expectedCert != receivedCert {
		t.Errorf("GetPermissioningCert did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedCert, receivedCert)
	}

	// Create an instance object, setting partial to be populated
	// and full to be empty
	partialNdfInstance := Instance{
		partial: secured,
	}

	// GetPermissioningCert from the instance and compare with the expected value
	receivedCert = partialNdfInstance.GetPermissioningCert()
	if expectedCert != receivedCert {
		t.Errorf("GetPermissioningCert did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedCert, receivedCert)
	}

	// Create an instance object, setting no ndf
	noNdfInstance := Instance{}

	// GetPermissioningCert, should be an empty string as no ndf's are set
	receivedCert = noNdfInstance.GetPermissioningCert()
	if receivedCert != "" {
		t.Errorf("GetPermissioningCert did not get expected value!"+
			"No ndf set, cert should be an empty string. "+
			"\n\tReceived: %+v", receivedCert)
	}

}

// Error path: nil ndf is in the instance should cause a seg fault
func TestInstance_GetPermissioningAddress_NilCase(t *testing.T) {
	// Handle expected seg fault here
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected error case, should seg fault when a nil ndf is passed through")
		}
	}()

	// Create a nil ndf
	nilNdf, _ := NewSecuredNdf(nil)

	// Create an instance object with this nil ndf
	nilNdfInstance := Instance{
		full:    nilNdf,
		partial: nilNdf,
	}

	// Attempt to call getter, should seg fault
	nilNdfInstance.GetPermissioningAddress()
}

// Error path: nil ndf is in the instance should cause a seg fault
func TestInstance_GetPermissioningCert_NilCase(t *testing.T) {
	// Handle expected seg fault here
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected error case, should seg fault when a nil ndf is passed through")
		}
	}()

	// Create a nil ndf
	nilNdf, _ := NewSecuredNdf(nil)

	// Create an instance object with this nil ndf
	nilNdfInstance := Instance{
		full:    nilNdf,
		partial: nilNdf,
	}

	// Attempt to call getter, should seg fault
	nilNdfInstance.GetPermissioningCert()
}

// Happy path: Tests GetEllipticPublicKey with the full ndf set, the partial ndf set
// and no ndf set
func TestInstance_GetEllipticPublicKey(t *testing.T) {

	// Create populated ndf (secured) and empty ndf
	secured, _ := NewSecuredNdf(testutils.NDF)
	// Create an instance object, setting full to be populated
	// and partial to be empty
	fullNdfInstance := Instance{
		full: secured,
	}

	// Expected cert gotten from testutils.NDF
	expectedKey := "MqaJJ3GjFisNRM6LRedRnooi14gepMaQxyWctXVU/w4="

	// GetEllipticPublicKey from the instance and compare with the expected value
	receivedKey := fullNdfInstance.GetEllipticPublicKey()
	if expectedKey != receivedKey {
		t.Errorf("GetEllipticPublicKey did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedKey, receivedKey)
	}

	// Create an instance object, setting partial to be populated
	// and full to be empty
	partialNdfInstance := Instance{
		partial: secured,
	}

	// GetEllipticPublicKey from the instance and compare with the expected value
	receivedKey = partialNdfInstance.GetEllipticPublicKey()
	if expectedKey != receivedKey {
		t.Errorf("GetEllipticPublicKey did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedKey, receivedKey)
	}

	// Create an instance object, setting no ndf
	noNdfInstance := Instance{}

	// GetEllipticPublicKey, should be an empty string as no ndf's are set
	receivedKey = noNdfInstance.GetEllipticPublicKey()
	if receivedKey != "" {
		t.Errorf("GetEllipticPublicKey did not get expected value!"+
			"No ndf set, cert should be an empty string. "+
			"\n\tReceived: %+v", receivedKey)
	}

}

// Error path: nil ndf is in the instance should cause a seg fault
func TestInstance_GetEllipticPublicKey_NilCase(t *testing.T) {
	// Handle expected seg fault here
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected error case, should seg fault when a nil ndf is passed through")
		}
	}()

	// Create a nil ndf
	nilNdf, _ := NewSecuredNdf(nil)

	// Create an instance object with this nil ndf
	nilNdfInstance := Instance{
		full:    nilNdf,
		partial: nilNdf,
	}

	// Attempt to call getter, should seg fault
	nilNdfInstance.GetEllipticPublicKey()
}

// GetPermissioningId should fetch the value of id.PERMISSIONING in primitives
func TestInstance_GetPermissioningId(t *testing.T) {
	// Create an instance object,
	instance := Instance{}

	receivedId := instance.GetPermissioningId()

	if receivedId != &id.Permissioning {
		t.Errorf("GetPermissioningId did not get value from primitives"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", id.Permissioning, receivedId)
	}
}

// Full smoke test for Node Event Model
func TestInstance_NodeEventModel(t *testing.T) {
	i, f := setupComm(t)

	// Set up the channels
	addNode := make(chan NodeGateway, 10)
	removeNode := make(chan *id.ID, 10)
	addGateway := make(chan NodeGateway, 10)
	removeGateway := make(chan *id.ID, 10)
	i.SetRemoveGatewayChan(removeGateway)
	i.SetRemoveNodeChan(removeNode)
	i.SetAddGatewayChan(addGateway)
	i.SetAddNodeChan(addNode)

	// Install the NDF
	err := i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf("Unable to initalize group: %+v", err)
	}

	// Set up channels that count number of items they receive
	counter := uint32(0)
	go func() {
		for range addNode {
			atomic.AddUint32(&counter, 1)
		}
	}()
	go func() {
		for range addGateway {
			atomic.AddUint32(&counter, 1)
		}
	}()

	// Trigger sending to channels
	err = i.UpdateNodeConnections()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = i.UpdateGatewayConnections()
	if err != nil {
		t.Errorf(err.Error())
	}

	// Get the NDF
	newNdf, err := ndf.Unmarshal(f.Ndf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// Verify channels received the correct amount of information
	timeout := 5
	for {
		counterVal := atomic.LoadUint32(&counter)
		totalNodeGateways := uint32(len(newNdf.Nodes) + len(newNdf.Gateways))
		if counterVal == totalNodeGateways {
			break
		} else {
			timeout -= 1
			if timeout == 0 {
				t.Errorf("Unable to properly add nodes and gateways! Got %d", atomic.LoadUint32(&counter))
				return
			}
		}
		time.Sleep(1 * time.Second)
	}

	// Remove all nodes and gateways and resign the NDF
	newNdf.Nodes = make([]ndf.Node, 0)
	newNdf.Gateways = make([]ndf.Gateway, 0)
	f.Ndf, err = newNdf.Marshal()
	if err != nil {
		t.Errorf(err.Error())
	}
	privKey, err := testutils.LoadPrivateKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load private key: %v", err)
		t.FailNow()
	}
	err = signature.SignRsa(f, privKey)
	if err != nil {
		t.Fatalf("Failed to sign round info: %v", err)
	}
	// Set up channels that reduce counter by the number of items they receive
	go func() {
		for range removeNode {
			atomic.AddUint32(&counter, ^uint32(0)) // decrement counter
		}
	}()
	go func() {
		for range removeGateway {
			atomic.AddUint32(&counter, ^uint32(0)) // decrement counter
		}
	}()

	// Install the newly-empty NDF
	err = i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Verify channels received the correct amount of information
	timeout = 5
	for {
		if atomic.LoadUint32(&counter) == 0 {
			break
		} else {
			timeout -= 1
			if timeout == 0 {
				t.Errorf("Unable to properly remove nodes and gateways! Got %d", counter)
				return
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// Happy path
func TestInstance_RoundUpdates(t *testing.T) {
	i, _ := setupComm(t)
	nwHealth := make(chan Heartbeat, 10)
	i.SetNetworkHealthChan(nwHealth)

	r := &mixmessages.RoundInfo{
		ID:         2,
		UpdateID:   4,
		State:      uint32(states.COMPLETED),
		Timestamps: []uint64{0, 0, 0, 0, 0},
	}
	err := testutils.SignRoundInfoRsa(r, t)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Set up a function to read channel output
	isFinished := false
	go func() {
		for heartbeat := range nwHealth {
			if !heartbeat.IsRoundComplete {
				t.Errorf("Round should have been complete")
			}
			if heartbeat.HasWaitingRound {
				t.Errorf("Should have had no waiting rounds")
			}
			isFinished = true
			break
		}
	}()

	// Send the round update
	err = i.RoundUpdates([]*mixmessages.RoundInfo{r})

	// Wait for other thread to finish
	for !isFinished {
		time.Sleep(50 * time.Millisecond)
	}
}

// Happy path
func TestInstance_UpdateGroup(t *testing.T) {
	i, f := setupComm(t)
	err := i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf("Unable to initalize group: %+v", err)
	}

	// Update with same values should not cause an error
	err = i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf("Unable to call update group with same values: %+v", err)
	}

}

// Error path: attempt to modify group once already initialized
func TestInstance_UpdateGroup_Error(t *testing.T) {
	i, f := setupComm(t)

	err := i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf("Unable to initalize group: %+v", err)
	}

	badNdf := createBadNdf(t)

	// Update with same values should not cause an error
	err = i.UpdateFullNdf(badNdf)
	if err != nil {
		return
	}

	t.Errorf("Expected error case: Should not be able to update instance's group once initialized!")

}

// Creates a bad ndf
func createBadNdf(t *testing.T) *mixmessages.NDF {
	f := &mixmessages.NDF{}

	badGrp := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	baseNDF := ndf.NetworkDefinition{
		E2E:  badGrp,
		CMIX: badGrp,
	}

	var err error
	f.Ndf, err = baseNDF.Marshal()
	if err != nil {
		t.Errorf("Could not generate serialized ndf: %s", err)
	}
	privKey, err := testutils.LoadPrivateKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load private key: %v", err)
		t.FailNow()
	}

	err = signature.SignRsa(f, privKey)
	if err != nil {
		t.Fatalf("Failed to sign ndf: %v", err)
	}
	return f
}

// Test that a new round update is inputted into the ERS map
func TestInstance_RoundUpdateAddsToERS(t *testing.T) {
	// Get signing certificates
	privKey, err := testutils.LoadPrivateKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load private key: %v", err)
		t.FailNow()
	}
	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
	if err != nil {
		t.Errorf("Could not get rsa key: %s", err)
	}

	// Create a basic testing NDF and sign it
	f := &mixmessages.NDF{}
	f.Ndf = []byte(testutils.ExampleJSON)
	baseNDF := testutils.NDF
	if err != nil {
		t.Errorf("Could not generate serialized ndf: %s", err)
	}
	err = signature.SignRsa(f, privKey)
	if err != nil {
		t.Fatalf("Failed to sign ndf: %v", err)
	}

	// Build the Instance object with an ERS memory map
	testManager := connect.NewManagerTesting(t)
	pc := &connect.ProtoComms{
		Manager: testManager,
	}
	var ers ds.ExternalRoundStorage = &ersMemMap{rounds: make(map[id.Round]*mixmessages.RoundInfo)}
	i, err := NewInstance(pc, baseNDF, baseNDF, ers, 0, false)
	if err != nil {
		t.Error(nil)
	}

	// Add a permissioning host
	_, err = i.comm.AddHost(&id.Permissioning, "0.0.0.0:4200", pub, connect.GetDefaultHostParams())
	if err != nil {
		t.Errorf("Failed to add permissioning host: %+v", err)
	}

	// Build a basic RoundInfo object and sign it
	r := &mixmessages.RoundInfo{
		ID:       2,
		UpdateID: 4,
	}
	err = signature.SignRsa(r, privKey)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Cause a RoundUpdate
	_, err = i.RoundUpdate(r)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Check that the round info was stored correctly
	rr, err := ers.Retrieve(id.Round(r.ID))
	if err != nil {
		t.Errorf(err.Error())
	}
	if rr == nil {
		t.Fatalf("returned round info was nil")
	}
	if rr.ID != r.ID || rr.UpdateID != r.UpdateID {
		t.Errorf("Second returned round and original mismatched IDs")
	}
}

// Happy path
func TestInstance_GetNodeAndGateway(t *testing.T) {
	// Get signing certificates
	privKey, err := testutils.LoadPrivateKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load private key: %v", err)
		t.FailNow()
	}

	// Create a basic testing NDF and sign it
	f := &mixmessages.NDF{}
	f.Ndf = []byte(testutils.ExampleJSON)
	baseNDF := testutils.NDF
	if err != nil {
		t.Errorf("Could not generate serialized ndf: %s", err)
	}
	err = signature.SignRsa(f, privKey)
	if err != nil {
		t.Fatalf("Failed to sign ndf: %v", err)
	}

	// Build the Instance object with an ERS memory map
	testManager := connect.NewManagerTesting(t)
	pc := &connect.ProtoComms{
		Manager: testManager,
	}
	var ers ds.ExternalRoundStorage = &ersMemMap{rounds: make(map[id.Round]*mixmessages.RoundInfo)}
	i, err := NewInstance(pc, baseNDF, baseNDF, ers, 0, false)
	if err != nil {
		t.Error(nil)
	}

	expectedGateway := baseNDF.Gateways[0]
	expectedNode := baseNDF.Nodes[0]
	ngid, err := id.Unmarshal(expectedGateway.ID)
	if err != nil {
		t.Errorf("Could not parse gateway id in NDF: %v", err)
	}
	nodeGw, err := i.GetNodeAndGateway(ngid)
	if err != nil {
		t.Errorf("Failed to get nodeGateway: %v", err)
	}

	if !reflect.DeepEqual(nodeGw.Gateway, expectedGateway) {
		t.Errorf("Unexpected value in gateway."+
			"\n\tExpected: %v\n\tReceived: %v", expectedGateway, nodeGw.Gateway)
	}

	if !reflect.DeepEqual(nodeGw.Node, expectedNode) {
		t.Errorf("Unexpected value in node."+
			"\n\tExpected: %v\n\tReceived: %v", expectedNode, nodeGw.Node)
	}
}
