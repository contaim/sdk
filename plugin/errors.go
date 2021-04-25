package plugin

import "errors"

const (
	errMsgPluginShutdown       = "plugin has shut down"
	errMsgPluginInvalidRuntime = "plugin reports invalid runtime mode"
)

var (
	// ErrPluginShutdown is returned when an action is attempted on a plugin that
	// has already shutdown.
	ErrPluginShutdown = errors.New(errMsgPluginShutdown)

	// ErrPluginInvalidRuntime is returned when a plugin reports an invalid runtime
	// mode. This is a fatal error and the plugin will not operate.
	ErrPluginInvalidRuntime = errors.New(errMsgPluginInvalidRuntime)
)
