package config

import (
	"crypto/ecdsa"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/pion/webrtc/v2"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/Kdag-K/kdag/src/common"
	"github.com/Kdag-K/kdag/src/proxy"
)

const (
	// DefaultKeyfile defines the default name of the file containing the
	// validator's private key
	DefaultKeyfile = "priv_key"

	// DefaultBadgerFile defines the default name of the folder containing the
	// Badger database
	DefaultBadgerFile = "badger_db"
	DefaultCertFile = "cert.pem"
)

// Default configuration values.
const (
	DefaultLogLevel             = "debug"
	DefaultBindAddr             = "127.0.0.1:1337"
	DefaultServiceAddr          = "127.0.0.1:8000"
	DefaultHeartbeatTimeout     = 10 * time.Millisecond
	DefaultSlowHeartbeatTimeout = 1000 * time.Millisecond
	DefaultTCPTimeout           = 1000 * time.Millisecond
	DefaultJoinTimeout          = 10000 * time.Millisecond
	DefaultCacheSize            = 10000
	DefaultSyncLimit            = 1000
	DefaultMaxPool              = 2
	DefaultStore                = false
	DefaultMaintenanceMode      = false
	DefaultSuspendLimit         = 100
	DefaultWebRTC               = false
	DefaultSignalAddr           = "127.0.0.1:2443"
	DefaultSignalRealm          = "office"
	DefaultSignalSkipVerify     = false
)

// Config contains all the configuration properties of a Kdag node.
type Config struct {
	// DataDir is the top-level directory containing Kdag configuration and
	// data
	DataDir string `mapstructure:"datadir"`

	// LogLevel determines the chattiness of the log output.
	LogLevel string `mapstructure:"log"`

	// BindAddr is the local address:port where this node gossips with other
	// nodes. By default, this is "0.0.0.0", meaning Kdag will bind to all
	// addresses on the local machine. However, in some cases, there may be a
	// routable address that cannot be bound. Use AdvertiseAddr to enable
	// gossiping a different address to support this. If this address is not
	// routable, the node will be in a constant flapping state as other nodes
	// will treat the non-routability as a failure
	BindAddr string `mapstructure:"listen"`

	// AdvertiseAddr is used to change the address that we advertise to other
	// nodes in the cluster
	AdvertiseAddr string `mapstructure:"advertise"`

	// NoService disables the HTTP API service.
	NoService bool `mapstructure:"no-service"`

	// ServiceAddr is the address:port that serves the user-facing API. If not
	// specified, and "no-service" is not set, the API handlers are registered
	// with the DefaultServerMux of the http package. It is possible that
	// another server in the same process is simultaneously using the
	// DefaultServerMux. In which case, the handlers will be accessible from
	// both servers. This is usefull when Kdag is used in-memory and expecpted
	// to use the same endpoint (address:port) as the application's API.
	ServiceAddr string `mapstructure:"service-listen"`

	// HeartbeatTimeout is the frequency of the gossip timer when the node has
	// something to gossip about.
	HeartbeatTimeout time.Duration `mapstructure:"heartbeat"`

	// SlowHeartbeatTimeout is the frequency of the gossip timer when the node
	// has nothing to gossip about.
	SlowHeartbeatTimeout time.Duration `mapstructure:"slow-heartbeat"`

	// MaxPool controls how many connections are pooled per target in the gossip
	// routines.
	MaxPool int `mapstructure:"max-pool"`

	// TCPTimeout is the timeout of gossip TCP connections.
	TCPTimeout time.Duration `mapstructure:"timeout"`

	// JoinTimeout is the timeout of Join Requests
	JoinTimeout time.Duration `mapstructure:"join_timeout"`

	// SyncLimit defines the max number of hashgraph events to include in a
	// SyncResponse or EagerSyncRequest
	SyncLimit int `mapstructure:"sync-limit"`

	// EnableFastSync determines whether or not to enable the FastSync protocol.
	EnableFastSync bool `mapstructure:"fast-sync"`

	// Store is a flag that determines whether or not to use persistant storage.
	Store bool `mapstructure:"store"`

	// DatabaseDir is the directory containing database files.
	DatabaseDir string `mapstructure:"db"`

	// CacheSize is the max number of items in in-memory caches.
	CacheSize int `mapstructure:"cache-size"`

	// Bootstrap determines whether or not to load Kdag from an existing
	// database file. Forces Store, ie. bootstrap only works with a persistant
	// database store.
	Bootstrap bool `mapstructure:"bootstrap"`

	// MaintenanceMode when set to true causes Kdag to initialise in a
	// suspended state. I.e. it does not start gossipping. Forces Bootstrap,
	// which itself forces Store. I.e. MaintenanceMode only works if the node is
	// bootstrapped from an existing database.
	MaintenanceMode bool `mapstructure:"maintenance-mode"`

	// SuspendLimit is the number of Undetermined Events (Events which haven't
	// reached consensus) that will cause the node to become suspended
	SuspendLimit int `mapstructure:"suspend-limit"`

	// Moniker defines the friendly name of this node
	Moniker string `mapstructure:"moniker"`

	// LoadPeers determines whether or not to attempt loading the peer-set from
	// a local json file.
	WebRTC bool `mapstructure:"webrtc"`

	// SignalAddr is the IP:PORT of the WebRTC signaling server. It is ignored
	// when WebRTC is not enabled. The connection is over secured web-sockets,
	// wss, and it possible to include a self-signed certificated in a file
	// called cert.pem in the datadir. If no self-signed certificate is found,
	// the server's certifacate signing authority better be trusted.
	SignalAddr string `mapstructure:"signal-addr"`

	// SignalRealm is an administrative domain within the WebRTC signaling
	// server. WebRTC signaling messages are only routed within a Realm.
	SignalRealm string `mapstructure:"signal-realm"`

	// SignalSkipVerify controls whether the signal client verifies the server's
	// certificate chain and host name. If SignalSkipVerify is true, TLS accepts
	// any certificate presented by the server and any host name in that
	// certificate. In this mode, TLS is susceptible to man-in-the-middle
	// attacks. This should be used only for testing.
	SignalSkipVerify bool `mapstructure:"signal-skip-verify"`

	// ICEServers defines a slice describing servers available to be used by
	// ICE, such as STUN and TURN servers.
	// https://developer.mozilla.org/en-US/docs/Web/API/RTCIceServer/urls
	ICEServers []webrtc.ICEServer

	// Proxy is the application proxy that enables Kdag to communicate with
	// application.
	Proxy proxy.AppGateway

	// Key is the private key of the validator.
	Key *ecdsa.PrivateKey

	logger *logrus.Logger
}

