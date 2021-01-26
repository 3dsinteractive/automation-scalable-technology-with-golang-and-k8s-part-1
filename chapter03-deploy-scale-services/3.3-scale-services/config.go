// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import "os"

// IConfig is interface for application config
type IConfig interface {
	ServiceID() string
	CacheServer() string
	MQServers() string
}

// Config implement IConfig
type Config struct{}

// NewConfig return new config instance
func NewConfig() *Config {
	return &Config{}
}

// ServiceID return ID of service
func (cfg *Config) ServiceID() string {
	return os.Getenv("SERVICE_ID")
}

// CacheServer return redis server
func (cfg *Config) CacheServer() string {
	return os.Getenv("CACHE_SERVER")
}

// MQServers return Kafka servers
func (cfg *Config) MQServers() string {
	return os.Getenv("MQ_SERVERS")
}
