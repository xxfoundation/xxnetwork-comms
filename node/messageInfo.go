package node

// MessageInfo contains information about a comm to be passed into every interface callback
// Specifically, contains sender ID and network address, signature and signature validity
type MessageInfo struct {
	SenderId       string
	Address        string
	Signature      []byte
	ValidSignature bool
}
