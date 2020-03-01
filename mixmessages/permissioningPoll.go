package mixmessages

// GetState gets the state of the node
func (m *PermissioningPoll) GetState() uint32 {
	if m != nil {
		return m.Activity
	}
	return 0
}
