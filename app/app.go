package app

import (
	"context"
	"fmt"
	"os"

	logger "github.com/izumin5210/gentleman-logger"
	"github.com/mikloslorinczi/todoadm/model"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/h2non/gentleman.v2"

	"github.com/donbattery/colorlog"
)

const (
	defaultConfFile = "todoadm_conf"
	confKey         = "ConfKey"
)

// App is the top level object of the application
type App struct {
	ctx    context.Context
	log    *colorlog.ColorLogger
	cmd    *cobra.Command
	client *gentleman.Client
	debug  bool
}

// Build build up the Application and returns it
func Build() *App {
	app := &App{
		log: colorlog.NewColorLogger("Cmd", false),
		ctx: context.Background(),
	}

	app.cmd = &cobra.Command{
		Long:               app.getLong(),
		PersistentPreRunE:  app.setup,
		RunE:               app.run,
		PersistentPostRunE: app.cleanup,
	}

	app.cmd.PersistentFlags().StringP("config_file", "c", "", "Configuration file path")
	app.cmd.PersistentFlags().StringP("todo_host", "H", "http://127.0.0.1:6969", "TODO Service URL")
	app.cmd.PersistentFlags().StringP("admin_header", "A", "TodoAdminHeader", "TODO Admin Header")
	app.cmd.PersistentFlags().StringP("admin_secret", "S", "TodoAdminSecret", "TODO Admin Secret")
	app.cmd.PersistentFlags().Bool("debug", false, "Debug mode")

	app.cmd.SilenceErrors = true
	app.cmd.SilenceUsage = true

	app.cmd.AddCommand(app.buildPingCmd())

	return app
}

// Execute runs the undelying root Cobra command of the application
// if any error occures during execution it will be logged and the program will exit with 1
func (app *App) Execute() {
	if err := app.cmd.Execute(); err != nil {
		app.log.Fatal(err)
	}
	os.Exit(0)
}

func (app *App) getLong() string {
	return `TODO Admin CLI`
}

// setup creates and validates the App's configuration object
// considering defaults, the config file, environment variables and command line flags
func (app *App) setup(cmd *cobra.Command, args []string) error {
	// Get the default configs
	conf := model.DefaultConfigs()

	// Bind Cobra flags to Viper keys
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return errors.Wrap(err, "Failed to bind Cobra flags to Viper keys")
	}

	// Search the environment for variables prefixed with TODOADM_
	viper.SetEnvPrefix("TODOADM")

	// Read in environment variables that match
	viper.AutomaticEnv()

	if viper.GetBool("debug") {
		app.debug = true
	}

	app.log = colorlog.NewColorLogger("Cmd", app.debug)
	app.log.Debug("Logger has been set up")

	// If a ConfigFile is set read it in
	if cPath := viper.GetString("config_file"); cPath != "" {
		// Use config file from the flag
		viper.SetConfigFile(cPath)
	} else {
		viper.SetConfigName(defaultConfFile)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		app.log.Debug(errors.Wrap(err, "Failed to load configs from file"))
	}

	app.log.Debugf("Configuration file used: %s", viper.ConfigFileUsed())

	// Load the Viper configs into the conf object
	if err := viper.Unmarshal(&conf); err != nil {
		return errors.Wrap(err, "Failed to decode config object")
	}

	app.log.Debugf("Configuration object:\n\n%+v\n", conf)

	// Validate the configurations
	if err := conf.Validate(); err != nil {
		return errors.Wrap(err, "Invalid Configuration")
	}

	if app.debug {
		app.log.Debug("Viper Settings:")
		viper.Debug()
		for _, key := range viper.AllKeys() {
			app.log.Debugf("Viper key: %s Value: %v", key, viper.Get(key))
		}
	}

	// Load the configs into the context
	app.setVal(confKey, conf)

	app.initClient()

	return nil
}

func (app *App) initClient() {
	cfg := app.getConf()

	app.client = gentleman.New()
	if app.debug {
		app.client.Use(logger.New(os.Stdout))
	}

	app.client.URL(fmt.Sprintf("%s/v1", cfg.TodoHost))

	app.client.AddHeader(cfg.AdminHeader, cfg.AdminSecret)
}

func (app *App) run(cmd *cobra.Command, args []string) error {
	return app.cmd.Usage()
}

func (app *App) cleanup(cmd *cobra.Command, args []string) error {
	return nil
}

func (app *App) setVal(key string, value interface{}) {
	app.ctx = context.WithValue(app.ctx, key, value)
}

func (app *App) getVal(key string) interface{} {
	return app.ctx.Value(key)
}

func (app *App) getConf() *model.Configs {
	if val, ok := app.ctx.Value(confKey).(*model.Configs); ok {
		return val
	}
	return nil
}
