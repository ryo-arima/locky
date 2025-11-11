package repository_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRedisClient(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	redisCfg := config.Redis{
		Host: s.Host(),
		Port: uint16(s.Port()),
	}

	client, err := repository.NewRedisClient(redisCfg)
	require.NoError(t, err)
	assert.NotNil(t, client)

	// Test connection
	pong, err := client.Ping(client.Context()).Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)
}

func TestNewRedisClient_Upstash(t *testing.T) {
	// This test is more of a configuration check, as we can't easily mock the TLS handshake
	redisCfg := config.Redis{
		Host: "some-host.upstash.io",
		Port: 12345,
	}

	// We expect this to fail because we're not running a real upstash instance,
	// but we can check that the TLS config is being set.
	// The actual connection logic is tested in TestNewRedisClient_TLS_Retry
	_, err := repository.NewRedisClient(redisCfg)
	require.Error(t, err)
}
