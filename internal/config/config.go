package config

import (
	"os"
)

const (
	httpPortKey = "HTTP_PORT"

	deepgramAPIKey = "DEEPGRAM_API_KEY"
	deepgramURL    = "DEEPGRAM_URL"

	geminiAPIKey = "GEMINI_API_KEY"
	geminiURL    = "GEMINI_URL"

	databaseURL      = "DATABASE_URL"
	postgresHost     = "POSTGRES_HOST"
	postgresPort     = "POSTGRES_PORT"
	postgresUser     = "POSTGRES_USER"
	postgresPassword = "POSTGRES_PASSWORD"
	postgresDBName   = "POSTGRES_DB_NAME"

	minioEndpoint      = "MINIO_ENDPOINT"
	minioAccessKey     = "MINIO_ACCESS_KEY"
	minioSecretKey     = "MINIO_SECRET_KEY"
	minioUseSSL        = "MINIO_USE_SSL"
	minioImagesBucket  = "MINIO_IMAGES_BUCKET"
	minioAnswersBucket = "MINIO_ANSWERS_BUCKET"
)

type Config struct {
	HTTPPort string

	Postgres *DB
	Deepgram *ExternalAPI
	Gemini   *ExternalAPI
	Minio    *Minio
}

func New() Config {
	HTTPPort := os.Getenv(httpPortKey)

	Deepgram := ExternalAPI{
		APIKey: os.Getenv(deepgramAPIKey),
		URL:    os.Getenv(deepgramURL),
	}

	Gemini := ExternalAPI{
		APIKey: os.Getenv(geminiAPIKey),
		URL:    os.Getenv(geminiURL),
	}

	Postgres := DB{
		URL:      os.Getenv(databaseURL),
		Host:     os.Getenv(postgresHost),
		Port:     os.Getenv(postgresPort),
		User:     os.Getenv(postgresUser),
		Password: os.Getenv(postgresPassword),
		DBName:   os.Getenv(postgresDBName),
	}

	Minio := Minio{
		Endpoint:        os.Getenv(minioEndpoint),
		AccessKey:       os.Getenv(minioAccessKey),
		SecretAccessKey: os.Getenv(minioSecretKey),
		UseSSL:          os.Getenv(minioUseSSL) == "true",

		ImagesBucket:  os.Getenv(minioImagesBucket),
		AnswersBucket: os.Getenv(minioAnswersBucket),
	}

	return Config{
		HTTPPort: HTTPPort,

		Deepgram: &Deepgram,
		Gemini:   &Gemini,
		Postgres: &Postgres,
		Minio:    &Minio,
	}
}

type ExternalAPI struct {
	APIKey string
	URL    string
}

type Minio struct {
	Endpoint        string
	AccessKey       string
	SecretAccessKey string
	UseSSL          bool

	ImagesBucket  string
	AnswersBucket string
}

type DB struct {
	URL      string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (db *DB) GetConnStr() string {
	// If DATABASE_URL is provided (Railway/Heroku), use it directly
	if db.URL != "" {
		return db.URL
	}

	// Otherwise build from individual components with SSL
	return "host=" + db.Host +
		" port=" + db.Port +
		" user=" + db.User +
		" password=" + db.Password +
		" dbname=" + db.DBName +
		" sslmode=require"
}
