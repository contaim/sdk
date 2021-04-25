// Copyright (c) 2021 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package plugin

import (
	"context"

	"github.com/contaim/sdk/plugin/proto/v1"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const (
	// PluginTypeBase is a plugin that strictly implements the base
	// requirements. This is basically a useless plugin and a warning
	// will be emitted if it is ever loaded.
	PluginTypeBase = "base"

	// DaemonPluginTypeUnspecified is reported when the plugin does not properly
	// report a plugin type. This plugin will not operate.
	DaemonPluginTypeUnspecified = "daemon.unspecified"

	// DaemonPluginTypeDriver is a plugin that implements the driver plugin
	// interface. This plugin will manage containers and networks on the
	// supported container runtime.
	DaemonPluginTypeDriver = "daemon.driver"

	// ServerPluginTypeUnspecified is reported when the plugin on the server
	// does not properly report a plugin type. This plugin will not operate.
	ServerPluginTypeUnspecified = "server.unspecified"

	// PluginTypeBuilder is a plugin that implements the builder plugin
	// interface. This plugin type manages building/provisioning servers
	// at various cloud server providers.
	ServerPluginTypeBuilder = "server.builder"
)

const (
	// RuntimeModeDaemon is reported by plugins that should be executed by a daemon.
	RuntimeModeDaemon = "daemon"

	// RuntimeModeServer is reported by plugins that should be executed by a server.
	RuntimeModeServer = "server"
)

var (
	// DaemonPluginHandshake is the shared handshake for plugins that should be executed
	// by the agent, such as driver plugins.
	DaemonPluginHandshake = plugin.HandshakeConfig{
		ProtocolVersion:  2,
		MagicCookieKey:   "CONTAIM_DARMON_MAGIC_COOKIE",
		MagicCookieValue: "VbAU6gtNdef9Xpeid7rJMpZn67uqJj8F6oiLBSmMK6pGDhICr8Se6xlx1H3kr4Jq8yx9dIo2Qs9AGmg5",
	}

	// ServerPluginHandshake is the shared handshake for plugins that should be executed
	// by the server, such as builder plugins.
	ServerPluginHandshake = plugin.HandshakeConfig{
		ProtocolVersion:  2,
		MagicCookieKey:   "CONTAIM_SERVER_MAGIC_COOKIE",
		MagicCookieValue: "Fr8JySTs186bMm2PaOiG1fIxVJOsn5PE9N0zGhX2yaaFJST5gVDQVvUVyg6IbfnCj1704weTPeQR7m9S",
	}
)

type PluginBase struct {
	plugin.NetRPCUnsupportedPlugin
	Impl BasePlugin
}

func (p *PluginBase) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	proto.RegisterBasePluginServiceServer(server, &basePluginServer{
		impl:   p.Impl,
		broker: broker,
	})

	return nil
}

func (p *PluginBase) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, clientConn *grpc.ClientConn) (interface{}, error) {
	return &BasePluginClient{
		Client:  proto.NewBasePluginServiceClient(clientConn),
		DoneCtx: ctx,
	}, nil
}
