package locale

import (
	"bytes"
	"errors"
	"html/template"
	"strings"

	"github.com/sevaho/goforms/src/pkg/logger"

	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed locale.yaml
var yamlFile []byte

type Localization struct {
	FR    string `yaml:"fr"`
	NL_NL string `yaml:"nl-NL"`
	NL_BE string `yaml:"nl-BE"`
	EN    string `yaml:"en"`
	DE    string `yaml:"de"`
}

func (l Localization) Get(lang string, country string) (string, error) {
	// default to EN lang
	if lang == "" {
		lang = "EN"
	}
	switch strings.ToUpper(lang) {
	case "FR":
		if l.FR != "" {
			return l.FR, nil
		} else {
			return "", errors.New("French (FR) localization not found for lang: " + lang)
		}
	case "NL":
		// default to BE country
		if country == "" {
			country = "BE"
		}
		switch strings.ToUpper(country) {
		case "BE":
			if l.NL_BE != "" {
				return l.NL_BE, nil
			} else {
				return "", errors.New("Dutch (NL-BE) localization not found for lang: " + lang)
			}
		case "NL":
			if l.NL_NL != "" {
				return l.NL_NL, nil
			} else {
				return "", errors.New("Dutch (NL-NL) localization not found for lang: " + lang)
			}
		default:
			return "", errors.New("Country: " + country + " not supported for language: " + lang + " in localization yaml")

		}
	case "EN":
		if l.EN != "" {
			return l.EN, nil
		} else {
			return "", errors.New("English (EN) localization not found for lang: " + lang)
		}
	case "DE":
		if l.DE != "" {
			return l.DE, nil
		} else {
			return "", errors.New("German (DE) localization not found for lang: " + lang)
		}
	default:
		return "", errors.New("Lang: " + lang + " not found in localization yaml")
	}
}

var LocaleMap map[string]Localization

// Translate
func T(key string, lang string, country string) string {

	locale, err := LocaleMap[key].Get(lang, country)
	if err != nil {
		logger.Logger.Err(err).Msgf("Something went while localizing key: %s for language: %s", key, lang)
		return key
	}
	return locale
}

// Translate and directly use the variables in a template to return
func T_WithTemplate(key string, lang string, country string, vars map[string]interface{}) string {
	templ := template.Must(template.New("").Parse(T(key, lang, country)))
	buf := &bytes.Buffer{}
	templ.Execute(buf, vars)
	return buf.String()
}

func init() {
	if err := yaml.Unmarshal(yamlFile, &LocaleMap); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Unable to load locale.yaml")
	}
}
