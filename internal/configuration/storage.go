package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func cfgPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".cgn", "api", "config.yaml"), nil
}

func lockPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".cgn", "api", ".lock"), nil
}

func grantConfigFile(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
		if err != nil {
			return err
		}

		file, err := os.Create(configPath)
		if err != nil {
			return err
		}
		defer file.Close()

		defaultConfig := Config{
			// Initialize with default values
		}
		data, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			return err
		}

		_, err = file.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func lock() error {
	lockFilePath, err := lockPath()
	if err != nil {
		return err
	}

	_, err = os.OpenFile(lockFilePath, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("lock file already exists, if stuck, remove manually %s", lockFilePath)
		}
		return err
	}
	return nil
}

func unlock() {
	lockFilePath, _ := lockPath()
	os.Remove(lockFilePath)
}

func readConfig() (Config, error) {
	configPath, err := cfgPath()
	if err != nil {
		return Config{}, err
	}

	err = grantConfigFile(configPath)
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var rv Config
	yaml.Unmarshal(data, &rv)

	return rv, nil
}

func writeConfig(cfg Config) error {
	configPath, err := cfgPath()
	if err != nil {
		return err
	}

	err = grantConfigFile(configPath)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func ListProfiles() (map[string]Profile, error) {
	err := lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}

	return cfg.Profiles, nil
}

func AddProfile(name, realm, clientID, authServer, refreshToken string) error {
	if realm == "" {
		return fmt.Errorf("realm cannot be empty")
	}
	if clientID == "" {
		return fmt.Errorf("clientID cannot be empty")
	}
	if authServer == "" {
		return fmt.Errorf("authServer cannot be empty")
	}

	err := lock()
	if err != nil {
		return err
	}
	defer unlock()

	cfg, err := readConfig()
	if err != nil {
		return err
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}

	if _, exists := cfg.Profiles[name]; exists {
		return fmt.Errorf("profile %s already exists", name)
	}

	cfg.Profiles[name] = Profile{
		AuthServer:   authServer,
		RefreshToken: refreshToken,
		Realm:        realm,
		ClientID:     clientID,
	}

	writeConfig(cfg)

	return nil
}

func RemoveProfile(name string) error {
	err := lock()
	if err != nil {
		return err
	}
	defer unlock()

	cfg, err := readConfig()
	if err != nil {
		return err
	}

	if _, exists := cfg.Profiles[name]; !exists {
		return fmt.Errorf("profile %s does not exist", name)
	}

	delete(cfg.Profiles, name)

	writeConfig(cfg)

	return nil
}

func InitToken(name, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("refreshToken cannot be empty")
	}

	err := lock()
	if err != nil {
		return err
	}
	defer unlock()

	cfg, err := readConfig()
	if err != nil {
		return err
	}

	if _, exists := cfg.Profiles[name]; !exists {
		return fmt.Errorf("profile %s does not exist", name)
	}

	profile := cfg.Profiles[name]
	profile.RefreshToken = refreshToken
	cfg.Profiles[name] = profile

	writeConfig(cfg)

	return nil
}

func GetProfile(name string) (Profile, error) {
	err := lock()
	if err != nil {
		return Profile{}, err
	}
	defer unlock()

	cfg, err := readConfig()
	if err != nil {
		return Profile{}, err
	}

	if profile, exists := cfg.Profiles[name]; exists {
		return profile, nil
	}

	return Profile{}, fmt.Errorf("profile %s does not exist", name)
}

func UpdateProfile(name string, profile Profile) error {
	err := lock()
	if err != nil {
		return err
	}
	defer unlock()

	cfg, err := readConfig()
	if err != nil {
		return err
	}

	if _, exists := cfg.Profiles[name]; !exists {
		return fmt.Errorf("profile %s does not exist", name)
	}

	cfg.Profiles[name] = profile

	writeConfig(cfg)

	return nil
}
