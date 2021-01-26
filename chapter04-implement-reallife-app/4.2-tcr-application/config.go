// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import "os"

// IConfig is interface for application config
type IConfig interface {
	ServiceID() string
	CacheServer() string
	MQServers() string
	CitizenRegisteredTopic() string
	CitizenConfirmedTopic() string
	CitizenValidationAPI() string
	CitizenDeliveryAPI() string
	BatchDeliverAPI() string
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

// CitizenRegisteredTopic return topic name for registered event
func (cfg *Config) CitizenRegisteredTopic() string {
	return "when-citizen-has-registered"
}

// CitizenConfirmedTopic return topic name for confirmed event
func (cfg *Config) CitizenConfirmedTopic() string {
	return "when-citizen-has-confirmed"
}

// CitizenValidationAPI return API to validate citizen information
func (cfg *Config) CitizenValidationAPI() string {
	return "http://external-api:8080/3rd-party/validate"
}

// CitizenDeliveryAPI return API to request delivery citizen ID card
func (cfg *Config) CitizenDeliveryAPI() string {
	return "http://external-api:8080/3rd-party/delivery"
}

// BatchDeliverAPI return API to batch delivery citizen ID card
func (cfg *Config) BatchDeliverAPI() string {
	return "http://batch-ptask-api:8080/ptask/delivery"
}
