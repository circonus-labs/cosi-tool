// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/circonus-labs/cosi-tool/internal/release"
	"github.com/circonus-labs/go-apiclient"
	"github.com/fatih/color"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var client *apiclient.API
var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cosi",
	Short: "CLI for managing a COSI registration",
	Long: `A command line tool for registering a system with Circonus
and managing the local registration.
`,
	PersistentPreRunE: initApp,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zlog := zerolog.New(zerolog.SyncWriter(os.Stderr)).With().Timestamp().Logger()
	log.Logger = zlog

	stdlog.SetFlags(0)
	stdlog.SetOutput(zlog)

	cobra.OnInitialize(initConfig)

	desc := func(desc, env string) string {
		return fmt.Sprintf("[ENV: %s] %s", env, desc)
	}

	// configuration file
	{
		var (
			longOpt     = "config"
			shortOpt    = "c"
			description = "config file (default: " + defaults.ConfigFile + "|.json|.toml)"
		)
		RootCmd.PersistentFlags().StringVarP(&cfgFile, longOpt, shortOpt, "", description)
	}

	// registration options configuration file
	{
		const (
			key         = config.KeyRegConf
			longOpt     = "regconf"
			envVar      = release.ENVPREFIX + "_REG_CONF"
			description = "Registration options configuration file"
		)
		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}

	//
	// Circonus API
	//
	{
		const (
			key         = config.KeyAPITokenKey
			longOpt     = "api-key"
			envVar      = release.ENVPREFIX + "_API_KEY"
			description = "Circonus API Token Key"
		)

		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}
	{
		const (
			key         = config.KeyAPITokenApp
			longOpt     = "api-app"
			envVar      = release.ENVPREFIX + "_API_APP"
			description = "Circonus API Token App Name"
		)

		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}
	{
		const (
			key         = config.KeyAPIURL
			longOpt     = "api-url"
			envVar      = release.ENVPREFIX + "_API_URL"
			description = "Circonus API URL"
		)

		RootCmd.PersistentFlags().String(longOpt, defaults.APIURL, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.APIURL)
	}
	{
		const (
			key         = config.KeyAPICAFile
			longOpt     = "api-ca-file"
			envVar      = release.ENVPREFIX + "_API_CA_FILE"
			description = "Circonus API Certificate CA file"
		)

		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}

	// cosi url
	{
		const (
			key         = config.KeyCosiURL
			longOpt     = "cosi-url"
			envVar      = release.ENVPREFIX + "_URL"
			description = "Circonus One Step Install (cosi server) URL"
		)

		RootCmd.PersistentFlags().String(longOpt, defaults.COSIURL, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.COSIURL)
	}

	//
	// Agent settings
	//
	{
		const (
			key         = config.KeyAgentMode
			longOpt     = "agent-mode"
			envVar      = release.ENVPREFIX + "_AGENT_MODE"
			description = "Agent mode for check (reverse|pull)"
		)

		RootCmd.PersistentFlags().String(longOpt, defaults.AgentMode, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.AgentMode)
	}
	{
		const (
			key         = config.KeyAgentURL
			longOpt     = "agent-url"
			envVar      = release.ENVPREFIX + "_AGENT_URL"
			description = "URL the Circonus Agent is listening on"
		)

		RootCmd.PersistentFlags().String(longOpt, defaults.AgentURL, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.AgentURL)
	}

	//
	// Host
	//
	{
		const (
			key         = config.KeyHostTarget
			longOpt     = "check-target"
			envVar      = release.ENVPREFIX + "_CHECK_TARGET"
			description = "Check target(host) to use when creating system check"
		)

		RootCmd.PersistentFlags().String(longOpt, defaults.HostCheckTarget, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.HostCheckTarget)
	}
	{
		const (
			key         = config.KeyHostGroupID
			longOpt     = "group-id"
			envVar      = release.ENVPREFIX + "_GROUP_ID"
			description = "Group ID for multi-system check"
		)

		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}
	{
		const (
			key         = config.KeyHostBrokerID
			longOpt     = "broker-id"
			envVar      = release.ENVPREFIX + "_BROKER_ID"
			description = "Broker ID to use when creating check [0=auto select] (default 0)"
		)

		RootCmd.PersistentFlags().Uint(longOpt, defaults.HostBrokerID, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.HostBrokerID)
	}
	{
		const (
			key         = config.KeyHostBrokerType
			longOpt     = "broker-type"
			envVar      = release.ENVPREFIX + "_BROKER_TYPE"
			description = "Limit automatic broker selection to a specific type of broker"
		)

		RootCmd.PersistentFlags().String(longOpt, defaults.HostBrokerType, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.HostBrokerType)
	}

	//
	// System - these are automatically generated by cosi-install
	//          if manual, they need to be added correctly for the system/os
	//
	{
		const (
			key         = config.KeySystemArch
			longOpt     = "sys-arch"
			envVar      = release.ENVPREFIX + "_SYS_ARCH"
			description = "System architecture (generated by cosi-install)"
		)

		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}
	{
		const (
			key         = config.KeySystemOSType
			longOpt     = "os-type"
			envVar      = release.ENVPREFIX + "_OS_TYPE"
			description = "OS type (generated by cosi-install)"
		)

		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}
	{
		const (
			key         = config.KeySystemOSDistro
			longOpt     = "os-distro"
			envVar      = release.ENVPREFIX + "_OS_DISTRO"
			description = "OS distribution (generated by cosi-install)"
		)

		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}
	{
		const (
			key         = config.KeySystemOSVersion
			longOpt     = "os-version"
			envVar      = release.ENVPREFIX + "_OS_VERSION"
			description = "OS distribution version (generated by cosi-install)"
		)

		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}
	{
		const (
			key         = config.KeySystemDMI
			longOpt     = "sys-dmi"
			envVar      = release.ENVPREFIX + "_SYS_DMI"
			description = "System dmi bios version (generated by cosi-install, only used in AWS)"
		)

		RootCmd.PersistentFlags().String(longOpt, "", desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
	}

	//
	// Miscellaneous settings
	//
	{
		const (
			key         = config.KeyDebug
			longOpt     = "debug"
			shortOpt    = "d"
			envVar      = release.ENVPREFIX + "_DEBUG"
			description = "Enable debug messages"
		)

		RootCmd.PersistentFlags().BoolP(longOpt, shortOpt, defaults.Debug, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.Debug)
	}

	{
		const (
			key         = config.KeyLogLevel
			longOpt     = "log-level"
			envVar      = release.ENVPREFIX + "_LOG_LEVEL"
			description = "Log level [(panic|fatal|error|warn|info|debug|disabled)]"
		)

		RootCmd.PersistentFlags().String(longOpt, defaults.LogLevel, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.LogLevel)
	}

	{
		const (
			key         = config.KeyLogPretty
			longOpt     = "log-pretty"
			envVar      = release.ENVPREFIX + "_LOG_PRETTY"
			description = "Output formatted/colored log lines [ignored on windows]"
		)

		RootCmd.PersistentFlags().Bool(longOpt, defaults.LogPretty, desc(description, envVar))
		viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longOpt))
		viper.BindEnv(key, envVar)
		viper.SetDefault(key, defaults.LogPretty)
	}
}

