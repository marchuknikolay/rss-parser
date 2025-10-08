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

			serverPort = 4321
			timeout    = 5 * time.Second
		)

		t.Cleanup(func() {
			os.Clearenv()
		})

		t.Setenv("DB_HOST", dbHost)
		t.Setenv("DB_USER", dbUser)
		t.Setenv("DB_PASSWORD", dbPassword)
		t.Setenv("DB_NAME", dbName)
		t.Setenv("DB_HOST_PORT", strconv.Itoa(dbPort))
		t.Setenv("DB_CONTAINER_PORT", strconv.Itoa(dbPort))

		t.Setenv("SERVER_PORT", strconv.Itoa(serverPort))
		t.Setenv("SERVER_SHUTDOWN_TIMEOUT", timeout.String())
		t.Setenv("SERVER_READ_HEADER_TIMEOUT", timeout.String())

		config, err := New()

		require.NoError(t, err)

		require.Equal(t, dbHost, config.DB.Host)
		require.Equal(t, dbUser, config.DB.User)
		require.Equal(t, dbPassword, config.DB.Password)
		require.Equal(t, dbName, config.DB.Name)
		require.Equal(t, dbPort, config.DB.HostPort)
		require.Equal(t, dbPort, config.DB.ContainerPort)

		require.Equal(t, serverPort, config.Server.Port)
		require.Equal(t, timeout, config.Server.ShutdownTimeout)
		require.Equal(t, timeout, config.Server.ReadHeaderTimeout)
	})

	t.Run("MissingEnvVariables", func(t *testing.T) {
		os.Clearenv()

		config, err := New()

		require.Error(t, err)
		require.Nil(t, config)
	})
}
