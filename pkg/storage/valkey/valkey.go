package valkeyStorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

// Storage interface that is implemented by storage providers.
type Storage struct {
	client valkey.Client
}

// New creates a new valkey storage.
func New(config ...Config) *Storage {
	// Set default config
	cfg := configDefault(config...)

	// Create new valkey universal client
	var client valkey.Client

	// Parse the URL and update config values accordingly
	if cfg.URL != "" {
		options, err := valkey.ParseURL(cfg.URL)
		if err != nil {
			panic(err)
		}

		// Update the config values with the parsed URL values
		cfg.Username = options.Username
		cfg.Password = options.Password
		cfg.Database = options.SelectDB
		cfg.Addrs = options.InitAddress

		// If cfg.TLSConfig is not provided, and options returns one, use it.
		if cfg.TLSConfig == nil && options.TLSConfig != nil {
			cfg.TLSConfig = options.TLSConfig
		}
	} else if len(cfg.Addrs) == 0 {
		// Fallback to Host and Port values if Addrs is empty
		cfg.Addrs = []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}
	}

	// Create Universal Client
	client, _ = valkey.NewClient(valkey.ClientOption{
		InitAddress: cfg.Addrs,
		//MasterName:       cfg.MasterName,
		ClientName: cfg.ClientName,
		Sentinel:   valkey.SentinelOption{Username: cfg.SentinelUsername, Password: cfg.SentinelPassword},
		SelectDB:   cfg.Database,
		Username:   cfg.Username,
		Password:   cfg.Password,
		TLSConfig:  cfg.TLSConfig,
		//PoolSize:         cfg.PoolSize,
	})

	// Test connection
	if err := client.Do(context.Background(), client.B().Ping().Build()).Error(); err != nil {
		panic(err)
	}

	// Empty collection if Clear is true
	if cfg.Reset {
		if err := client.Do(context.Background(), client.B().Flushdb().Build()).Error(); err != nil {
			panic(err)
		}
	}

	// Create new store
	return &Storage{
		client: client,
	}
}

// Get value by key.
func (s *Storage) Get(key string) ([]byte, error) {
	if len(key) == 0 {
		return nil, nil
	}

	val, err := s.client.Do(context.Background(), s.client.B().Get().Key(key).Build()).AsBytes()
	if errors.Is(err, valkey.Nil) {
		return nil, nil
	}

	return val, err
}

// Set key with value.
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	if len(key) == 0 || len(val) == 0 {
		return nil
	}

	err := s.client.Do(context.Background(), s.client.B().Set().Key(key).Value(string(val)).Ex(exp).Build()).Error()
	return err
}

// Delete key by key.
func (s *Storage) Delete(key string) error {
	if len(key) == 0 {
		return nil
	}
	return s.client.Do(context.Background(), s.client.B().Del().Key(key).Build()).Error()
}

// Reset all keys.
func (s *Storage) Reset() error {
	return s.client.Do(context.Background(), s.client.B().Flushdb().Build()).Error()
}

// Close the database.
func (s *Storage) Close() error {
	s.client.Close()
	return nil
}

// Return database client.
func (s *Storage) Conn() valkey.Client {
	return s.client
}

// Return all the keys.
func (s *Storage) Keys() ([]string, error) {
	var keys []string
	var err error
	var resp valkey.ScanEntry

	for {
		if resp, err = s.client.Do(context.Background(),
			s.client.B().Scan().Cursor(resp.Cursor).Match("*").Count(10).Build()).AsScanEntry(); err != nil {
			return nil, err
		}

		keys = append(keys, resp.Elements...)

		if resp.Cursor == 0 {
			break
		}
	}

	if len(resp.Elements) == 0 {
		return nil, nil
	}

	return keys, nil
}
