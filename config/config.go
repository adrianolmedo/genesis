package config

import (
	"encoding/json"
	"net"
	"os"
	"strings"
)

// Driver type for enum driver engines in func main.
type Driver string

// Server configuration for RESTful API.
type Config struct {
	// LocalHost set true if you want the server runing on 127.0.0.1 by default,
	// false if you want the server run it using IPv4 address (local IP).
	LocalHost bool `json:"localhost"`

	// Port server, if is empty by default will are 8080.
	Port string `json:"port"`

	// CORS directive, add address separated by comma.
	CORS string `json:"cors"`

	Database `json:"database"`
}

// Database config.
type Database struct {
	// Engine eg.: "mysql" or "postgres".
	Engine Driver `json:"engine"`

	// Server when is running the database Engine.
	Server string `json:"server"`

	// Port of database Engine server.
	Port string `json:"port"`

	// User of database, eg.: "root".
	User string `json:"user"`

	// Password of User database
	Password string `json:"password"`

	// Name of SQL database.
	Name string `json:"name"`
}

// Init load .json file configuration from root and dump it in Config structure.
func Init(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	err = json.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	cfg.CORS = strings.Join(strings.Fields(cfg.CORS), "") // remove whitespaces

	return cfg, nil
}

// Address return address server eg.: "127.0.0.1:8080".
// If the "localhost" field is true in config.json, the IP address it will be 127.0.0.1 by default,
// otherwise it will try to take the IPv4 (if exists), eg.: "192.168.0.107:8080".
func (cfg *Config) Address() string {
	IP := "127.0.0.1"
	if !cfg.LocalHost {
		IP = getHostIP()
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return IP + ":" + cfg.Port
}

// getHostIP return local IP. If you are not connected to IPv4 it will return empty string.
func getHostIP() string {
	netInterfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIP, ok := netInterfaceAddress.(*net.IPNet)

		if ok && !networkIP.IP.IsLoopback() && networkIP.IP.To4() != nil {
			ip := networkIP.IP.String()
			return ip
		}
	}
	return ""
}
