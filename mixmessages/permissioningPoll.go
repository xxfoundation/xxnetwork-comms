package mixmessages

// GetActivity gets the state of the node
func (m *PermissioningPoll) GetActivity() uint32 {
	if m != nil {
		return m.NodeState
	}
	return 0
}
