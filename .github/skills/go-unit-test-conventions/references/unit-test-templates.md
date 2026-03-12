# Unit Test Templates

Use these templates to keep unit test structure consistent across this repository.

## 1) Use case/service template (with gomock suite)

```go
package <package>_test

import (
	"testing"

	"<module>/src/<domain>/<package>"
	"<module>/tests/<mock_package>"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type When<Behavior>TestingSuite struct {
	sut *<package>.<TypeUnderTest>

	<dependencyName>Mock *<mock_package>.Mock<DependencyInterface>
}

func When<Behavior>BeforeEach(t *testing.T) *When<Behavior>TestingSuite {
	mockController := gomock.NewController(t)

	<dependencyName>Mock := <mock_package>.NewMock<DependencyInterface>(mockController)

	sut := <package>.New<TypeUnderTest>(
		<dependencyName>Mock,
	)

	return &When<Behavior>TestingSuite{
		sut: sut,

		<dependencyName>Mock: <dependencyName>Mock,
	}
}

func TestWhen<Behavior>(t *testing.T) {
	t.Parallel()

	t.Run("Given <context>", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *When<Behavior>TestingSuite) {
			suite.<dependencyName>Mock.EXPECT().
				<Method>(gomock.Any()).
				Return(<value>, nil).
				Times(1)
		}

		t.Run("Should <expected result>", func(t *testing.T) {
			t.Parallel()

			suite := When<Behavior>BeforeEach(t)
			initSuite(suite)

			result, err := suite.sut.<Action>(<input>)
			assert.NoError(t, err)
			assert.Equal(t, <expected>, result)
		})
	})
}
```

## 2) Request validation template (no suite, no mocks)

```go
package requests_test

import (
	"testing"

	"<module>/src/<domain>/use_cases/requests"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidating<ReqName>Request(t *testing.T) {
	t.Parallel()

	generate<ReqName>Request := func() *requests.<ReqName>Request {
		return &requests.<ReqName>Request{}
	}

	t.Run("Given <invalid condition>", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()

			req := generate<ReqName>Request()
			// mutate req

			err := req.Validate()
			assert.Error(t, err)
		})
	})

	t.Run("Given request is valid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			req := generate<ReqName>Request()
			err := req.Validate()
			assert.NoError(t, err)
		})
	})
}
```

## 3) Settings loader template (koanf/env)

```go
package settings_test

import (
	"testing"

	"<module>/src/<domain>/settings"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoading<SettingsName>Settings(t *testing.T) {
	t.Run("Given settings provided through koanf", func(t *testing.T) {
		k := koanf.New(".")
		k.Set("<key>", "<value>")

		settingsValue := settings.New<SettingsName>Settings(zap.NewNop(), k)

		assert.NotNil(t, settingsValue)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("<ENV_KEY>", "<value>")

		k := koanf.New(".")
		settingsValue := settings.New<SettingsName>Settings(zap.NewNop(), k)

		assert.NotNil(t, settingsValue)
	})
}
```

## 4) Mock generation checklist

Run from repository root when interfaces changed:

```bash
deepsearch-mockgen -S ./src -O ./tests -A -P
```

CLI repository: https://github.com/FournyP/deepsearch-mockgen-cli

After generation:
- keep generated mock files under existing `tests/*_mocks/**` packages
- do not hand-edit generated files
- rerun targeted `go test` for affected unit test packages
- keep iterating (fix + rerun) until relevant tests pass
