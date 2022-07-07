package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// AppDir is the name of the directory where the config file is stored.
const AppDir = ".jelliflix"

// FileName is the name of the config file that gets created.
const FileName = "config.yml"

// SettingsConfig struct represents the config for the settings.
type SettingsConfig struct {
	Timeout         string `yaml:"timeout"`
	RefreshRate     string `yaml:"refresh_rate"`
	EnableLogging   bool   `yaml:"enable_logging"`
	StorageDataDir  string `yaml:"storage_data_dir"`
	DownloadQuality string `yaml:"download_quality"`
}

// TorrentConfig represents the config for torrent.
type TorrentConfig struct {
	Port      int      `yaml:"port"`
	Providers []string `yaml:"providers"`
}

// ServiceConfig represents the config for media services.
type ServiceConfig struct {
	Enabled    bool   `yaml:"enabled"`
	DataDir    string `yaml:"data_dir"`
	MoviesPath string `yaml:"movies_path"`
	SeriesPath string `yaml:"series_path"`
}

// ProviderConfig represents the config for watchlist providers.
type ProviderConfig struct {
	ID string `yaml:"id"`
}

// Config represents the main config for the application.
type Config struct {
	Settings SettingsConfig `yaml:"settings"`
	Exporter TorrentConfig  `yaml:"exporter"`
	Services struct {
		Jellyfin ServiceConfig `yaml:"jellyfin"`
		Plex     ServiceConfig `yaml:"plex"`
		Emby     ServiceConfig `yaml:"emby"`
	} `yaml:"services"`
	Providers struct {
		IMDb ProviderConfig `yaml:"imdb"`
	} `yaml:"providers"`
}

// configError represents an error that occurred while parsing the config file.
type configError struct {
	configDir string
	parser    Parser
	err       error
}

// Parser is the parser for the config file.
type Parser struct{}

// getDefaultConfig returns the default config for the application.
func (parser Parser) getDefaultConfig() Config {
	return Config{
		Settings: SettingsConfig{
			Timeout:         "10s",
			RefreshRate:     "10s",
			EnableLogging:   false,
			StorageDataDir:  ".",
			DownloadQuality: "720p",
		},
		Exporter: TorrentConfig{
			Port:      42342,
			Providers: []string{"RARBG", "YTS", "TPB"},
		},
	}
}

// getDefaultConfigYamlContents returns the default config file contents.
func (parser Parser) getDefaultConfigYamlContents() string {
	defaultConfig := parser.getDefaultConfig()
	yml, _ := yaml.Marshal(defaultConfig)

	return string(yml)
}

// Error returns the error message for when a config file is not found.
func (e configError) Error() string {
	return fmt.Sprintf(`
Couldn't find a config.yml configuration file.
Create one under: %s
Example of a config.yml file:
%s
For more info, go to https://github.com/jeliiflix/jelliflix
press q to exit.
Original error: %v`,
		path.Join(e.configDir, AppDir, FileName),
		e.parser.getDefaultConfigYamlContents(),
		e.err,
	)
}

// writeDefaultConfigContents writes the default config file contents to the given file.
func (parser Parser) writeDefaultConfigContents(newConfigFile *os.File) error {
	_, err := newConfigFile.WriteString(parser.getDefaultConfigYamlContents())

	if err != nil {
		return err
	}

	return nil
}

// createConfigFileIfMissing creates the config file if it doesn't exist.
func (parser Parser) createConfigFileIfMissing(configFilePath string) error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		newConfigFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}

		defer func() { _ = newConfigFile.Close() }()
		return parser.writeDefaultConfigContents(newConfigFile)
	}

	return nil
}

// getConfigFileOrCreateIfMissing returns the config file path or creates the config file if it doesn't exist.
func (parser Parser) getConfigFileOrCreateIfMissing() (*string, error) {
	var err error
	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		configDir, err = os.UserHomeDir()
		if err != nil {
			return nil, configError{parser: parser, configDir: configDir, err: err}
		}
	}

	prsConfigDir := filepath.Join(configDir, AppDir)
	err = os.MkdirAll(prsConfigDir, os.ModePerm)
	if err != nil {
		return nil, configError{parser: parser, configDir: configDir, err: err}
	}

	configFilePath := filepath.Join(prsConfigDir, FileName)
	err = parser.createConfigFileIfMissing(configFilePath)
	if err != nil {
		return nil, configError{parser: parser, configDir: configDir, err: err}
	}

	return &configFilePath, nil
}

// parsingError represents an error that occurred while parsing the config file.
type parsingError struct {
	err error
}

// Error represents an error that occurred while parsing the config file.
func (e parsingError) Error() string {
	return fmt.Sprintf("failed parsing config.yml: %v", e.err)
}

// readConfigFile reads the config file and returns the config.
func (parser Parser) readConfigFile(path string) (Config, error) {
	config := parser.getDefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return config, configError{parser: parser, configDir: path, err: err}
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}

// initParser initializes the parser.
func initParser() Parser {
	return Parser{}
}

// ParseConfig parses the config file and returns the config.
func ParseConfig() (Config, error) {
	var config Config
	var err error

	parser := initParser()

	configFilePath, err := parser.getConfigFileOrCreateIfMissing()
	if err != nil {
		return config, parsingError{err: err}
	}

	config, err = parser.readConfigFile(*configFilePath)
	if err != nil {
		return config, parsingError{err: err}
	}

	return config, nil
}
