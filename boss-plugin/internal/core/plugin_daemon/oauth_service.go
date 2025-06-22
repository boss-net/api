package plugin_daemon

import (
	"github.com/boss-net/api/boss-plugin/internal/core/session_manager"
	"github.com/boss-net/api/boss-plugin/internal/utils/stream"
	"github.com/boss-net/api/boss-plugin/pkg/entities/oauth_entities"
	"github.com/boss-net/api/boss-plugin/pkg/entities/requests"
)

func OAuthGetAuthorizationURL(
	session *session_manager.Session,
	request *requests.RequestOAuthGetAuthorizationURL,
) (*stream.Stream[oauth_entities.OAuthGetAuthorizationURLResult], error) {
	return GenericInvokePlugin[requests.RequestOAuthGetAuthorizationURL, oauth_entities.OAuthGetAuthorizationURLResult](
		session,
		request,
		1,
	)
}

func OAuthGetCredentials(
	session *session_manager.Session,
	request *requests.RequestOAuthGetCredentials,
) (*stream.Stream[oauth_entities.OAuthGetCredentialsResult], error) {
	return GenericInvokePlugin[requests.RequestOAuthGetCredentials, oauth_entities.OAuthGetCredentialsResult](
		session,
		request,
		1,
	)
}
