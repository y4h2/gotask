package config

type Config struct {
	DB DB
}

type HTTP struct {
	Host string
	Port int
}

type DB struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func LoadConfig() Config {

	return Config{
		DB: DB{},
	}
}
