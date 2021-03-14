package providers

// Provider is the default interface for all providers
type Provider interface {
	UpdateIP(subdomain, ip string) error
}
