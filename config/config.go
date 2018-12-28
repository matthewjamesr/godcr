package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/jessevdk/go-flags"
)

const (
	defaultConfigFilename    = "dcrcli.conf"
	defaultHTTPServerAddress = "127.0.0.1:1234"
)

var (
	defaultAppDataDir          = dcrutil.AppDataDir("dcrcli", false)
	defaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(defaultDcrwalletAppDataDir, "rpc.cert")
	defaultConfigFile          = filepath.Join(defaultAppDataDir, defaultConfigFilename)
)

// Config holds the top-level options for the CLI program.
type Config struct {
	ShowVersion       bool   `short:"v" long:"version" description:"Display version information and exit"`
	ConfigFile        string `short:"C" long:"configfile" description:"Path to configuration file"`
	TestNet           bool   `short:"t" long:"testnet" description:"Connects to testnet wallet instead of mainnet"`
	RPCUser           string `short:"u" long:"rpcuser" description:"RPC username"`
	RPCPassword       string `short:"p" long:"rpcpass" default-mask:"-" description:"RPC password"`
	WalletRPCServer   string `short:"w" long:"walletrpcserver" description:"Wallet RPC server to connect to"`
	RPCCert           string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	HTTPServerAddress string `short:"s" long:"serveraddress" description:"Address and port of the HTTP server."`
	HTTPMode          bool   `long:"http" description:"Run in HTTP mode."`
	NoDaemonTLS       bool   `long:"nodaemontls" description:"Disable TLS"`
}

// defaultConfig an instance of Config with the defaults set.
func defaultConfig() Config {
	return Config{
		ConfigFile:        defaultConfigFile,
		RPCCert:           defaultRPCCertFile,
		HTTPServerAddress: defaultHTTPServerAddress,
	}
}

// AppName returns the name of the program binary file that started the process.
func AppName() string {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	return appName
}

// LoadConfig parses program configuration from both the CLI flags and the config file.
func LoadConfig() (Config, *flags.Parser, error) {
	// load defaults first
	config := defaultConfig()

	parser := flags.NewParser(&config, flags.HelpFlag)
	// stub out the command handler so that the commands are not run at while loading configuration.
	//parser.CommandHandler = func(command flags.Commander, args []string) error {
	//	return nil
	//}

	_, err := parser.Parse()
	if err != nil && !IsFlagErrorType(err, flags.ErrHelp) {
		return config, parser, err
	}

	if config.ShowVersion {
		return config, parser, fmt.Errorf(AppVersion())
	}

	// Load additional config from file
	err = flags.NewIniParser(parser).ParseFile(config.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return config, parser, fmt.Errorf("Error parsing configuration file: %v", err.Error())
		}
		return config, parser, err
	}

	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil && !IsFlagErrorType(err, flags.ErrHelp) {
		return config, parser, err
	}

	return config, parser, nil
}
