package gossip

// Passed into NewGossip to specify how Gossip messages will be handled
// the byte slice will be used to pass in a merkle tree and signature on the
// trees root for multi-part gossips sent over streaming when streaming is
// implemented. Ignore it for non streaming implementations.
type Receiver func(*GossipMsg) error

// Passed into NewGossip to specify how Gossip message signatures will be verified
type SignatureVerification func(*GossipMsg, []byte) error
