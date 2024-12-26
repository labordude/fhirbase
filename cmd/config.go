package cmd

import (
	"fmt"

	"github.com/labordude/fhirbase/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initializeConfig(cmd *cobra.Command) error {
	opts := shared.FhirbaseConfig{}
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.SetDefault("display.cursor", ">")
	viper.SetDefault("keys.up", "up")
	viper.SetDefault("keys.down", "down")
	viper.SetDefault("keys.quit", "q")
	viper.SetDefault("keys.help", "?")
	viper.SetDefault("keys.select", "enter")
	viper.SetDefault("keys.back", "backspace")
	viper.SetDefault("keys.toggle", "space")
	viper.SetDefault("debug.log_messages", false)
	viper.SetDefault("debug.filepath", "debug.log")
	// Defaults
	viper.SetDefault("fhir", "4.0.0")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.db", "fhirbase")

	// Bind flags to viper
	viper.BindPFlag("fhir", cmd.PersistentFlags().Lookup("fhir"))
	viper.BindPFlag("database.db", cmd.PersistentFlags().Lookup("db"))
	viper.BindPFlag("database.port", cmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("database.host", cmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("database.username", cmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("database.password", cmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("database.sslmode", cmd.PersistentFlags().Lookup("sslmode"))
	/* Look for a config file in the current directory */
	configFile, _ := cmd.Flags().GetString("config")
	if configFile == "" {
		configFile = "."
	}
	viper.AddConfigPath(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// try to create the file
			if err := viper.WriteConfigAs("config.yaml"); err != nil {
				return fmt.Errorf("error creating config file: %s", err)
			}
		} else {
			return fmt.Errorf("error reading config: %s", err)
		}
	}

	if err := viper.Unmarshal(&opts); err != nil {
		return fmt.Errorf("error reading config: %s", err)
	}

	shared.FhirbaseOptions = opts
	return nil
}
