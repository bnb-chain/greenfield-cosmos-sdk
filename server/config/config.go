package config

import (
	"errors"
	"fmt"
	"math"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"

	clientflags "github.com/cosmos/cosmos-sdk/client/flags"
	pruningtypes "github.com/cosmos/cosmos-sdk/pruning/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	defaultMinGasPrices = ""

	// DefaultAPIAddress defines the default address to bind the API server to.
	DefaultAPIAddress = "tcp://0.0.0.0:1317"

	// DefaultGRPCAddress defines the default address to bind the gRPC server to.
	DefaultGRPCAddress = "0.0.0.0:9090"

	// DefaultGRPCWebAddress defines the default address to bind the gRPC-web server to.
	DefaultGRPCWebAddress = "0.0.0.0:9091"

	// DefaultJSONRPCAddress is the default address the JSON-RPC server binds to.
	DefaultJSONRPCAddress = "0.0.0.0:8545"

	// DefaultJSONRPCWsAddress is the default address the JSON-RPC WebSocket server binds to.
	DefaultJSONRPCWsAddress = "0.0.0.0:8546"

	// DefaultGRPCMaxRecvMsgSize defines the default gRPC max message size in
	// bytes the server can receive.
	DefaultGRPCMaxRecvMsgSize = 1024 * 1024 * 10

	// DefaultGRPCMaxSendMsgSize defines the default gRPC max message size in
	// bytes the server can send.
	DefaultGRPCMaxSendMsgSize = math.MaxInt32

	DefaultHTTPTimeout = 30 * time.Second

	DefaultHTTPIdleTimeout = 120 * time.Second

	// DefaultMaxOpenConnections represents the amount of open connections (unlimited = 0)
	DefaultMaxOpenConnections = 0
)

// BaseConfig defines the server's basic configuration
type BaseConfig struct {
	// The minimum gas prices a validator is willing to accept for processing a
	// transaction. A transaction's fees must meet the minimum of any denomination
	// specified in this config (e.g. 0.25token1;0.0001token2).
	MinGasPrices string `mapstructure:"minimum-gas-prices"`

	Pruning           string `mapstructure:"pruning"`
	PruningKeepRecent string `mapstructure:"pruning-keep-recent"`
	PruningInterval   string `mapstructure:"pruning-interval"`

	// HaltHeight contains a non-zero block height at which a node will gracefully
	// halt and shutdown that can be used to assist upgrades and testing.
	//
	// Note: Commitment of state will be attempted on the corresponding block.
	HaltHeight uint64 `mapstructure:"halt-height"`

	// HaltTime contains a non-zero minimum block time (in Unix seconds) at which
	// a node will gracefully halt and shutdown that can be used to assist
	// upgrades and testing.
	//
	// Note: Commitment of state will be attempted on the corresponding block.
	HaltTime uint64 `mapstructure:"halt-time"`

	// MinRetainBlocks defines the minimum block height offset from the current
	// block being committed, such that blocks past this offset may be pruned
	// from Tendermint. It is used as part of the process of determining the
	// ResponseCommit.RetainHeight value during ABCI Commit. A value of 0 indicates
	// that no blocks should be pruned.
	//
	// This configuration value is only responsible for pruning Tendermint blocks.
	// It has no bearing on application state pruning which is determined by the
	// "pruning-*" configurations.
	//
	// Note: Tendermint block pruning is dependant on this parameter in conunction
	// with the unbonding (safety threshold) period, state pruning and state sync
	// snapshot parameters to determine the correct minimum value of
	// ResponseCommit.RetainHeight.
	MinRetainBlocks uint64 `mapstructure:"min-retain-blocks"`

	// InterBlockCache enables inter-block caching.
	InterBlockCache bool `mapstructure:"inter-block-cache"`

	// IndexEvents defines the set of events in the form {eventType}.{attributeKey},
	// which informs Tendermint what to index. If empty, all events will be indexed.
	IndexEvents []string `mapstructure:"index-events"`

	// IavlCacheSize set the size of the iavl tree cache.
	IAVLCacheSize uint64 `mapstructure:"iavl-cache-size"`

	// IAVLDisableFastNode enables or disables the fast sync node.
	IAVLDisableFastNode bool `mapstructure:"iavl-disable-fastnode"`

	// AppDBBackend defines the type of Database to use for the application and snapshots databases.
	// An empty string indicates that the Tendermint config's DBBackend value should be used.
	AppDBBackend string `mapstructure:"app-db-backend"`
}

