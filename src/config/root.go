package config

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/pbkdf2"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RELEASE        string        `default:"0.0.1"`
	LOG_LEVEL      zerolog.Level `default:"1"`
	IS_DEVELOPMENT bool          `default:"false"`

	DB_DSN             string `required:"True"`
	RUN_IN_TRANSACTION bool   `default:"false"`
	SECRET_KEY         string `required:"True"`
	API_KEY            string

	// Telegram
	TELEGRAM_BOT_API_KEY string `required:"True"`
	TELEGRAM_BOT_CHAT_ID int64  `required:"True"`

	// mailing
	MAILERSEND_API_KEY string `required:"True"`

	// GOOGLE RECAPTCHA
	GOOGLE_RECAPTCHA_SECRET_KEY string `required:"True"`

	// Static assets, it should not be necessary to change these defaults.
	STATIC_DIRECTORY           string `default:"static"`
	TEMPLATES_DIRECTORY        string `default:"templates"`
	HOT_RELOAD_TEMPLATE_FOLDER bool

	FORMS_CONFIG_BASE64    string
	FORMS_CONFIG_FILE_PATH string

	// FormsConfig models.FormsConfig

	apiKeySalt        []byte
	apiKeyHash        []byte
	apiKey_keyLength  int
	apiKey_iterations int
}

// TODO:  <25-09-25, Sebastiaan Van Hoecke> // These functions should be moved to somewhere else, not be here
func (c *Config) GenerateOrSetApiKey() {
	c.apiKey_iterations = 100_000
	c.apiKey_keyLength = 32

	if c.API_KEY == "" {
		c.API_KEY = uuid.NewString()
		fmt.Printf("Your API-KEY: %s\n", c.API_KEY)
	}

	c.apiKeySalt = make([]byte, 16)
	_, err := rand.Read(c.apiKeySalt)

	if err != nil {
		panic(err)
	}

	c.apiKeyHash = pbkdf2.Key([]byte(c.API_KEY), c.apiKeySalt, c.apiKey_iterations, c.apiKey_keyLength, sha256.New)

}

// TODO:  <25-09-25, Sebastiaan Van Hoecke> // These functions should be moved to somewhere else, not be here
func (c *Config) VerifyApiKey(apiKey string) bool {
	derived := pbkdf2.Key([]byte(apiKey), c.apiKeySalt, c.apiKey_iterations, c.apiKey_keyLength, sha256.New)
	if len(derived) != len(c.apiKeyHash) {
		return false
	}
	var result byte
	for i := range derived {
		result |= derived[i] ^ c.apiKeyHash[i]
	}
	return result == 0

}

func New() *Config {

	config := Config{}

	// parse env vars to struct
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("Failed to decode env vars to struct config/config.go: %s", err)
	}

	config.GenerateOrSetApiKey()

	return &config
}
