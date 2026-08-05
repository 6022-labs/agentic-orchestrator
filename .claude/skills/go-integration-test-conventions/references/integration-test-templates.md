# Integration Test Templates

Use these templates to keep integration test structure consistent across projects.
Replace all `<...>` placeholders with your repository package names, helper names, and types.

## 1) PostgreSQL repository integration template

```go
package <repository_package>_test

import (
	"testing"

	dbclient "<module>/src/<db_client_package>"
	repositories "<module>/src/<repository_package_path>"
	settings "<module>/src/<settings_package_path>"
	"<module>/tests/<mocks_package_path>/metrics_mocks"
	test_helpers "<module>/tests/<shared_test_helpers_package>"
	"<module>/tests/<shared_test_helpers_package>/container_builders"
	"<module>/tests/<shared_test_helpers_package>/model_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type When<Behavior>TestingSuite struct {
	dbName string
	sut    *repositories.<RepositoryType>
}

func When<Behavior>BeforeEach(t *testing.T) *When<Behavior>TestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	migrationFolderPath, err := test_helpers.FindMigrationFolderPath()
	if err != nil {
		panic(err)
	}

	dbClient, err := dbclient.New<DbClientType>(
		connectionString,
		zap.NewNop(),
		&settings.<MigrationsSettingsType>{
			<MigrationsFolderPathField>: *migrationFolderPath,
		},
	)
	if err != nil {
		panic(err)
	}
	dbClient.Migrate()

	mockController := gomock.NewController(t)
	metricsMock := metrics_mocks.NewMockDbQueryMetricsInterface(mockController)
	metricsMock.EXPECT().Record(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	sut := repositories.New<RepositoryType>(dbClient, metricsMock)

	return &When<Behavior>TestingSuite{
		dbName: dbName,
		sut:    sut,
	}
}

func When<Behavior>AfterEach(suite *When<Behavior>TestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhen<Behavior>(t *testing.T) {
	t.Parallel()

	t.Run("Given <context>", func(t *testing.T) {
		t.Parallel()

		t.Run("Should <expected result>", func(t *testing.T) {
			t.Parallel()

			suite := When<Behavior>BeforeEach(t)
			defer When<Behavior>AfterEach(suite)

			entity := model_builders.New<Entity>Builder().Build()
			result, err := suite.sut.<Action>(entity)

			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	})
}
```

## 2) Redis integration template

```go
package services

import (
	"testing"

	services "<module>/src/<redis_service_package_path>"
	"<module>/tests/<mocks_package_path>/metrics_mocks"
	"<module>/tests/<shared_test_helpers_package>/container_builders"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type When<Behavior>TestingSuite struct {
	redisClient *redis.Client
	sut         *services.<RedisStorageType>
}

func When<Behavior>BeforeEach(t *testing.T) *When<Behavior>TestingSuite {
	redisClient := container_builders.GetRedisConnection()

	mockController := gomock.NewController(t)
	cacheMetricsMock := metrics_mocks.NewMockCacheMetricsInterface(mockController)
	cacheMetricsMock.EXPECT().RecordRead(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	cacheMetricsMock.EXPECT().RecordWrite(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	sut := services.New<RedisStorageType>(redisClient, cacheMetricsMock)

	return &When<Behavior>TestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func When<Behavior>AfterEach(suite *When<Behavior>TestingSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhen<Behavior>(t *testing.T) {
	t.Parallel()

	t.Run("Given <context>", func(t *testing.T) {
		t.Parallel()

		t.Run("Should <expected result>", func(t *testing.T) {
			t.Parallel()

			suite := When<Behavior>BeforeEach(t)
			defer When<Behavior>AfterEach(suite)

			err := suite.sut.<Action>(<args>)
			assert.NoError(t, err)
		})
	})
}
```

## 3) Factory/infrastructure integration template

```go
package <factory_package>_test

import (
	"testing"

	"<module>/src/<domain>/factories"
	"<module>/tests/<shared_test_helpers_package>/container_builders"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type When<Behavior>TestingSuite struct {
	dbName string
	sut    *factories.<FactoryType>
}

func When<Behavior>BeforeEach() *When<Behavior>TestingSuite {
	dbName, connectionString := container_builders.CreatePostgreSQLDatabase()

	sut := factories.New<FactoryType>(
		zap.NewNop(),
		<settings using connectionString>,
	)

	return &When<Behavior>TestingSuite{
		dbName: dbName,
		sut:    sut,
	}
}

func When<Behavior>AfterEach(suite *When<Behavior>TestingSuite) {
	container_builders.DestroyPostgreSQLDatabase(suite.dbName)
}

func TestWhen<Behavior>(t *testing.T) {
	t.Parallel()

	t.Run("Should <expected result>", func(t *testing.T) {
		t.Parallel()

		suite := When<Behavior>BeforeEach()
		defer When<Behavior>AfterEach(suite)

		result, err := suite.sut.<Action>()
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})
}
```

## Integration test checklist

- place files under `tests/*_integration_tests/**`
- keep file name `when_<behavior>_test.go`
- use `Given ...` / `Should ...` nested subtests
- use `defer When<Behavior>AfterEach(suite)` in each test scenario
- use real PostgreSQL/Redis containers through your shared `container_builders` test helpers
- run targeted package tests: `go test ./tests/<integration_package>...`
- keep iterating (fix + rerun) until relevant tests pass