// APIConfig defines the API listener configuration.
type APIConfig struct {
	// Enable defines if the API server should be enabled.
	Enable bool `mapstructure:"enable"`

	// Swagger defines if swagger documentation should automatically be registered.
	Swagger bool `mapstructure:"swagger"`

	// EnableUnsafeCORS defines if CORS should be enabled (unsafe - use it at your own risk)
	EnableUnsafeCORS bool `mapstructure:"enabled-unsafe-cors"`

	// Address defines the API server to listen on
	Address string `mapstructure:"address"`

	// MaxOpenConnections defines the number of maximum open connections
	MaxOpenConnections uint `mapstructure:"max-open-connections"`

	// RPCReadTimeout defines the Tendermint RPC read timeout (in seconds)
	RPCReadTimeout uint `mapstructure:"rpc-read-timeout"`

	// RPCWriteTimeout defines the Tendermint RPC write timeout (in seconds)
	RPCWriteTimeout uint `mapstructure:"rpc-write-timeout"`

	// RPCMaxBodyBytes defines the Tendermint maximum response body (in bytes)
	RPCMaxBodyBytes uint `mapstructure:"rpc-max-body-bytes"`

	// TODO: TLS/Proxy configuration.
	//
	// Ref: https://github.com/cosmos/cosmos-sdk/issues/6420
}

// RosettaConfig defines the Rosetta API listener configuration.
type RosettaConfig struct {
	// Address defines the API server to listen on
	Address string `mapstructure:"address"`

	// Blockchain defines the blockchain name
	// defaults to DefaultBlockchain
	Blockchain string `mapstructure:"blockchain"`

	// Network defines the network name
	Network string `mapstructure:"network"`

	// Retries defines the maximum number of retries
	// rosetta will do before quitting
	Retries int `mapstructure:"retries"`

	// Enable defines if the API server should be enabled.
	Enable bool `mapstructure:"enable"`

	// Offline defines if the server must be run in offline mode
	Offline bool `mapstructure:"offline"`

	// EnableFeeSuggestion defines if the server should suggest fee by default
	EnableFeeSuggestion bool `mapstructure:"enable-fee-suggestion"`

	// GasToSuggest defines gas limit when calculating the fee
	GasToSuggest int `mapstructure:"gas-to-suggest"`

	// DenomToSuggest defines the defult denom for fee suggestion
	DenomToSuggest string `mapstructure:"denom-to-suggest"`
}

// GRPCConfig defines configuration for the gRPC server.
type GRPCConfig struct {
	// Enable defines if the gRPC server should be enabled.
	Enable bool `mapstructure:"enable"`

	// Address defines the API server to listen on
	Address string `mapstructure:"address"`

	// MaxRecvMsgSize defines the max message size in bytes the server can receive.
	// The default value is 10MB.
	MaxRecvMsgSize int `mapstructure:"max-recv-msg-size"`

	// MaxSendMsgSize defines the max message size in bytes the server can send.
	// The default value is math.MaxInt32.
	MaxSendMsgSize int `mapstructure:"max-send-msg-size"`
}

// GRPCWebConfig defines configuration for the gRPC-web server.
type GRPCWebConfig struct {
	// Enable defines if the gRPC-web should be enabled.
	Enable bool `mapstructure:"enable"`

	// Address defines the gRPC-web server to listen on
	Address string `mapstructure:"address"`

	// EnableUnsafeCORS defines if CORS should be enabled (unsafe - use it at your own risk)
	EnableUnsafeCORS bool `mapstructure:"enable-unsafe-cors"`
}

// StateSyncConfig defines the state sync snapshot configuration.
type StateSyncConfig struct {
	// SnapshotInterval sets the interval at which state sync snapshots are taken.
	// 0 disables snapshots.
	SnapshotInterval uint64 `mapstructure:"snapshot-interval"`

	// SnapshotKeepRecent sets the number of recent state sync snapshots to keep.
	// 0 keeps all snapshots.
	SnapshotKeepRecent uint32 `mapstructure:"snapshot-keep-recent"`
}

