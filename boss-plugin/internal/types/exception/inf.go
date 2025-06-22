package exception

import "github.com/boss-net/api/boss-plugin/pkg/entities"

type PluginDaemonError interface {
	error

	ToResponse() *entities.Response
}
