// Code generated by Kitex v0.6.0. DO NOT EDIT.

package gopushservice

import (
	server "github.com/cloudwego/kitex/server"
	kitex_gen "test/kitex-issue/kitex_gen"
)

// NewInvoker creates a server.Invoker with the given handler and options.
func NewInvoker(handler kitex_gen.GoPushService, opts ...server.Option) server.Invoker {
	var options []server.Option

	options = append(options, opts...)

	s := server.NewInvoker(options...)
	if err := s.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	return s
}