// JSONRPCConfig defines configuration for the EVM RPC server.
type JSONRPCConfig struct {
	// API defines a list of JSON-RPC namespaces that should be enabled
	API []string `mapstructure:"api"`
	// Address defines the HTTP server to listen on
	Address string `mapstructure:"address"`
	// WsAddress defines the WebSocket server to listen on
	WsAddress string `mapstructure:"ws-address"`
	// Enable defines if the EVM RPC server should be enabled.
	Enable bool `mapstructure:"enable"`
	// HTTPTimeout is the read/write timeout of http json-rpc server.
	HTTPTimeout time.Duration `mapstructure:"http-timeout"`
	// HTTPIdleTimeout is the idle timeout of http json-rpc server.
	HTTPIdleTimeout time.Duration `mapstructure:"http-idle-timeout"`
	// MaxOpenConnections sets the maximum number of simultaneous connections
	// for the server listener.
	MaxOpenConnections int `mapstructure:"max-open-connections"`
}

// TLSConfig defines the certificate and matching private key for the server.
type TLSConfig struct {
	// CertificatePath the file path for the certificate .pem file
	CertificatePath string `mapstructure:"certificate-path"`
	// KeyPath the file path for the key .pem file
	KeyPath string `mapstructure:"key-path"`
}

// UpgradeConfig defines the upgrading configuration.
type UpgradeConfig struct {
	Name   string `mapstructure:"name"`
	Height int64  `mapstructure:"height"`
	Info   string `mapstructure:"info"`
}

// Config defines the server's top level configuration
type Config struct {
	BaseConfig `mapstructure:",squash"`

	Upgrade []UpgradeConfig `mapstructure:"upgrade"`
	// Telemetry defines the application telemetry configuration
	Telemetry telemetry.Config `mapstructure:"telemetry"`
	API       APIConfig        `mapstructure:"api"`
	GRPC      GRPCConfig       `mapstructure:"grpc"`
	Rosetta   RosettaConfig    `mapstructure:"rosetta"`
	GRPCWeb   GRPCWebConfig    `mapstructure:"grpc-web"`
	StateSync StateSyncConfig  `mapstructure:"state-sync"`
	JSONRPC   JSONRPCConfig    `mapstructure:"json-rpc"`
	TLS       TLSConfig        `mapstructure:"tls"`
}

// SetMinGasPrices sets the validator's minimum gas prices.
func (c *Config) SetMinGasPrices(gasPrices sdk.DecCoins) {
	c.MinGasPrices = gasPrices.String()
}

// GetMinGasPrices returns the validator's minimum gas prices based on the set
// configuration.
func (c *Config) GetMinGasPrices() sdk.DecCoins {
	if c.MinGasPrices == "" {
		return sdk.DecCoins{}
	}

	gasPricesStr := strings.Split(c.MinGasPrices, ";")
	gasPrices := make(sdk.DecCoins, len(gasPricesStr))

	for i, s := range gasPricesStr {
		gasPrice, err := sdk.ParseDecCoin(s)
		if err != nil {
			panic(fmt.Errorf("failed to parse minimum gas price coin (%s): %s", s, err))
		}

		gasPrices[i] = gasPrice
	}

	return gasPrices
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: BaseConfig{
			MinGasPrices:        defaultMinGasPrices,
			InterBlockCache:     true,
			Pruning:             pruningtypes.PruningOptionDefault,
			PruningKeepRecent:   "0",
			PruningInterval:     "0",
			MinRetainBlocks:     0,
			IndexEvents:         make([]string, 0),
			IAVLCacheSize:       781250, // 50 MB
			IAVLDisableFastNode: false,
			AppDBBackend:        "",
		},
		Telemetry: telemetry.Config{
			Enabled:      false,
			GlobalLabels: [][]string{},
		},
		API: APIConfig{
			Enable:             false,
			Swagger:            false,
			Address:            DefaultAPIAddress,
			MaxOpenConnections: 1000,
			RPCReadTimeout:     10,
			RPCMaxBodyBytes:    1000000,
		},
		GRPC: GRPCConfig{
			Enable:         true,
			Address:        DefaultGRPCAddress,
			MaxRecvMsgSize: DefaultGRPCMaxRecvMsgSize,
			MaxSendMsgSize: DefaultGRPCMaxSendMsgSize,
		},
		Rosetta: RosettaConfig{
			Enable:              false,
			Address:             ":8080",
			Blockchain:          "app",
			Network:             "network",
			Retries:             3,
			Offline:             false,
			EnableFeeSuggestion: false,
			GasToSuggest:        clientflags.DefaultGasLimit,
			DenomToSuggest:      "uatom",
		},
		GRPCWeb: GRPCWebConfig{
			Enable:  true,
			Address: DefaultGRPCWebAddress,
		},
		StateSync: StateSyncConfig{
			SnapshotInterval:   0,
			SnapshotKeepRecent: 2,
		},
		JSONRPC: *DefaultJSONRPCConfig(),
		TLS:     *DefaultTLSConfig(),
	}
}

