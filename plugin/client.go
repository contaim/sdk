// Copyright (c) 2021 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package plugin

import (
	"context"

	"github.com/contaim/sdk/plugin/proto/v1"
	"github.com/contaim/sdk/protos/hclspec/v1"
	"github.com/contaim/sdk/utils/grpcutil"
)

// BasePluginClient implements the client interface for the remote base plugin, using
// gRPC to communicate with the plugin.
type BasePluginClient struct {
	Client proto.BasePluginServiceClient

	// DoneCtx is closed when the plugin is exiting.
	DoneCtx context.Context
}

func (b *BasePluginClient) daemonPluginInfo(res *proto.PluginInfoResponse) (*PluginInfoResponse, error) {
	var pluginType string
	switch res.DaemonType {
	case proto.DaemonPluginType_DAEMON_PLUGIN_TYPE_DRIVER:
		pluginType = DaemonPluginTypeDriver
	default:
		pluginType = DaemonPluginTypeUnspecified
	}

	return &PluginInfoResponse{
		Runtime:           RuntimeModeDaemon,
		Type:              pluginType,
		PluginApiVersions: res.PluginApiVersions,
		PluginVersion:     res.PluginVersion,
		Name:              res.Name,
	}, nil
}

func (b *BasePluginClient) serverPluginInfo(res *proto.PluginInfoResponse) (*PluginInfoResponse, error) {
	var pluginType string
	switch res.ServerType {
	case proto.ServerPluginType_SERVER_PLUGIN_TYPE_BUILDER:
		pluginType = ServerPluginTypeBuilder
	default:
		pluginType = ServerPluginTypeUnspecified
	}

	return &PluginInfoResponse{
		Runtime:           RuntimeModeServer,
		Type:              pluginType,
		PluginApiVersions: res.PluginApiVersions,
		PluginVersion:     res.PluginVersion,
		Name:              res.Name,
	}, nil
}

func (b *BasePluginClient) PluginInfo() (*PluginInfoResponse, error) {
	res, err := b.Client.PluginInfo(b.DoneCtx, &proto.PluginInfoRequest{})
	if err != nil {
		return nil, grpcutil.HandleGrpcErr(err, b.DoneCtx)
	}

	switch res.Runtime {
	case proto.PluginRuntime_PLUGIN_RUNTIME_SERVER:
		return b.serverPluginInfo(res)
	case proto.PluginRuntime_PLUGIN_RUNTIME_DAEMON:
		return b.daemonPluginInfo(res)
	default:
		return nil, ErrPluginInvalidRuntime
	}
}

func (b *BasePluginClient) ConfigSchema() (*hclspec.Spec, error) {
	res, err := b.Client.ConfigSchema(b.DoneCtx, &proto.ConfigSchemaRequest{})
	if err != nil {
		return nil, grpcutil.HandleGrpcErr(err, b.DoneCtx)
	}

	return res.GetSpec(), nil
}

func (b *BasePluginClient) SetConfig(c *Config) error {
	_, err := b.Client.SetConfig(b.DoneCtx, &proto.SetConfigRequest{
		PluginConfig:        c.PluginConfig,
		DaemonRuntimeConfig: c.DaemonConfig.toProto(),
		ServerRuntimeConfig: c.ServerConfig.toProto(),
		PluginApiVersion:    c.ApiVersion,
	})

	return grpcutil.HandleGrpcErr(err, b.DoneCtx)
}
