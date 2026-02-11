package pluginmsg

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/edition/java/proxy/message"
)

var Plugin = proxy.Plugin{
	Name: "PluginMsg",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log, err := logr.FromContext(ctx)
		if err != nil {
			return err
		}

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
	//p.registerPluginChannels()
	return nil
}

// Register event subscribers
func (p *PluginMsg) registerSubscribers() {
	// Send message on server switch.
	event.Subscribe(p.Event(), 0, p.onServerSwitch)

	// Change the MOTD response.
	//event.Subscribe(p.Event(), 0, pingHandler())

	// Show a boss bar to all players on this proxy.
	//event.Subscribe(p.Event(), 0, p.bossBarDisplay())

	// Listen for plugin messages.
	event.Subscribe(p.Event(), 0, p.onPluginMessage)

	// Listen for plugin messages during login.
	event.Subscribe(p.Event(), 0, p.onServerLoginPluginMessage)

	p.log.Info("Registered plugin message event subscriber")
}

func (p *PluginMsg) registerPluginChannels() {
	// Register a plugin channel for sending messages to the client.
	//p.Proxy().RegisterPluginChannel("my:channel")

	luckId, err := message.NewChannelIdentifier("pluginmsg", "main")
	if err != nil {
		p.log.Error(err, "Failed to create plugin channel identifier")
		return
	}

	p.ChannelRegistrar().Register(luckId)
}

func (p *PluginMsg) onServerLoginPluginMessage(e proxy.ServerLoginPluginMessageEvent) {
	p.log.Info("ServerLoginPluginMessageEvent fired!")

	res := e.Result()

	p.log.Info("Plugin message received during login", "allowed", res.Allowed())
}

func (p *PluginMsg) onPluginMessage(e proxy.PluginMessageEvent) {

	p.log.Info("PluginMessageEvent fired!")

	// Another plugin may have already cancelled the event.
	if !e.Allowed() {
		p.log.Info("Plugin message event already cancelled, not responding")
		return
	}

	// if _, ok := e.Target().(proxy.Player); ok {
	// 	// e.Source() IS a proxy.Player
	// 	// 'player' is now the concrete value
	// 	p.log.Info("Plugin message receieved with target player, do nothing")
	// 	return
	// }

	// e.SetForward(false)

	// if e.Source() == nil {
	// 	p.log.Info("Plugin message received", "source type", "<nil>", "length", len(e.Data()))
	// } else {
	// 	p.log.Info("Plugin message received", "source type", fmt.Sprintf("%T", e.Source()), "length", len(e.Data()))
	// }
}

func (p *PluginMsg) onServerPreConnect(e proxy.ServerPreConnectEvent) {
	p.log.Info("ServerPreConnectEvent fired!")

	//p.registerPluginChannels()
}
