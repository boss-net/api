package real

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/boss-net/api/boss-plugin/internal/core/boss_invocation"
)

type NewBossInvocationDaemonPayload struct {
	BaseUrl      string
	CallingKey   string
	WriteTimeout int64
	ReadTimeout  int64
}

func NewBossInvocationDaemon(payload NewBossInvocationDaemonPayload) (boss_invocation.BackwardsInvocation, error) {
	var err error
	invocation := &RealBackwardsInvocation{}
	baseurl, err := url.Parse(payload.BaseUrl)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 120 * time.Second,
			}).Dial,
			IdleConnTimeout: 120 * time.Second,
		},
	}

	invocation.bossInnerApiBaseurl = baseurl
	invocation.client = client
	invocation.bossInnerApiKey = payload.CallingKey
	invocation.writeTimeout = payload.WriteTimeout
	invocation.readTimeout = payload.ReadTimeout

	return invocation, nil
}