// GetConfig returns a fully parsed Config object.
func GetConfig(v *viper.Viper) (Config, error) {
	globalLabelsRaw, ok := v.Get("telemetry.global-labels").([]interface{})
	if !ok {
		return Config{}, fmt.Errorf("failed to parse global-labels config")
	}

	globalLabels := make([][]string, 0, len(globalLabelsRaw))
	for idx, glr := range globalLabelsRaw {
		labelsRaw, ok := glr.([]interface{})
		if !ok {
			return Config{}, fmt.Errorf("failed to parse global label number %d from config", idx)
		}
		if len(labelsRaw) == 2 {
			globalLabels = append(globalLabels, []string{labelsRaw[0].(string), labelsRaw[1].(string)})
		}
	}

	var upgrade []UpgradeConfig
	v.UnmarshalKey("upgrade", &upgrade)
	return Config{
		BaseConfig: BaseConfig{
			MinGasPrices:        v.GetString("minimum-gas-prices"),
			InterBlockCache:     v.GetBool("inter-block-cache"),
			Pruning:             v.GetString("pruning"),
			PruningKeepRecent:   v.GetString("pruning-keep-recent"),
			PruningInterval:     v.GetString("pruning-interval"),
			HaltHeight:          v.GetUint64("halt-height"),
			HaltTime:            v.GetUint64("halt-time"),
			IndexEvents:         v.GetStringSlice("index-events"),
			MinRetainBlocks:     v.GetUint64("min-retain-blocks"),
			IAVLCacheSize:       v.GetUint64("iavl-cache-size"),
			IAVLDisableFastNode: v.GetBool("iavl-disable-fastnode"),
			AppDBBackend:        v.GetString("app-db-backend"),
		},
		Upgrade: upgrade,
		Telemetry: telemetry.Config{
			ServiceName:             v.GetString("telemetry.service-name"),
			Enabled:                 v.GetBool("telemetry.enabled"),
			EnableHostname:          v.GetBool("telemetry.enable-hostname"),
			EnableHostnameLabel:     v.GetBool("telemetry.enable-hostname-label"),
			EnableServiceLabel:      v.GetBool("telemetry.enable-service-label"),
			PrometheusRetentionTime: v.GetInt64("telemetry.prometheus-retention-time"),
			GlobalLabels:            globalLabels,
		},
		API: APIConfig{
			Enable:             v.GetBool("api.enable"),
			Swagger:            v.GetBool("api.swagger"),
			Address:            v.GetString("api.address"),
			MaxOpenConnections: v.GetUint("api.max-open-connections"),
			RPCReadTimeout:     v.GetUint("api.rpc-read-timeout"),
			RPCWriteTimeout:    v.GetUint("api.rpc-write-timeout"),
			RPCMaxBodyBytes:    v.GetUint("api.rpc-max-body-bytes"),
			EnableUnsafeCORS:   v.GetBool("api.enabled-unsafe-cors"),
		},
		Rosetta: RosettaConfig{
			Enable:              v.GetBool("rosetta.enable"),
			Address:             v.GetString("rosetta.address"),
			Blockchain:          v.GetString("rosetta.blockchain"),
			Network:             v.GetString("rosetta.network"),
			Retries:             v.GetInt("rosetta.retries"),
			Offline:             v.GetBool("rosetta.offline"),
			EnableFeeSuggestion: v.GetBool("rosetta.enable-fee-suggestion"),
			GasToSuggest:        v.GetInt("rosetta.gas-to-suggest"),
			DenomToSuggest:      v.GetString("rosetta.denom-to-suggest"),
		},
		GRPC: GRPCConfig{
			Enable:         v.GetBool("grpc.enable"),
			Address:        v.GetString("grpc.address"),
			MaxRecvMsgSize: v.GetInt("grpc.max-recv-msg-size"),
			MaxSendMsgSize: v.GetInt("grpc.max-send-msg-size"),
		},
		GRPCWeb: GRPCWebConfig{
			Enable:           v.GetBool("grpc-web.enable"),
			Address:          v.GetString("grpc-web.address"),
			EnableUnsafeCORS: v.GetBool("grpc-web.enable-unsafe-cors"),
		},
		StateSync: StateSyncConfig{
			SnapshotInterval:   v.GetUint64("state-sync.snapshot-interval"),
			SnapshotKeepRecent: v.GetUint32("state-sync.snapshot-keep-recent"),
		},
		JSONRPC: JSONRPCConfig{
			API:                v.GetStringSlice("json-rpc.api"),
			Address:            v.GetString("json-rpc.address"),
			WsAddress:          v.GetString("json-rpc.ws-address"),
			Enable:             v.GetBool("json-rpc.enable"),
			HTTPTimeout:        v.GetDuration("json-rpc.http-timeout"),
			HTTPIdleTimeout:    v.GetDuration("json-rpc.http-idle-timeout"),
			MaxOpenConnections: v.GetInt("json-rpc.max-open-connections"),
		},
		TLS: TLSConfig{
			CertificatePath: v.GetString("tls.certificate-path"),
			KeyPath:         v.GetString("tls.key-path"),
		},
	}, nil
}

