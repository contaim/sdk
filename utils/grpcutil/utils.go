package grpcutil

import (
	"context"
	"time"

	"github.com/contaim/sdk/plugin/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errorWaitTime = 5 * time.Second
)

// HandleGrpcErr handles various errors returned from gRPC-based remote plugins. This method
// will handle detecting if the plugin has shutdown and gracefully handle waiting for a proper
// shutdown if the gRPC remote plugin is unavailable.
//
// This method may return ErrPluginShutdown instead of the original error if the plugin has
// actually shutdown.
func HandleGrpcErr(err error, pluginCtx context.Context) error {
	if err == nil {
		return nil
	}

	if errStatus := status.FromContextError(pluginCtx.Err()); errStatus.Code() == codes.Canceled {
		select {
		case <-pluginCtx.Done():
			err = errs.ErrPluginShutdown
		default:
		}
	}

	if errStatus, ok := status.FromError(err); ok && errStatus.Code() == codes.Unavailable {
		select {
		case <-pluginCtx.Done():
			err = errs.ErrPluginShutdown

		// Give the plugin the length of errorWaitTime to report that it was shutdown. If the wait time
		// has been hit, continue returning the error. This should never reach the timeout if the plugin
		// is truly shutting down.
		case <-time.After(errorWaitTime):
		}

		return err
	}

	return err
}
