package builder

import (
	"errors"
	"testing"

	"github.com/BeInBloom/big-deal/services/order/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBuildWithInMemoryRepo(t *testing.T) {
	t.Setenv("ORDER_HTTP_ADDR", "127.0.0.1:18081")
	t.Setenv("ORDER_REPO_KIND", "hash_map")

	server, closeResources, err := Build(loadConfig(t))
	require.NoError(t, err)

	require.Equal(t, "127.0.0.1:18081", server.Addr)
	require.NotNil(t, server.Handler)
	require.NoError(t, closeResources())
}

func TestBuildReleasesResourcesOnFailure(t *testing.T) {
	t.Setenv("ORDER_REPO_KIND", "postgres")
	t.Setenv("ORDER_REPO_POSTGRES_URL", "postgres://order:secret@127.0.0.1:1/orders?sslmode=disable")

	b := &builder{cfg: loadConfig(t)}

	_, err := b.httpServer()
	require.Error(t, err)
	require.Len(t, b.closers, 1)
	require.NoError(t, b.close())

	server, closeResources, err := Build(loadConfig(t))
	require.Error(t, err)
	require.Nil(t, server)
	require.Nil(t, closeResources)
}

func TestBuilderCloseRunsClosersInReverseOrder(t *testing.T) {
	var closed []string
	errFirst := errors.New("first")

	b := &builder{}
	b.onClose(func() error {
		closed = append(closed, "first")
		return errFirst
	})
	b.onClose(func() error {
		closed = append(closed, "second")
		return nil
	})

	err := b.close()

	require.Equal(t, []string{"second", "first"}, closed)
	require.ErrorIs(t, err, errFirst)
}

func loadConfig(t *testing.T) config.Config {
	t.Helper()

	cfg, err := config.LoadFromEnv()
	require.NoError(t, err)

	return cfg
}