// ValidateBasic returns an error if something is wrong with the config.
func (c Config) ValidateBasic() error {
	if c.BaseConfig.MinGasPrices == "" {
		return sdkerrors.ErrAppConfig.Wrap("set min gas price in app.toml or flag or env variable")
	}
	if c.Pruning == pruningtypes.PruningOptionEverything && c.StateSync.SnapshotInterval > 0 {
		return sdkerrors.ErrAppConfig.Wrapf(
			"cannot enable state sync snapshots with '%s' pruning setting", pruningtypes.PruningOptionEverything,
		)
	}

	if err := c.JSONRPC.Validate(); err != nil {
		return sdkerrors.ErrAppConfig.Wrapf("invalid json-rpc config value: %s", err.Error())
	}

	if err := c.TLS.Validate(); err != nil {
		return sdkerrors.ErrAppConfig.Wrapf("invalid tls config value: %s", err.Error())
	}

	return nil
}

// GetDefaultAPINamespaces returns the default list of JSON-RPC namespaces that should be enabled
func GetDefaultAPINamespaces() []string {
	return []string{"eth", "net"}
}

// DefaultJSONRPCConfig returns an EVM config with the JSON-RPC API enabled by default
func DefaultJSONRPCConfig() *JSONRPCConfig {
	return &JSONRPCConfig{
		Enable:             true,
		API:                GetDefaultAPINamespaces(),
		Address:            DefaultJSONRPCAddress,
		WsAddress:          DefaultJSONRPCWsAddress,
		HTTPTimeout:        DefaultHTTPTimeout,
		HTTPIdleTimeout:    DefaultHTTPIdleTimeout,
		MaxOpenConnections: DefaultMaxOpenConnections,
	}
}

// Validate returns an error if the JSON-RPC configuration fields are invalid.
func (c JSONRPCConfig) Validate() error {
	if c.Enable && len(c.API) == 0 {
		return errors.New("cannot enable JSON-RPC without defining any API namespace")
	}

	if c.HTTPTimeout < 0 {
		return errors.New("JSON-RPC HTTP timeout duration cannot be negative")
	}

	if c.HTTPIdleTimeout < 0 {
		return errors.New("JSON-RPC HTTP idle timeout duration cannot be negative")
	}

	// check for duplicates
	seenAPIs := make(map[string]bool)
	for _, api := range c.API {
		if seenAPIs[api] {
			return fmt.Errorf("repeated API namespace '%s'", api)
		}

		seenAPIs[api] = true
	}

	return nil
}

// DefaultTLSConfig returns the default TLS configuration
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		CertificatePath: "",
		KeyPath:         "",
	}
}

// Validate returns an error if the TLS certificate and key file extensions are invalid.
func (c TLSConfig) Validate() error {
	certExt := path.Ext(c.CertificatePath)

	if c.CertificatePath != "" && certExt != ".pem" {
		return fmt.Errorf("invalid extension %s for certificate path %s, expected '.pem'", certExt, c.CertificatePath)
	}

	keyExt := path.Ext(c.KeyPath)

	if c.KeyPath != "" && keyExt != ".pem" {
		return fmt.Errorf("invalid extension %s for key path %s, expected '.pem'", keyExt, c.KeyPath)
	}

	return nil
}
