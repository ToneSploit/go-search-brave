package functions

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration values loaded from env or .env file.
type Config struct {
	BraveAPIKey   string
	Keywords      []string
	SearchOptions SearchOptions
	Pushover      PushoverConfig
}

// LoadConfig reads configuration from a .env file (if present) and OS environment variables.
// If envFile is non-empty, that file is loaded instead of the default ".env".
// OS environment variables take precedence over .env file values.
func LoadConfig(envFile string) (*Config, error) {
	if envFile != "" {
		viper.SetConfigFile(envFile)
		viper.SetConfigType("env")
	} else {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
	}

	// Read .env file if it exists; ignore error if missing.
	_ = viper.ReadInConfig()

	// Override with actual OS env vars.
	viper.AutomaticEnv()

	apiKey := viper.GetString("BRAVE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("BRAVE_API_KEY is required but not set")
	}

	count := viper.GetInt("SEARCH_COUNT")
	if count == 0 {
		count = 20
	}
	freshness := viper.GetString("SEARCH_FRESHNESS")
	if freshness == "" {
		freshness = "pd"
	}

	var keywords []string
	if raw := viper.GetString("SEARCH_KEYWORDS"); raw != "" {
		for _, k := range strings.Split(raw, ",") {
			if k = strings.TrimSpace(k); k != "" {
				keywords = append(keywords, k)
			}
		}
	}

	return &Config{
		BraveAPIKey: apiKey,
		Keywords:    keywords,
		SearchOptions: SearchOptions{
			Freshness:  freshness,
			Count:      count,
			Offset:     viper.GetInt("SEARCH_OFFSET"),
			Country:    viper.GetString("SEARCH_COUNTRY"),
			SearchLang: viper.GetString("SEARCH_LANG"),
		},
		Pushover: PushoverConfig{
			AppToken: viper.GetString("PUSHOVER_APP_TOKEN"),
			UserKey:  viper.GetString("PUSHOVER_USER_KEY"),
			Enabled:  viper.GetBool("PUSHOVER_ENABLED"),
		},
	}, nil
}
