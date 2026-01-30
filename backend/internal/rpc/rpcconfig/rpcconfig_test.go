package rpcconfig_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/luckymaks/bm_backend/backend/internal/rpc/rpcconfig"
)

func TestNewLoader(t *testing.T) {
	t.Parallel()

	t.Run("creates loader with default values", func(t *testing.T) {
		t.Parallel()

		loader, err := rpcconfig.NewLoader()

		require.NoError(t, err)
		require.NotNil(t, loader)
	})
}

func TestLoaderLoad(t *testing.T) {
	t.Parallel()

	t.Run("returns config with environment values", func(t *testing.T) {
		t.Parallel()
		loader, err := rpcconfig.NewLoader()
		require.NoError(t, err)

		cfg, err := loader.Load(context.Background())

		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.Equal(t, "Dev", cfg.Env.DeploymentIdent)
	})
}

func TestEnvConfigDefaults(t *testing.T) {
	t.Parallel()

	t.Run("deployment ident defaults to Dev", func(t *testing.T) {
		t.Parallel()
		loader, err := rpcconfig.NewLoader()
		require.NoError(t, err)

		cfg, err := loader.Load(context.Background())

		require.NoError(t, err)
		require.Equal(t, "Dev", cfg.Env.DeploymentIdent)
	})
}
