package hera

import (
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sevenNt/hera/file"
	"github.com/spf13/cast"
)

var cfg = New()
var loadedFileMap = make(map[string]bool)

// Provider provides configuration meta-data.
type Provider interface {
	Read() (map[string]interface{}, error)
	Watch(func(map[string]interface{}))
}

// Config provides configuration for application.
type Config struct {
	provider Provider
	mu       sync.RWMutex
	override map[string]interface{}
	keyDelim string

	plugins []Plugin
}

const (
	defaultKeyDelim = "."
)

// New constructs a new Config with provider.
func New() *Config {
	return &Config{
		override: make(map[string]interface{}),
		keyDelim: defaultKeyDelim,
	}
}

// Reset resets all to default settings.
func Reset() {
	cfg = New()
	loadedFileMap = make(map[string]bool)
}

// SetKeyDelim set keyDelim of a cfg instance.
func (c *Config) SetKeyDelim(delim string) {
	c.keyDelim = delim
}

// Sub returns new Config instance representing a sub tree of this instance.
func (c *Config) Sub(key string) *Config {
	return &Config{
		keyDelim: c.keyDelim,
		override: c.GetStringMap(key),
	}
}

// LoadFromFile loads configuration from file.
func LoadFromFile(path string, watch bool) error {
	if _, loaded := loadedFileMap[path]; loaded {
		return nil
	}
	err := Load(file.NewProvider(path, watch))
	if err == nil {
		loadedFileMap[path] = true
	}

	return err
}

// MustLoadFromFile is like LoadFromFile but panics if the configuration file cannot be parsed.
func MustLoadFromFile(path string, watch bool) {
	if _, loaded := loadedFileMap[path]; loaded {
		return
	}
	err := Load(file.NewProvider(path, watch))
	if err != nil {
		panic(err)
	}
	loadedFileMap[path] = true
}

// RegisterPlugin Registers callback plugin.
func RegisterPlugin(p Plugin) {
	cfg.plugins = append(cfg.plugins, p)
}

// Load loads configuration from provided provider with default cfg.
func Load(p Provider) error {
	return cfg.Load(p)
}

// Load loads configuration from provided provider.
func (c *Config) Load(p Provider) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, err := p.Read()
	if err != nil {
		return err
	}

	mergeStringMap(c.override, data)
	return nil
}

// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface{} {
	return cfg.Get(key)
}

// Get returns the value associated with the key
func (c *Config) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.find(key)
}

// GetString returns the value associated with the key as a string with default cfg.
func GetString(key string) string {
	return cfg.GetString(key)
}

// GetString returns the value associated with the key as a string.
func (c *Config) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

// GetBool returns the value associated with the key as a boolean with default cfg.
func GetBool(key string) bool {
	return cfg.GetBool(key)
}

// GetBool returns the value associated with the key as a boolean.
func (c *Config) GetBool(key string) bool {
	return cast.ToBool(c.Get(key))
}

// GetInt returns the value associated with the key as an integer with default cfg.
func GetInt(key string) int {
	return cfg.GetInt(key)
}

// GetInt returns the value associated with the key as an integer.
func (c *Config) GetInt(key string) int {
	return cast.ToInt(c.Get(key))
}

// GetInt64 returns the value associated with the key as an integer with default cfg.
func GetInt64(key string) int64 {
	return cfg.GetInt64(key)
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Config) GetInt64(key string) int64 {
	return cast.ToInt64(c.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64 with default cfg.
func GetFloat64(key string) float64 {
	return cfg.GetFloat64(key)
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.Get(key))
}

// GetTime returns the value associated with the key as time with default cfg.
func GetTime(key string) time.Time {
	return cfg.GetTime(key)
}

// GetTime returns the value associated with the key as time.
func (c *Config) GetTime(key string) time.Time {
	return cast.ToTime(c.Get(key))
}

// GetDuration returns the value associated with the key as a duration with default cfg.
func GetDuration(key string) time.Duration {
	return cfg.GetDuration(key)
}

// GetDuration returns the value associated with the key as a duration.
func (c *Config) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings with default cfg.
func GetStringSlice(key string) []string {
	return cfg.GetStringSlice(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.Get(key))
}

// GetSlice returns the value associated with the key as a slice of strings with default cfg.
func GetSlice(key string) []interface{} {
	return cfg.GetSlice(key)
}

// GetSlice returns the value associated with the key as a slice of strings.
func (c *Config) GetSlice(key string) []interface{} {
	return cast.ToSlice(c.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces with default cfg.
func GetStringMap(key string) map[string]interface{} {
	return cfg.GetStringMap(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Config) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(c.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings with default cfg.
func GetStringMapString(key string) map[string]string {
	return cfg.GetStringMapString(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Config) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(c.Get(key))
}

// GetSliceStringMap returns the value associated with the slice of maps.
//func (c *Config) GetSliceStringMap(key string) []map[string]interface{} {
//	return cast.ToSliceStringMap(c.Get(key))
//}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings with default cfg.
func GetStringMapStringSlice(key string) map[string][]string {
	return cfg.GetStringMapStringSlice(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Config) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(c.Get(key))
}

// UnmarshalKey takes a single key and unmarshals it into a Struct with default cfg.
func UnmarshalKey(key string, rawVal interface{}) error {
	return cfg.UnmarshalKey(key, rawVal)
}

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     rawVal,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	if key == "" {
		c.mu.RLock()
		defer c.mu.RUnlock()
		return decoder.Decode(c.override)
	}
	return decoder.Decode(c.Get(key))
}

func (c *Config) find(key string) interface{} {
	paths := strings.Split(key, c.keyDelim)
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := deepSearchInMap(c.override, paths[:len(paths)-1]...)
	return m[paths[len(paths)-1]]
}
