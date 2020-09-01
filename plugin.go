package plugin

import (
	"context"
	"net/http"
)

type Handler func(handler http.Handler) http.Handler

type Plugin interface {
	Init(ctx context.Context) error
	Handler() Handler
	String() string
}

type Manager interface {
	Plugins(...PluginOption) []Plugin
	Register(Plugin, ...PluginOption) error
}

type PluginOptions struct {
	Module string
}

type PluginOption func(o *PluginOptions)

// Module will scope the plugin to a specific module, e.g. the "api".
func Module(m string) PluginOption {
	return func(o *PluginOptions) {
		o.Module = m
	}
}

type plugin struct {
	opts    Options
	init    func(ctx *context.Context) error
	handler Handler
}

func (p *plugin) Handler() Handler {
	return p.handler
}

func (p *plugin) Init(ctx context.Context) error {
	return p.opts.Init(ctx)
}

func (p *plugin) String() string {
	return p.opts.Name
}

func newPlugin(opts ...Option) Plugin {
	options := Options{
		Name: "default",
		Init: func(ctx context.Context) error { return nil },
	}

	for _, o := range opts {
		o(&options)
	}

	handler := func(hdlr http.Handler) http.Handler {
		for _, h := range options.Handlers {
			hdlr = h(hdlr)
		}

		return hdlr
	}

	return &plugin{
		opts:    options,
		handler: handler,
	}
}

// Plugins lists the global plugins.
func Plugins(opts ...PluginOption) []Plugin {
	return defaultManager.Plugins(opts...)
}

// Register registers a global plugins.
func Register(plugin Plugin, opts ...PluginOption) error {
	return defaultManager.Register(plugin, opts...)
}

// IsRegistered check plugin whether registered global.
// Notice plugin is not check whether is nil.
func IsRegistered(plugin Plugin, opts ...PluginOption) bool {
	return defaultManager.isRegistered(plugin, opts...)
}

// NewManager creates a new plugin manager.
func NewManager() Manager {
	return newManager()
}

// NewPlugin makes it easy to create a new plugin.
func NewPlugin(opts ...Option) Plugin {
	return newPlugin(opts...)
}
