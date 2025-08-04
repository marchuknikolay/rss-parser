package config

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		const (
			dbHost     = "host"
			dbUser     = "user"
			dbPassword = "password"
			dbName     = "name"
			dbPort     = 1234

			serverPort            = 4321
			serverShutdownTimeout = 5 * time.Second
		)

		t.Cleanup(func() {
			os.Clearenv()
		})

		require.NoError(t, os.Setenv("DB_HOST", dbHost))
		require.NoError(t, os.Setenv("DB_USER", dbUser))
		require.NoError(t, os.Setenv("DB_PASSWORD", dbPassword))
		require.NoError(t, os.Setenv("DB_NAME", dbName))
		require.NoError(t, os.Setenv("DB_HOST_PORT", strconv.Itoa(dbPort)))
		require.NoError(t, os.Setenv("DB_CONTAINER_PORT", strconv.Itoa(dbPort)))

		require.NoError(t, os.Setenv("SERVER_PORT", strconv.Itoa(serverPort)))
		require.NoError(t, os.Setenv("SERVER_SHUTDOWN_TIMEOUT", serverShutdownTimeout.String()))

		config, err := New()

		require.NoError(t, err)

		require.Equal(t, dbHost, config.DB.Host)
		require.Equal(t, dbUser, config.DB.User)
		require.Equal(t, dbPassword, config.DB.Password)
		require.Equal(t, dbName, config.DB.Name)
		require.Equal(t, dbPort, config.DB.HostPort)
		require.Equal(t, dbPort, config.DB.ContainerPort)

		require.Equal(t, serverPort, config.Server.Port)
		require.Equal(t, serverShutdownTimeout, config.Server.ShutdownTimeout)
	})

	t.Run("MissingEnvVariables", func(t *testing.T) {
		os.Clearenv()

		config, err := New()

		require.Error(t, err)
		require.Nil(t, config)
	})
}
