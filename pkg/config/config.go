package config

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

// constants
const configName = "visualization-api"
const configDirPath = "/etc/platformvisibility/visualization-api"

/*
  CONFIG NAMES
  The following logic has 3 options to specify config value
    1 - conf file
    2 - env variable
    3 - command line argument
  Conf file values are overriden by env variables,
  Env variables are overriden by command line arguments

  Example of usage:
    given name "mysql.port":
      * config parser would split this string by dot and would try
        to find "port" option in "mysql" section
      * env variables parser would replace "." to "_" and would
        capitalize the strigs, then it would try to find the variable
        named "MYSQL_PORT"
      * command line parser would replace "." to "-" and look for
        "--mysql-port" option
*/
const logLevelConfigName = "log.level"
const logFileConfigName = "log.path"

const mysqlPortConfigName = "mysql.port"
const mysqlPasswordConfigName = "mysql.password"
const mysqlHostConfigName = "mysql.host"
const mysqlUserConfigName = "mysql.username"
const mysqlDatabaseConfigName = "mysql.database"

// VisualizationAPIConfig is a struct that keeps all application config options
type VisualizationAPIConfig struct {
	// logging settings
	LogFilePath  string
	LogLevel     string
	ConsoleDebug bool

	// mysql settings
	MysqlHost         string
	MysqlPassword     string
	MysqlUsername     string
	MysqlDatabaseName string
	MysqlPort         int
}

var (
	singleToneConfig *VisualizationAPIConfig
)

var flagReplacer = strings.NewReplacer(".", "-", "_", "-")
var envReplacer = strings.NewReplacer(".", "_")

func initializeCommandLineFlags() error {
	// define flags here
	flag.String(flagReplacer.Replace(logFileConfigName), "",
		"Path to log file")
	flag.String(flagReplacer.Replace(mysqlPortConfigName), "",
		"Port mysql server is listening on")
	flag.String(flagReplacer.Replace(mysqlPasswordConfigName), "",
		"Password to authenticate on mysql server")
	flag.String(flagReplacer.Replace(mysqlHostConfigName), "",
		"Host mysql server is running on")
	flag.String(flagReplacer.Replace(mysqlUserConfigName), "",
		"Username to authenticate on mysql server")
	flag.String(flagReplacer.Replace(mysqlDatabaseConfigName), "",
		"Database to use on mysql server")

	flag.Bool("debug", false, "display debug messages in stdout")

	flag.Parse()

	flagsToBind := []string{
		logFileConfigName,
		mysqlHostConfigName,
		mysqlPortConfigName,
		mysqlUserConfigName,
		mysqlDatabaseConfigName,
		mysqlPasswordConfigName,
	}
	for _, configName := range flagsToBind {
		err := viper.BindPFlag(configName, flag.Lookup(
			flagReplacer.Replace(configName)))
		if err != nil {
			return err
		}
	}

	err := viper.BindPFlag("logging.consoleDebug", flag.Lookup(
		"debug"))

	return err
}

// InitializeConfig parses application configuration from config file, env
// variables and console flags. parsed configs are stored in module level variable
func InitializeConfig() error {
	var err error

	// assign address of default initialized VisualizationApiConfig
	// to global pointer
	singleToneConfig = &VisualizationAPIConfig{}

	// initialize path to conf files
	viper.SetConfigName(configName)
	viper.AddConfigPath(configDirPath)

	// set automatic parse of environment variables
	// env variables have higher priority then config file values,
	// but are overriden with command line flags
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(envReplacer)

	// initialize viper to read parameters passed by commandline via
	// argv[], commandline variables have higher priority
	// then env variables and config file values
	err = initializeCommandLineFlags()
	if err != nil {
		return err
	}

	// parse configs using viper lib
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	// check if config value was captured
	logFileConfigValue := viper.GetString(
		logFileConfigName)
	if logFileConfigValue == "" {
		return NewParseError(
			"logPath", "path", "log", "LOG_PATH", "--log-path")
	}
	singleToneConfig.LogFilePath = logFileConfigValue

	logLevelConfigValue := viper.GetString(
		logLevelConfigName)
	if logLevelConfigValue == "" {
		return NewParseError(
			"logLevel", "level", "log", "LOG_LEVEL", "")
	}
	singleToneConfig.LogLevel = logLevelConfigValue

	mysqlHostConfigValue := viper.GetString(
		mysqlHostConfigName)
	if mysqlHostConfigValue == "" {
		return NewParseError(
			"mysqlHost", "host", "mysql", "MYSQL_HOST", "--mysql-host")
	}
	singleToneConfig.MysqlHost = mysqlHostConfigValue

	mysqlUserConfigValue := viper.GetString(
		mysqlUserConfigName)
	if mysqlUserConfigValue == "" {
		return NewParseError(
			"mysqlUser", "username", "mysql", "MYSQL_USERNAME", "--mysql-username")
	}
	singleToneConfig.MysqlUsername = mysqlUserConfigValue

	mysqlPasswordConfigValue := viper.GetString(
		mysqlPasswordConfigName)
	if mysqlPasswordConfigValue == "" {
		return NewParseError(
			"MysqlPassword", "password", "mysql", "MYSQL_PASSWORD", "--mysql-password")
	}
	singleToneConfig.MysqlPassword = mysqlPasswordConfigValue

	mysqlDatabaseConfigValue := viper.GetString(
		mysqlDatabaseConfigName)
	if mysqlDatabaseConfigValue == "" {
		return NewParseError(
			"MysqlDatabase", "database", "mysql", "MYSQL_DATABASE", "--mysql-database")
	}
	singleToneConfig.MysqlDatabaseName = mysqlDatabaseConfigValue

	mysqlPortConfigValue := viper.GetInt(
		mysqlPortConfigName)
	if mysqlPortConfigValue == 0 {
		return NewParseError(
			"MysqlPort", "port", "mysql", "MYSQL_PORT", "--mysql-port")
	}
	singleToneConfig.MysqlPort = mysqlPortConfigValue

	// console debug has default values - no need to check
	singleToneConfig.ConsoleDebug = viper.GetBool(
		"logging.consoleDebug")
	return nil
}

// GetConfig returns structure with all application configuration
func GetConfig() *VisualizationAPIConfig {
	return singleToneConfig
}
