package connect

// Params object for host creation
type HostParams struct {
	MaxRetries  uint32
	AuthEnabled bool
}

// Get default set of host params
func GetDefaultHostParams() HostParams {
	return HostParams{
		MaxRetries:  100,
		AuthEnabled: true,
	}
}
