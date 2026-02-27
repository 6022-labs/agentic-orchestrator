package settings

import (
	"errors"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const defaultConfigFile = "appsettings.json"

func LoadKoanfConfig(configFile string) (*koanf.Koanf, error) {
	k := koanf.New(".")

	if configFile == "" {
		configFile = defaultConfigFile
	}

	if _, err := os.Stat(configFile); err == nil {
		if err := k.Load(file.Provider(configFile), json.Parser()); err != nil {
			return nil, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if err := loadEnvIntoKoanf(k); err != nil {
		return nil, err
	}

	return k, nil
}

func loadEnvIntoKoanf(k *koanf.Koanf) error {
	envMap := mapEnvEntries(os.Environ())
	if len(envMap) == 0 {
		return nil
	}

	return k.Load(confmap.Provider(envMap, "."), nil)
}

func mapEnvEntries(entries []string) map[string]interface{} {
	envMap := map[string]interface{}{}

	for _, entry := range entries {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]
		if !strings.Contains(key, "__") {
			continue
		}

		normalized := strings.ToLower(strings.ReplaceAll(key, "__", "."))
		envMap[normalized] = value
	}

	return envMap
}
