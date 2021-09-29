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
type Server struct {
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

// NewServer load .json file configuration from root and dump it in Server structure.
func NewServer(path string) (*Server, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	srv := new(Server)
	err = json.NewDecoder(file).Decode(srv)
	if err != nil {
		return nil, err
	}

	if srv.Port == "" {
		srv.Port = "8080"
	}

	srv.CORS = strings.Join(strings.Fields(srv.CORS), "") // remove whitespaces

	return srv, nil
}

// Address return string for address server eg.: "127.0.0.1:8080",
// if the "localhost" field is true in config.json, the IP address it will be 127.0.0.1 by default
// otherwise it will try to take the IPv4 (if exists), eg.: "192.168.0.107:8080".
func (srv *Server) Address() string {
	IP := "127.0.0.1"
	if !srv.LocalHost {
		IP = getHostIP()
	}

	if srv.Port == "" {
		srv.Port = "8080"
	}

	return IP + ":" + srv.Port
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
