package noop

// NoopFirewallClient performs no operations for exposing or unexposing ports
// It's typically used in development environments where there is no need to expose firewall ports
type NoopFirewallClient struct {
}

// ExposePort performs no operation for this NoopFirewallClient
func (pe NoopFirewallClient) ExposePort(port int32) error {
	return nil
}

// UnexposePort performs no operation for this NoopFirewallClient
func (pe NoopFirewallClient) UnexposePort(port int32) error {
	return nil
}