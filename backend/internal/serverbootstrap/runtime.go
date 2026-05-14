package serverbootstrap

import (
	"path"

	"github.com/thunder-id/thunderid/internal/system/config"
)

// InitializeRuntime loads deployment.yaml and default.json under serverHome and
// installs the global server runtime (same as the Thunder server binary).
func InitializeRuntime(serverHome string) (*config.Config, error) {
	configFilePath := path.Join(serverHome, "repository/conf/deployment.yaml")
	defaultConfigPath := path.Join(serverHome, "repository/resources/conf/default.json")
	cfg, err := config.LoadConfig(configFilePath, defaultConfigPath, serverHome)
	if err != nil {
		return nil, err
	}
	if err := config.InitializeServerRuntime(serverHome, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
