package pluginmsg

import (
	"context"
	"math"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/edition/java/proxy/message"
	"go.minekube.com/gate/pkg/util/permission"
)

const (
	Name = "lazygate" // Name represents plugin name.
)

var Plugin = proxy.Plugin{
	Name: "PluginMsg",
	Init: func(ctx context.Context, proxy *proxy.Proxy) error {
		//logger, err := logr.FromContext(ctx)
		// if err != nil {
		// 	return err
		// }

		return newPluginMsg(ctx, proxy).init()
	},
}

type PluginMsg struct {
	ctx   context.Context // Plugin context.
	log   logr.Logger     // Plugin logger.
	proxy *proxy.Proxy    // Gate proxy instance.
}

func newPluginMsg(ctx context.Context, proxy *proxy.Proxy) *PluginMsg {
	return &PluginMsg{
		ctx:   ctx,
		proxy: proxy,
	}
}

// initialize the plugin, e.g. register commands and event handlers
func (p *PluginMsg) init() error {
	p.log = logr.FromContextOrDiscard(p.ctx).WithName("PluginMsg")

	//p.registerCommands()
	p.initPluginChannels()

	if err := p.initHandlers(); err != nil {
		return err
	}
	return nil
}

// initHandlers subscribes event handlers.
func (p *PluginMsg) initHandlers() error {
	eventMgr := p.proxy.Event()

	event.Subscribe(eventMgr, 0, p.onPermissionsSetup)
	p.log.Info("Registered permissions setup event subscriber")

	// Change the MOTD response.
	//event.Subscribe(p.Event(), 0, pingHandler())

	// Show a boss bar to all players on this proxy.
	//event.Subscribe(p.Event(), 0, p.bossBarDisplay())

	// Listen for plugin messages.
	event.Subscribe(eventMgr, math.MaxInt, onPluginMessage(p.log))

	p.log.Info("Registered plugin message event subscriber")

	// Listen for plugin messages during login.
	event.Subscribe(eventMgr, math.MaxInt, p.onServerLoginPluginMessage)

	p.log.Info("Registered server login plugin message event subscriber")

	return nil
}

func (p *PluginMsg) initPluginChannels() {
	// Register a plugin channel for sending messages to the client.
	//p.Proxy().RegisterPluginChannel("my:channel")

	luckId, err := message.NewChannelIdentifier("luckperms", "update")
	if err != nil {
		p.log.Error(err, "Failed to create plugin channel identifier")
		return
	}

	p.proxy.ChannelRegistrar().Register(luckId)
}

func onPluginMessage(logger logr.Logger) func(*proxy.PluginMessageEvent) {

	return func(e *proxy.PluginMessageEvent) {
		logger.Info("[PLUGIN MESSAGE] Received plugin message", "ID", e.Identifier().ID())

		// Another plugin may have already cancelled the event.
		if !e.Allowed() {
			logger.Info("[PLUGIN MESSAGE] Plugin message event already cancelled, not responding")
			return
		}

		// Check the channel and respond to a specific one.
		if e.Identifier().ID() == "luckperms:update" {
			logger.Info("[PLUGIN MESSAGE] Received plugin message on luckperms:update")
			//e.SetForward(false)
		}
	}
}

func (p *PluginMsg) onServerLoginPluginMessage(e *proxy.ServerLoginPluginMessageEvent) {

	p.log.Info("ServerLoginPluginMessageEvent fired!")

	res := e.Result()

	p.log.Info("Plugin message received during login", "allowed", res.Allowed())

}

func (p *PluginMsg) onPermissionsSetup(e *proxy.PermissionsSetupEvent) {

	if player, ok := e.Subject().(proxy.Player); ok {
		// e.Subject() IS a proxy.Player
		// 'player' is now the concrete value
		p.log.Info("Setting up permissions for player", "player", player.Username())

		// type Func func(permission string) TriState
		e.SetFunc(func(permission string) permission.TriState {
			p.log.Info("Permission check", "player", player.Username(), "permission", permission)
			// For demonstration purposes, we allow the "example.permission" permission and deny all others.
			// if permission == "example.permission" {
			// 	return 1 // permission.True
			// }
			if player.Username() == "NerdByNature" {
				return 1 // permission.True
			}
			return 0 // permission.False
		})

	}
}
