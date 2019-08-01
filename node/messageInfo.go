package node

type MessageInfo struct {
	SenderId       string
	Address        string
	Signature      []byte
	ValidSignature bool
}
