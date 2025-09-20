package configuration

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}
// MySQLConfig holds MySQL database configuration
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

var AppConfig *Config

func setConfig() {
	AppConfig.MySQL.Host = viper.GetString("mysql_host")
	AppConfig.MySQL.Port = viper.GetString("mysql_port")
	AppConfig.MySQL.User = viper.GetString("mysql_user")
	AppConfig.MySQL.Password = viper.GetString("mysql_password")
	AppConfig.MySQL.Database = viper.GetString("mysql_database")
}

func init() {
	err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

}

func GetEnv(key string) string {
	return viper.GetString(key)
}

// LoadConfig loads configuration from various sources
func LoadConfig() error {
	if AppConfig == nil {
		AppConfig = &Config{}
	}
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Bind environment variables to config keys
	bindEnvironmentVariables()

	// Read configuration file (if it exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	setConfig()
	return nil
}

// bindEnvironmentVariables binds environment variables to configuration keys
func bindEnvironmentVariables() {
	viper.BindEnv("mysql.host", "MYSQL_HOST")
	viper.BindEnv("mysql.port", "MYSQL_PORT")
	viper.BindEnv("mysql.user", "MYSQL_USER")
	viper.BindEnv("mysql.password", "MYSQL_PASSWORD", "MYSQL_PASS") // Support both variants
	viper.BindEnv("mysql.database", "MYSQL_DB", "MYSQL_DATABASE")   // Support both variants
}

// GetMySQLConfig returns the MySQL configuration
func GetMySQLConfig() MySQLConfig {
	return AppConfig.MySQL
}

// GetMySQLConnectionString returns a formatted MySQL connection string
func GetMySQLConnectionString() string {
	mysql := AppConfig.MySQL
	//mysql://[user]:[password]@[host][:port]/[database][?options]
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Database)
}

// ReloadConfig reloads the configuration
func ReloadConfig() error {
	err := LoadConfig()
	if err != nil {
		return err
	}
	return nil
}

// Deprecated: Use GetMySQLConfig().Host instead
func GetMySQLHost() string {
	return AppConfig.MySQL.Host
}

// Deprecated: Use GetMySQLConfig().Port instead
func GetMySQLPort() string {
	return AppConfig.MySQL.Port
}

// Deprecated: Use GetMySQLConfig().User instead
func GetMySQLUser() string {
	return AppConfig.MySQL.User
}

// Deprecated: Use GetMySQLConfig().Password instead
func GetMySQLPassword() string {
	return AppConfig.MySQL.Password
}

// Deprecated: Use GetMySQLConfig().Database instead
func GetMySQLDatabase() string {
	return AppConfig.MySQL.Database
}
