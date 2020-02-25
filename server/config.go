// Package server contains interfaces and functions dealing with setting up proxy core,
// including code construct the module pipeline.
package server

import (
	"os"
	"strconv"
	"time"
)

//Config describe parameters need to make a connection to a Mongo database
type Config struct {
	Scheme    string        `json:"scheme"`
	Hosts     string        `json:"hosts"`
	TLS       bool          `json:"tls"`
	Database  string        `json:"database"`
	Username  string        `json:"-"`
	Password  string        `json:"-"`
	Timeout   time.Duration `json:"timeout"`
	OptParams string        `json:"optParams"`
	ReadOnly  bool          `json:"readOnly"`
	Port      int           `json:"port"`
}

// FromEnv populates Config from the environment
func (c *Config) FromEnv() {
	c.Scheme = os.Getenv("SCHEME")
	c.Hosts = os.Getenv("ADDRESSES")
	c.Username = os.Getenv("USERNAME")
	c.Password = os.Getenv("PASSWORD")
	c.Database = os.Getenv("DATABASE")
	c.OptParams = os.Getenv("OPT_PARAMS")
	c.TLS = os.Getenv("TLS") == "true"
	timeoutStr := os.Getenv("TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if (timeoutStr == "") || (err != nil) {
		c.Timeout = time.Duration(20 * time.Second)
	} else {
		c.Timeout = time.Duration(timeout) * time.Second
	}
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if (portStr == "") || (err != nil) {
		c.Port = 27017
	} else {
		c.Port = port
	}
	c.ReadOnly = os.Getenv("READONLY") == "true"
}

// AsConnectionString constructs a MongoDB connection string from a Config
func (c *Config) AsConnectionString() string {
	var url string
	if c.Scheme != "" {
		url += c.Scheme + "://"
	} else {
		url += "mongodb://"
	}

	if c.Username != "" {
		url += c.Username
		if c.Password != "" {
			url += ":"
			url += c.Password
		}
		url += "@"
	}
	url += c.Hosts
	url += "/"
	url += c.Database
	if c.TLS {
		url += "?ssl=true"
	} else {
		url += "?ssl=false"
	}
	if c.OptParams != "" {
		url += c.OptParams
	}

	return url
}