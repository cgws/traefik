package zk

import (
	"github.com/abronan/valkeyrie/store"
	"github.com/cgws/traefik/v2/pkg/provider"
	"github.com/cgws/traefik/v2/pkg/provider/kv"
)

var _ provider.Provider = (*Provider)(nil)

// Provider holds configurations of the provider.
type Provider struct {
	kv.Provider
}

// SetDefaults sets the default values.
func (p *Provider) SetDefaults() {
	p.Provider.SetDefaults()
	p.Endpoints = []string{"127.0.0.1:2181"}
}

// Init the provider
func (p *Provider) Init() error {
	return p.Provider.Init(store.ZK, "zookeeper")
}
