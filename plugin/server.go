// Copyright (c) 2021 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package plugin

import (
	"context"
	"fmt"

	"github.com/contaim/sdk/plugin/proto/v1"
	"github.com/hashicorp/go-plugin"
)

type basePluginServer struct {
	proto.UnimplementedBasePluginServiceServer

	broker *plugin.GRPCBroker
	impl   BasePlugin
}

func (b *basePluginServer) DaemonPluginInfo(res *PluginInfoResponse) (*proto.PluginInfoResponse, error) {
	var pluginType proto.DaemonPluginType
	switch res.Type {
	case DaemonPluginTypeDriver:
		pluginType = proto.DaemonPluginType_DAEMON_PLUGIN_TYPE_DRIVER
	default:
		return nil, fmt.Errorf("plugin reports an unspcified type: %q", res.Type)
	}

	return &proto.PluginInfoResponse{
		Runtime:           proto.PluginRuntime_PLUGIN_RUNTIME_DAEMON,
		DaemonType:        pluginType,
		ServerType:        proto.ServerPluginType_SERVER_PLUGIN_TYPE_UNSPECIFIED,
		PluginApiVersions: res.PluginApiVersions,
		PluginVersion:     res.PluginVersion,
		Name:              res.Name,
	}, nil
}

func (b *basePluginServer) serverPluginInfo(res *PluginInfoResponse) (*proto.PluginInfoResponse, error) {
	var pluginType proto.ServerPluginType
	switch res.Type {
	case ServerPluginTypeBuilder:
		pluginType = proto.ServerPluginType_SERVER_PLUGIN_TYPE_BUILDER
	default:
		return nil, fmt.Errorf("plugin reports an unknown type: %q", res.Type)
	}

	return &proto.PluginInfoResponse{
		Runtime:           proto.PluginRuntime_PLUGIN_RUNTIME_SERVER,
		DaemonType:        proto.DaemonPluginType_DAEMON_PLUGIN_TYPE_UNSPECIFIED,
		ServerType:        pluginType,
		PluginApiVersions: res.PluginApiVersions,
		PluginVersion:     res.PluginVersion,
		Name:              res.Name,
	}, nil
}

func (b *basePluginServer) PluginInfo(context.Context, *proto.PluginInfoRequest) (*proto.PluginInfoResponse, error) {
	res, err := b.impl.PluginInfo()
	if err != nil {
		return nil, err
	}

	switch res.Runtime {
	case RuntimeModeDaemon:
		return b.DaemonPluginInfo(res)
	case RuntimeModeServer:
		return b.serverPluginInfo(res)
	default:
		return nil, ErrPluginInvalidRuntime
	}
}

func (b *basePluginServer) ConfigSchema(context.Context, *proto.ConfigSchemaRequest) (*proto.ConfigSchemaResponse, error) {
	spec, err := b.impl.ConfigSchema()
	if err != nil {
		return nil, err
	}

	return &proto.ConfigSchemaResponse{
		Spec: spec,
	}, nil
}

func (b *basePluginServer) SetConfig(ctx context.Context, req *proto.SetConfigRequest) (*proto.SetConfigResponse, error) {
	info, err := b.impl.PluginInfo()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ApiVersion:   req.GetPluginApiVersion(),
		PluginConfig: req.GetPluginConfig(),
	}

	// when the plugin is an Daemon plugin
	if info.Runtime == RuntimeModeDaemon {
		daemonConfig := protoToDaemonConfig(req.GetDaemonRuntimeConfig())

		if daemonConfig != nil {
			switch info.Type {
			case DaemonPluginTypeDriver:
				cfg.DaemonConfig = &DaemonConfig{}
			}
		}
	}

	// when the plugin is a server plugin
	if info.Runtime == RuntimeModeServer {
		serverConfig := protoToServerConfig(req.GetServerRuntimeConfig())

		if serverConfig != nil {
			switch info.Type {
			case ServerPluginTypeBuilder:
				cfg.ServerConfig = &ServerConfig{}
			}
		}
	}

	if err := b.impl.SetConfig(cfg); err != nil {
		return nil, fmt.Errorf("impl.SetConfig: %w", err)
	}

	return &proto.SetConfigResponse{}, nil
}
