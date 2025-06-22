package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/boss-net/api/boss-plugin/internal/core/boss_invocation"
	"github.com/boss-net/api/boss-plugin/internal/core/plugin_daemon"
	"github.com/boss-net/api/boss-plugin/internal/core/plugin_daemon/access_types"
	"github.com/boss-net/api/boss-plugin/internal/core/plugin_manager"
	"github.com/boss-net/api/boss-plugin/internal/core/session_manager"
	"github.com/boss-net/api/boss-plugin/internal/db"
	"github.com/boss-net/api/boss-plugin/internal/service/install_service"
	"github.com/boss-net/api/boss-plugin/internal/types/exception"
	"github.com/boss-net/api/boss-plugin/internal/types/models"
	"github.com/boss-net/api/boss-plugin/internal/utils/encryption"
	"github.com/boss-net/api/boss-plugin/internal/utils/routine"
	"github.com/boss-net/api/boss-plugin/pkg/entities"
	"github.com/boss-net/api/boss-plugin/pkg/entities/endpoint_entities"
	"github.com/boss-net/api/boss-plugin/pkg/entities/plugin_entities"
	"github.com/boss-net/api/boss-plugin/pkg/entities/requests"
)

func copyRequest(req *http.Request, hookId string, path string) (*bytes.Buffer, error) {
	newReq := req.Clone(context.Background())
	// get query params
	queryParams := req.URL.Query()

	// replace path with endpoint path
	newReq.URL.Path = path
	// set query params
	newReq.URL.RawQuery = queryParams.Encode()

	// read request body until complete, max 10MB
	body, err := io.ReadAll(io.LimitReader(req.Body, 10*1024*1024))
	if err != nil {
		return nil, err
	}

	// replace with a new reader
	newReq.Body = io.NopCloser(bytes.NewReader(body))
	newReq.ContentLength = int64(len(body))
	newReq.TransferEncoding = nil

	// remove ip traces for security
	newReq.Header.Del("X-Forwarded-For")
	newReq.Header.Del("X-Real-IP")
	newReq.Header.Del("X-Forwarded")
	newReq.Header.Del("X-Original-Forwarded-For")
	newReq.Header.Del("X-Original-Url")
	newReq.Header.Del("X-Original-Host")

	// check if X-Original-Host is set
	if originalHost := req.Header.Get(endpoint_entities.HeaderXOriginalHost); originalHost != "" {
		// replace host with original host
		newReq.Host = originalHost
	}

	// setup hook id to request
	newReq.Header.Set("Boss-Hook-Id", hookId)
	// check if Boss-Hook-Url is set
	if url := req.Header.Get("Boss-Hook-Url"); url == "" {
		newReq.Header.Set(
			"Boss-Hook-Url",
			fmt.Sprintf("http://%s/e/%s%s", newReq.Host, hookId, path),
		)
	}

	var buffer bytes.Buffer
	err = newReq.Write(&buffer)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}

