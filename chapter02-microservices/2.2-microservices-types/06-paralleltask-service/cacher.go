// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

// ICacher is interface for redis cache
type ICacher interface {
	Set(key string, value interface{}, expire time.Duration) error
	SetS(key string, value string, expire time.Duration) error
	Get(key string) (string, error)
	HasChanged(key string, value string) (bool, error)
}

// Cacher implement ICacher to connect with Redis
type Cacher struct {
	server string
	client *redis.Client
}

// NewCacher return new instance of Cacher
func NewCacher(server string) *Cacher {
	return &Cacher{
		server: server,
	}
}

// Set object into cache
func (cache *Cacher) Set(key string, value interface{}, expire time.Duration) error {
	c, err := cache.getClient()
	if err != nil {
		return err
	}

	str, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.Set(key, str, expire).Err()
	if err != nil {
		return err
	}

	return nil
}

// SetS set string into cache
func (cache *Cacher) SetS(key string, value string, expire time.Duration) error {
	c, err := cache.getClient()
	if err != nil {
		return err
	}

	err = c.Set(key, value, expire).Err()
	if err != nil {
		return err
	}

	return nil
}

// Get object from cache
func (cache *Cacher) Get(key string) (string, error) {
	c, err := cache.getClient()
	if err != nil {
		return "", err
	}

	val, err := c.Get(key).Result()
	if err == redis.Nil {
		// Key does not exists
		return "", nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}

// HasChanged detect if value of key has changed it will return true
// If get and error it will return true with error
// If get the same value it will return false
func (cache *Cacher) HasChanged(key string, value string) (bool, error) {
	current, err := cache.Get(key)
	if err != nil {
		return true, err
	}
	if current != value {
		return true, nil
	}
	return false, nil
}

func (cache *Cacher) getClient() (*redis.Client, error) {
	client := cache.client
	if client == nil {
		client = cache.newClient(cache.server)
		cache.client = client
	}

	retry := 3 // Retry connecting 3 times, if cannot ping
	for true {
		_, err := client.Ping().Result()
		if err != nil {
			retry--
			if retry < 0 {
				return nil, err
			}

			// Wait before reconnecting
			continue
		}
		// If we can PING without error, just return
		break
	}

	return client, nil
}

func (cache *Cacher) newClient(server string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: server,
		DB:   0,
	})
}