// initLogging initializes zerolog
func initLogging() error {
	//
	// Enable formatted output
	//
	if viper.GetBool(config.KeyLogPretty) {
		if runtime.GOOS != "windows" {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		} else {
			log.Warn().Msg("log-pretty not applicable on this platform")
		}
	}

	//
	// Enable debug logging, if requested
	// otherwise, default to info level and set custom level, if specified
	//
	if viper.GetBool(config.KeyDebug) {
		viper.Set(config.KeyLogLevel, "debug")
		log.Info().Msg("setting log level to debug")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		if viper.IsSet(config.KeyLogLevel) {
			level := viper.GetString(config.KeyLogLevel)

			switch level {
			case "panic":
				zerolog.SetGlobalLevel(zerolog.PanicLevel)
			case "fatal":
				zerolog.SetGlobalLevel(zerolog.FatalLevel)
			case "error":
				zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			case "warn":
				zerolog.SetGlobalLevel(zerolog.WarnLevel)
			case "info":
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			case "debug":
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			case "disabled":
				zerolog.SetGlobalLevel(zerolog.Disabled)
			default:
				return errors.Errorf("Unknown log level (%s)", level)
			}

			log.Debug().Str("log-level", level).Msg("Logging level")
		}
	}

	return nil
}

// initAPIClient initializes global api client used by majority of commands
func initAPIClient() error {
	opt := &apiclient.Config{
		URL:      viper.GetString(config.KeyAPIURL),
		TokenKey: viper.GetString(config.KeyAPITokenKey),
		TokenApp: viper.GetString(config.KeyAPITokenApp),
		Debug:    viper.GetBool(config.KeyDebug),
		Log:      stdlog.New(log.With().Str("pkg", "cgm.api").Logger(), "", 0),
	}

	if viper.GetString(config.KeyAPICAFile) != "" {
		data, err := ioutil.ReadFile(viper.GetString(config.KeyAPICAFile))
		if err != nil {
			return errors.Wrap(err, "reading api ca file")
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(data) {
			return errors.Errorf("unable to add cert from api ca file")
		}
		opt.TLSConfig = &tls.Config{
			RootCAs: cp,
		}
	}

	c, err := apiclient.New(opt)
	if err != nil {
		return errors.Wrap(err, "unable to initialize api client")
	}
	client = c
	return nil
}

// verifyConfig checks required configuration settings
func verifyConfig() error {
	if viper.GetString(config.KeyAPITokenKey) == "" {
		return errors.Errorf("invalid API Token Key - missing")
	}
	if viper.GetString(config.KeyAPITokenApp) == "" {
		return errors.Errorf("invalid API Token App name - missing")
	}

	// set defaults if settings with defaults are intentionally set to blank
	if viper.GetString(config.KeyAPIURL) == "" {
		viper.Set(config.KeyAPIURL, defaults.APIURL)
	}
	if viper.GetString(config.KeyAgentMode) == "" {
		viper.Set(config.KeyAgentMode, defaults.AgentMode)
	}
	if viper.GetString(config.KeyAgentURL) == "" {
		viper.Set(config.KeyAgentURL, defaults.AgentURL)
	}
	if viper.GetString(config.KeyUIBaseURL) != "" && !strings.HasSuffix(viper.GetString(config.KeyUIBaseURL), "/") {
		viper.Set(config.KeyUIBaseURL, viper.GetString(config.KeyUIBaseURL)+"/")
	}

	cosiID, err := getCosiID(path.Join(defaults.EtcPath, ".cosi_id"))
	if err != nil {
		return err
	}
	viper.Set(config.KeyCosiID, cosiID)

	return nil
}

func initApp(cmd *cobra.Command, args []string) error {

	if err := initLogging(); err != nil {
		return err
	}

	fullcmd := strings.Join(os.Args, " ")
	if !strings.Contains(fullcmd, "config init") && !strings.Contains(fullcmd, "config show") {
		if err := verifyConfig(); err != nil {
			return err
		}

		if err := initAPIClient(); err != nil {
			return err
		}
		if viper.GetString(config.KeyUIBaseURL) == "" {
			acct, err := client.FetchAccount(nil)
			if err != nil {
				return errors.Wrap(err, "fetching account to set ui base url")
			}
			viper.Set(config.KeyUIBaseURL, acct.UIBaseURL)
			f := viper.ConfigFileUsed()
			if f != "" {
				// if there was a config file, save the ui base url for future
				// so cosi doesn't hit the api for every command
				// NOTE: but ONLY save what is in the config file, not the merged
				//       configuration - leave command line parameters, ENV vars,
				//       and defaults as-is.
				var cfg map[string]interface{} //config.Config
				if err := config.LoadConfigFile(f, &cfg); err != nil {
					return errors.Wrap(err, "reading config to set ui base url")
				}

				// cfg.UIBaseURL = acct.UIBaseURL
				cfg[config.KeyUIBaseURL] = acct.UIBaseURL

				if err := config.SaveConfigFile(f, cfg, true); err != nil {
					return errors.Wrap(err, "saving updated configuration file")
				}
			}
		}
	}

	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(defaults.EtcPath)
		viper.AddConfigPath(".")
		viper.SetConfigName(release.NAME)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		f := viper.ConfigFileUsed()
		if f != "" { // explicit config file specified
			fullcmd := strings.Join(os.Args, " ")
			// error if the command is not 'config init'
			if !strings.Contains(fullcmd, "config init") {
				yellow := color.New(color.FgYellow).SprintFunc()
				red := color.New(color.FgRed).SprintFunc()
				fmt.Printf("%s: %s\n", red("Unable to load configuration file"), yellow(err.Error()))
				os.Exit(1)
			}
		}
	}
}

func getCosiID(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
	}

	id := strings.TrimSpace(string(data))
	if id != "" {
		return id, nil
	}

	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(file, []byte(u.String()), 0644); err != nil {
		return "", err
	}
	return u.String(), nil
}