func Endpoint(
	ctx *gin.Context,
	endpoint *models.Endpoint,
	pluginInstallation *models.PluginInstallation,
	maxExecutionTime time.Duration,
	path string,
) {
	if !endpoint.Enabled {
		ctx.JSON(404, exception.NotFoundError(errors.New("endpoint not found")).ToResponse())
		return
	}

	buffer, err := copyRequest(ctx.Request, endpoint.HookID, path)
	if err != nil {
		ctx.JSON(500, exception.InternalServerError(err).ToResponse())
		return
	}

	identifier, err := plugin_entities.NewPluginUniqueIdentifier(pluginInstallation.PluginUniqueIdentifier)
	if err != nil {
		ctx.JSON(400, exception.UniqueIdentifierError(err).ToResponse())
		return
	}

	// fetch plugin
	manager := plugin_manager.Manager()
	runtime, err := manager.Get(identifier)
	if err != nil {
		ctx.JSON(404, exception.ErrPluginNotFound().ToResponse())
		return
	}

	// fetch endpoint declaration
	endpointDeclaration := runtime.Configuration().Endpoint
	if endpointDeclaration == nil {
		ctx.JSON(404, exception.ErrPluginNotFound().ToResponse())
		return
	}

	// decrypt settings
	settings, err := manager.BackwardsInvocation().InvokeEncrypt(&boss_invocation.InvokeEncryptRequest{
		BaseInvokeBossRequest: boss_invocation.BaseInvokeBossRequest{
			TenantId: endpoint.TenantID,
			UserId:   "",
			Type:     boss_invocation.INVOKE_TYPE_ENCRYPT,
		},
		InvokeEncryptSchema: boss_invocation.InvokeEncryptSchema{
			Opt:       boss_invocation.ENCRYPT_OPT_DECRYPT,
			Namespace: boss_invocation.ENCRYPT_NAMESPACE_ENDPOINT,
			Identity:  endpoint.ID,
			Data:      endpoint.Settings,
			Config:    endpointDeclaration.Settings,
		},
	})

	if err != nil {
		ctx.JSON(500, exception.InternalServerError(err).ToResponse())
		return
	}

	session := session_manager.NewSession(
		session_manager.NewSessionPayload{
			TenantID:               endpoint.TenantID,
			UserID:                 "",
			PluginUniqueIdentifier: identifier,
			ClusterID:              ctx.GetString("cluster_id"),
			InvokeFrom:             access_types.PLUGIN_ACCESS_TYPE_ENDPOINT,
			Action:                 access_types.PLUGIN_ACCESS_ACTION_INVOKE_ENDPOINT,
			Declaration:            runtime.Configuration(),
			BackwardsInvocation:    manager.BackwardsInvocation(),
			IgnoreCache:            false,
			EndpointID:             &endpoint.ID,
		},
	)
	defer session.Close(session_manager.CloseSessionPayload{
		IgnoreCache: false,
	})

	session.BindRuntime(runtime)

	statusCode, headers, response, err := plugin_daemon.InvokeEndpoint(
		session, &requests.RequestInvokeEndpoint{
			RawHttpRequest: hex.EncodeToString(buffer.Bytes()),
			Settings:       settings,
		},
	)
	if err != nil {
		ctx.JSON(500, exception.InternalServerError(err).ToResponse())
		return
	}
	defer response.Close()

	done := make(chan bool)
	closed := new(int32)

	ctx.Status(statusCode)
	for k, v := range *headers {
		if len(v) > 0 {
			ctx.Writer.Header().Set(k, v[0])
		}
	}

	close := func() {
		if atomic.CompareAndSwapInt32(closed, 0, 1) {
			close(done)
		}
	}
	defer close()

	routine.Submit(map[string]string{
		"module":   "service",
		"function": "Endpoint",
	}, func() {
		defer close()
		for response.Next() {
			chunk, err := response.Read()
			if err != nil {
				ctx.Writer.Write([]byte(err.Error()))
				ctx.Writer.Flush()
				return
			}
			ctx.Writer.Write(chunk)
			ctx.Writer.Flush()
		}
	})

	select {
	case <-ctx.Writer.CloseNotify():
	case <-done:
	case <-time.After(maxExecutionTime):
		ctx.JSON(500, exception.InternalServerError(errors.New("killed by timeout")).ToResponse())
	}
}

func EnableEndpoint(endpoint_id string, tenant_id string) *entities.Response {

	if err := install_service.EnabledEndpoint(endpoint_id, tenant_id); err != nil {
		return exception.InternalServerError(errors.New("failed to enable endpoint")).ToResponse()
	}

	return entities.NewSuccessResponse(true)
}

func DisableEndpoint(endpoint_id string, tenant_id string) *entities.Response {

	if err := install_service.DisabledEndpoint(endpoint_id, tenant_id); err != nil {
		return exception.InternalServerError(errors.New("failed to disable endpoint")).ToResponse()
	}

	return entities.NewSuccessResponse(true)
}

