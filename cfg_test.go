package hera

import (
	"fmt"
	"testing"
	"time"

	"github.com/sevenNt/hera/file"
	"github.com/stretchr/testify/assert"
)

func initToml() {
	Reset()
	MustLoadFromFile("./example/cfg.toml", false)
}

func initYaml() {
	Reset()
	MustLoadFromFile("./example/cfg.yaml", false)
}

func initJSON() {
	Reset()
	MustLoadFromFile("./example/cfg.json", false)
}

func TestLoadFromFile(t *testing.T) {
	assert.NoError(t, LoadFromFile("./example/cfg.toml", false))
	assert.NoError(t, LoadFromFile("./example/cfg.json", false))
	assert.NoError(t, LoadFromFile("./example/cfg.yaml", false))
	assert.Error(t, LoadFromFile("./example/cfg-123.toml", false))
}

func TestToml(t *testing.T) {
	initToml()
	assert.Equal(t, true, Get("debug"))
	assert.Equal(t, []interface{}{"127.0.0.1:2379", "127.0.0.1:2479"}, Get("app.registry.etcd.endpoints"))
}

func TestYaml(t *testing.T) {
	initYaml()
	assert.Equal(t, "brown", GetString("a.b.c"))
}

func TestJson(t *testing.T) {
	initJSON()
	assert.Equal(t, "brown", Get("a.b.c"))
}

func TestGet(t *testing.T) {
	initToml()
	assert.Equal(t, true, Get("debug"))
}

func TestGetString(t *testing.T) {
	initToml()
	assert.Equal(t, "./log", GetString("log.dir"))
}

func TestGetBool(t *testing.T) {
	initToml()
	assert.Equal(t, true, GetBool("debug"))
}

func TestGetInt(t *testing.T) {
	initToml()
	assert.Equal(t, 51020, GetInt("server.http.port"))
}

func TestGetInt64(t *testing.T) {
	initToml()
	assert.Equal(t, int64(51020), GetInt64("server.http.port"))
}

func TestGetFloat64(t *testing.T) {
	initToml()
	assert.Equal(t, 64.09, GetFloat64("coordinate.longitude"))
}

func TestGetTime(t *testing.T) {
	initToml()
	tm, err := time.Parse(time.RFC3339, "2015-01-01T20:17:05Z")
	assert.NoError(t, err)
	assert.Equal(t, tm, GetTime("date"))
}

func TestGetDuration(t *testing.T) {
	initToml()
	assert.Equal(t, 2*time.Second, GetDuration("app.registry.etcd.timeout"))
}

func TestGetStringSlice(t *testing.T) {
	initToml()
	assert.Equal(t, []string{"127.0.0.1:2379", "127.0.0.1:2479"}, GetStringSlice("app.registry.etcd.endpoints"))
}

func TestGetSlice(t *testing.T) {
	initToml()
	assert.Equal(t, []interface{}{"127.0.0.1:2379", "127.0.0.1:2479"}, GetSlice("app.registry.etcd.endpoints"))
}

func TestGetStringMap(t *testing.T) {
	initToml()
	assert.Equal(t, map[string]interface{}{
		"dir":   "./log",
		"level": "Info | Warn | Error | Panic | Fatal",
	}, GetStringMap("log"))
}

func TestGetStringMapString(t *testing.T) {
	initToml()
	assert.Equal(t, map[string]string{
		"dir":   "./log",
		"level": "Info | Warn | Error | Panic | Fatal",
	}, GetStringMapString("log"))
}

func TestGetStringMapSlice(t *testing.T) {
	initToml()
	assert.Equal(t, map[string][]string{
		"Hubei":     {"Wuhan", "Tianmen"},
		"Guangdong": {"Guangzhou"},
	}, GetStringMapStringSlice("province"))
}

func TestUnmarshalKey(t *testing.T) {
	initToml()
	type Log struct {
		Dir   string
		Level string
	}
	var log Log
	UnmarshalKey("log", &log)
	assert.Equal(t, Log{Dir: "./log", Level: "Info | Warn | Error | Panic | Fatal"}, log)
}

func Test_GetStringMapSliceTOML(t *testing.T) {
	config := New()
	config.Load(file.NewProvider("/tmp/toml.toml", false))
	if rep := config.GetStringMapStringSlice("province"); rep != nil {
		fmt.Printf("rep: %#v\n", rep)
	}
}
func Test_GetStringMapSliceYAML(t *testing.T) {
	config := New()
	config.Load(file.NewProvider("/tmp/yaml.yaml", false))
	if rep := config.GetStringMapStringSlice("province"); rep != nil {
		fmt.Printf("rep: %#v\n", rep)
	}
}
