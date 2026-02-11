package pluginmsg

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/edition/java/proxy/message"
)

var Plugin = proxy.Plugin{
	Name: "PluginMsg",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log := logr.FromContextOrDiscard(ctx).WithName("PluginMsg")

		return newPluginMsg(p, &log).init()
	},
}

type PluginMsg struct {
	*proxy.Proxy
	log *logr.Logger
}

func newPluginMsg(proxy *proxy.Proxy, log *logr.Logger) *PluginMsg {
	return &PluginMsg{
		Proxy: proxy,
		log:   log,
	}
}

// initialize the plugin, e.g. register commands and event handlers
func (p *PluginMsg) init() error {
	//p.registerCommands()
	p.registerSubscribers()
	p.registerPluginChannels()
	return nil
}

func (p *PluginMsg) registerPluginChannels() {
	// Register a plugin channel for sending messages to the client.
	//p.Proxy().RegisterPluginChannel("my:channel")
	p.ChannelRegistrar().Register(message.LegacyChannelIdentifier("luckperms:update"))
}

// Register event subscribers
func (p *PluginMsg) registerSubscribers() {
	// Send message on server switch.
	//event.Subscribe(p.Event(), 0, p.onServerSwitch)

	// Change the MOTD response.
	//event.Subscribe(p.Event(), 0, pingHandler())

	// Show a boss bar to all players on this proxy.
	//event.Subscribe(p.Event(), 0, p.bossBarDisplay())

	// Listen for plugin messages.
	event.Subscribe(p.Event(), 0, p.onPluginMessage)
}

func (p *PluginMsg) onPluginMessage(e proxy.PluginMessageEvent) {

	// Another plugin may have already cancelled the event.
	if !e.Allowed() {
		p.log.Info("Plugin message event already cancelled, not responding")
		return
	}

	if e.Source() == nil {
		p.log.Info("Plugin message received", "source type", "<nil>", "length", len(e.Data()))
	} else {
		p.log.Info("Plugin message received", "source type", fmt.Sprintf("%T", e.Source()), "length", len(e.Data()))
	}
}
