package plugin

import "context"

// Options are used as part of a new plugin.
type Options struct {
	Name     string
	Handlers []Handler
	Init     func(context.Context) error
}

type Option func(o *Options)

// WithHandler adds middleware handlers to.
func WithHandler(h ...Handler) Option {
	return func(o *Options) {
		o.Handlers = append(o.Handlers, h...)
	}
}

// WithName defines the name of the plugin.
func WithName(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// WithInit sets the init function.
func WithInit(fn func(context.Context) error) Option {
	return func(o *Options) {
		o.Init = fn
	}
}