// NewDefaultConfig returns the a config object with default values.
func NewDefaultConfig() *Config {
	config := &Config{
		DataDir:              DefaultDataDir(),
		LogLevel:             DefaultLogLevel,
		BindAddr:             DefaultBindAddr,
		ServiceAddr:          DefaultServiceAddr,
		HeartbeatTimeout:     DefaultHeartbeatTimeout,
		SlowHeartbeatTimeout: DefaultSlowHeartbeatTimeout,
		TCPTimeout:           DefaultTCPTimeout,
		JoinTimeout:          DefaultJoinTimeout,
		CacheSize:            DefaultCacheSize,
		SyncLimit:            DefaultSyncLimit,
		MaxPool:              DefaultMaxPool,
		Store:                DefaultStore,
		MaintenanceMode:      DefaultMaintenanceMode,
		DatabaseDir:          DefaultDatabaseDir(),
		SuspendLimit:         DefaultSuspendLimit,
		WebRTC:               DefaultWebRTC,
		SignalAddr:           DefaultSignalAddr,
		SignalRealm:          DefaultSignalRealm,
		SignalSkipVerify:     DefaultSignalSkipVerify,
		ICEServers:           DefaultICEServers(),
	}

	return config
}

// NewTestConfig returns a config object with default values and a special
// logger. the logger forces formatting and colors even when there is no tty
// attached, which makes for more readable logs. The logger also provides info
// about the calling function.
func NewTestConfig(t testing.TB, level logrus.Level) *Config {
	config := NewDefaultConfig()
	config.logger = common.NewTestLogger(t, level)
	return config
}

// SetDataDir sets the top-level Kdag directory, and updates the database
// directory if it is currently set to the default value. If the database
// directory is not currently the default, it means the user has explicitely set
// it to something else, so avoid changing it again here.
func (c *Config) SetDataDir(dataDir string) {
	c.DataDir = dataDir
	if c.DatabaseDir == DefaultDatabaseDir() {
		c.DatabaseDir = filepath.Join(dataDir, DefaultBadgerFile)
	}
}

// Keyfile returns the full path of the file containing the private key.
func (c *Config) Keyfile() string {
	return filepath.Join(c.DataDir, DefaultKeyfile)
}

// Logger returns a formatted logrus Entry, with prefix set to "kdag".
func (c *Config) CertFile() string {
	return filepath.Join(c.DataDir, DefaultCertFile)
}
func (c *Config) Logger() *logrus.Entry {
	if c.logger == nil {
		c.logger = logrus.New()
		c.logger.Level = LogLevel(c.LogLevel)
		c.logger.Formatter = new(prefixed.TextFormatter)
	}
	return c.logger.WithField("prefix", "kdag")
}

// DefaultDatabaseDir returns the default path for the badger database files.
func DefaultDatabaseDir() string {
	return filepath.Join(DefaultDataDir(), DefaultBadgerFile)
}

// DefaultDataDir return the default directory name for top-level Kdag config
// based on the underlying OS, attempting to respect conventions.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".Kdag")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Kdag")
		} else {
			return filepath.Join(home, ".kdag")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

// HomeDir returns the user's home directory.
func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

// LogLevel parses a string into a Logrus log level.
func LogLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}

// DefaultICEServers returns a list containing a single ICEServer which
// points to a public STUN server provided by Google. This default configuration
// does not include a TURN server, so not all p2p connections will be possible.
func DefaultICEServers() []webrtc.ICEServer {
	return []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	}
}
