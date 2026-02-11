package pluginmsg

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "PluginMsg",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log := logr.FromContextOrDiscard(ctx).WithName("pluginMsg")
		log.Info("Hello from PluginMsg plugin!")

		pl := &plugin{proxy: p, log: &log}
		event.Subscribe(p.Event(), 0, pl.onPluginMessage)
		return nil
	},
}

type plugin struct {
	proxy *proxy.Proxy
	log   *logr.Logger
}

func (p *plugin) onPluginMessage(e proxy.PluginMessageEvent) {

	if e.Source() == nil {
		p.log.Info("Plugin message received", "source type", "<nil>", "length", len(e.Data()))
	} else {
		p.log.Info("Plugin message received", "source type", fmt.Sprintf("%T", e.Source()), "length", len(e.Data()))
	}

	// Another plugin may have already cancelled the event.
	if !e.Allowed() {
		p.log.Info("Plugin message event already cancelled, not responding")
		return
	}
}