func ListEndpoints(tenant_id string, page int, page_size int) *entities.Response {
	endpoints, err := db.GetAll[models.Endpoint](
		db.Equal("tenant_id", tenant_id),
		db.OrderBy("created_at", true),
		db.Page(page, page_size),
	)
	if err != nil {
		return exception.InternalServerError(fmt.Errorf("failed to list endpoints: %v", err)).ToResponse()
	}

	manager := plugin_manager.Manager()
	if manager == nil {
		return exception.InternalServerError(errors.New("failed to get plugin manager")).ToResponse()
	}

	// decrypt settings
	for i, endpoint := range endpoints {
		pluginInstallation, err := db.GetOne[models.PluginInstallation](
			db.Equal("plugin_id", endpoint.PluginID),
			db.Equal("tenant_id", tenant_id),
		)
		if err != nil {
			// use empty settings and declaration for uninstalled plugins
			endpoint.Settings = map[string]any{}
			endpoint.Declaration = &plugin_entities.EndpointProviderDeclaration{
				Settings:      []plugin_entities.ProviderConfig{},
				Endpoints:     []plugin_entities.EndpointDeclaration{},
				EndpointFiles: []string{},
			}
			endpoints[i] = endpoint
			continue
		}

		pluginUniqueIdentifier, err := plugin_entities.NewPluginUniqueIdentifier(
			pluginInstallation.PluginUniqueIdentifier,
		)
		if err != nil {
			return exception.UniqueIdentifierError(
				fmt.Errorf("failed to parse plugin unique identifier: %v", err),
			).ToResponse()
		}

		pluginDeclaration, err := manager.GetDeclaration(
			pluginUniqueIdentifier,
			tenant_id,
			plugin_entities.PluginRuntimeType(pluginInstallation.RuntimeType),
		)
		if err != nil {
			return exception.InternalServerError(
				fmt.Errorf("failed to get plugin declaration: %v", err),
			).ToResponse()
		}

		if pluginDeclaration.Endpoint == nil {
			return exception.NotFoundError(errors.New("plugin does not have an endpoint")).ToResponse()
		}

		decryptedSettings, err := manager.BackwardsInvocation().InvokeEncrypt(&boss_invocation.InvokeEncryptRequest{
			BaseInvokeBossRequest: boss_invocation.BaseInvokeBossRequest{
				TenantId: tenant_id,
				UserId:   "",
				Type:     boss_invocation.INVOKE_TYPE_ENCRYPT,
			},
			InvokeEncryptSchema: boss_invocation.InvokeEncryptSchema{
				Opt:       boss_invocation.ENCRYPT_OPT_DECRYPT,
				Namespace: boss_invocation.ENCRYPT_NAMESPACE_ENDPOINT,
				Identity:  endpoint.ID,
				Data:      endpoint.Settings,
				Config:    pluginDeclaration.Endpoint.Settings,
			},
		})
		if err != nil {
			return exception.InternalServerError(
				fmt.Errorf("failed to decrypt settings: %v", err),
			).ToResponse()
		}

		// mask settings
		decryptedSettings = encryption.MaskConfigCredentials(decryptedSettings, pluginDeclaration.Endpoint.Settings)

		endpoint.Settings = decryptedSettings
		endpoint.Declaration = pluginDeclaration.Endpoint

		endpoints[i] = endpoint
	}

	return entities.NewSuccessResponse(endpoints)
}

func ListPluginEndpoints(tenant_id string, plugin_id string, page int, page_size int) *entities.Response {
	endpoints, err := db.GetAll[models.Endpoint](
		db.Equal("plugin_id", plugin_id),
		db.Equal("tenant_id", tenant_id),
		db.OrderBy("created_at", true),
		db.Page(page, page_size),
	)
	if err != nil {
		return exception.InternalServerError(
			fmt.Errorf("failed to list endpoints: %v", err),
		).ToResponse()
	}

	manager := plugin_manager.Manager()
	if manager == nil {
		return exception.InternalServerError(
			errors.New("failed to get plugin manager"),
		).ToResponse()
	}

	// decrypt settings
	for i, endpoint := range endpoints {
		// get installation
		pluginInstallation, err := db.GetOne[models.PluginInstallation](
			db.Equal("plugin_id", plugin_id),
			db.Equal("tenant_id", tenant_id),
		)
		if err != nil {
			return exception.NotFoundError(
				fmt.Errorf("failed to find plugin installation: %v", err),
			).ToResponse()
		}

		pluginUniqueIdentifier, err := plugin_entities.NewPluginUniqueIdentifier(
			pluginInstallation.PluginUniqueIdentifier,
		)

		if err != nil {
			return exception.UniqueIdentifierError(
				fmt.Errorf("failed to parse plugin unique identifier: %v", err),
			).ToResponse()
		}

		pluginDeclaration, err := manager.GetDeclaration(
			pluginUniqueIdentifier,
			tenant_id,
			plugin_entities.PluginRuntimeType(pluginInstallation.RuntimeType),
		)
		if err != nil {
			return exception.InternalServerError(
				fmt.Errorf("failed to get plugin declaration: %v", err),
			).ToResponse()
		}

		decryptedSettings, err := manager.BackwardsInvocation().InvokeEncrypt(&boss_invocation.InvokeEncryptRequest{
			BaseInvokeBossRequest: boss_invocation.BaseInvokeBossRequest{
				TenantId: tenant_id,
				UserId:   "",
				Type:     boss_invocation.INVOKE_TYPE_ENCRYPT,
			},
			InvokeEncryptSchema: boss_invocation.InvokeEncryptSchema{
				Opt:       boss_invocation.ENCRYPT_OPT_DECRYPT,
				Namespace: boss_invocation.ENCRYPT_NAMESPACE_ENDPOINT,
				Identity:  endpoint.ID,
				Data:      endpoint.Settings,
				Config:    pluginDeclaration.Endpoint.Settings,
			},
		})
		if err != nil {
			return exception.InternalServerError(
				fmt.Errorf("failed to decrypt settings: %v", err),
			).ToResponse()
		}

		// mask settings
		decryptedSettings = encryption.MaskConfigCredentials(decryptedSettings, pluginDeclaration.Endpoint.Settings)

		endpoint.Settings = decryptedSettings
		endpoint.Declaration = pluginDeclaration.Endpoint

		endpoints[i] = endpoint
	}

	return entities.NewSuccessResponse(endpoints)
}
