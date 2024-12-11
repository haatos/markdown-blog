package internal

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var Settings *AppSettings

func NewSettings() *AppSettings {
	settings := AppSettings{
		SuperuserEmail:                    os.Getenv("SUPERUSER_EMAIL"),
		ContactEmail:                      getEnvOrDefault("CONTACT_EMAIL", ""),
		SessionExpires:                    30 * 24 * time.Hour,
		PasswordlessAuthenticationExpires: 10 * time.Minute,
		Domain:                            getEnvOrDefault("APP_DOMAIN", "localhost"),
		Port:                              getEnvOrDefault("APP_PORT", ":8080"),
		SQLiteDatabase:                    getEnvOrDefault("SQLITE_DATABASE", "file:.///db.sqlite"),
		AWSAccessKeyID:                    os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:                os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:                         getEnvOrDefault("AWS_REGION", "eu-west-1"),
		GoogleClientID:                    os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:                os.Getenv("GOOGLE_CLIENT_SECRET"),
		Bucket:                            os.Getenv("APP_BUCKET"),
		CloudfrontDomain:                  getEnvOrDefault("CLOUDFRONT_DOMAIN", ""),
	}
	if !strings.HasPrefix(settings.Port, ":") {
		settings.Port = ":" + settings.Port
	}
	return &settings
}

func getEnvOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

type AppSettings struct {
	SuperuserEmail                    string
	ContactEmail                      string
	SessionExpires                    time.Duration
	PasswordlessAuthenticationExpires time.Duration
	SQLiteDatabase                    string
	Domain                            string
	Port                              string
	AWSAccessKeyID                    string
	AWSSecretAccessKey                string
	AWSRegion                         string
	GoogleClientID                    string
	GoogleClientSecret                string
	Bucket                            string
	CloudfrontDomain                  string
}

func (as *AppSettings) BaseURL() string {
	if as.Domain == "localhost" {
		return fmt.Sprintf("http://%s%s", as.Domain, as.Port)
	} else {
		return fmt.Sprintf("https://%s", as.Domain)
	}
}

func (as *AppSettings) SQLiteDbString(readonly bool) string {
	params := make(url.Values)
	params.Add("_journal_mode", "WAL")
	params.Add("_busy_timeout", "5000")
	params.Add("_synchronous", "NORMAL")
	params.Add("_cache_size", "-20000")
	params.Add("_foreign_keys", "true")
	if readonly {
		params.Add("mode", "ro")
	} else {
		params.Add("_txlock", "IMMEDIATE")
		params.Add("mode", "rwc")
	}

	return as.SQLiteDatabase + "?" + params.Encode()
}

func ReadDotenv() {
	path := "./.env"
	re := regexp.MustCompile(`^[^0-9][A-Z0-9_]+=.+$`)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("err opening dotenv: ", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > 0 && line[0] != '#' && re.Match(line) {
			split := strings.Split(string(line), "=")
			name := strings.TrimSpace(split[0])
			value := strings.TrimSpace(split[1])
			value = strings.Trim(value, `"`)
			os.Setenv(name, value)
		} else {
			log.Println("not including invalid or empty line", "line", string(line))
		}
	}
}
