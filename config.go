package gobase

import (
	"errors"
	"github.com/BurntSushi/toml"
	"strings"
)

type GoConfig struct {
	ConfFile  string
	Version   string
	PathLevel int
	Item      map[string]interface{}
}

func NewConfig(file string, level int) (*GoConfig, error) {
	var tmp interface{}
	if _, err := toml.DecodeFile(file, &tmp); err != nil {
		return nil, err
	}

	c := new(GoConfig)
	c.ConfFile = file
	c.Version = "0.1.1"
	c.PathLevel = level
	c.Item = make(map[string]interface{})
	c.loadConfig(tmp, []string{})

	return c, nil
}

func (c *GoConfig) ReloadConfig() error {
	var tmp interface{}
	if _, err := toml.DecodeFile(c.ConfFile, &tmp); err != nil {
		return err
	}

	c.Item = make(map[string]interface{})
	c.loadConfig(tmp, []string{})
	return nil
}

func (c *GoConfig) loadConfig(tree interface{}, path []string) {
	if c.PathLevel > 0 && len(path) >= c.PathLevel {
		return
	}
	for key, value := range tree.(map[string]interface{}) {
		fullPath := append(path, key)
		pathKey := strings.Join(fullPath, ".")
		switch orig := value.(type) {
		case map[string]interface{}:
			c.loadConfig(orig, fullPath)
		default:
			c.Item[pathKey] = orig
		}

	}
}

func (c *GoConfig) Int(key string, def int) (int, error) {
	value, ok := c.Item[key]
	if !ok {
		return def, nil
	}
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	}
	return def, errors.New("Type Not Match Int")
}

func (c *GoConfig) Int64(key string, def int64) (int64, error) {
	value, ok := c.Item[key]
	if !ok {
		return def, nil
	}
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	}
	return def, errors.New("Type Not Match Int64")
}

func (c *GoConfig) Float64(key string, def float64) (float64, error) {
	value, ok := c.Item[key]
	if !ok {
		return def, nil
	}
	switch v := value.(type) {
	case float64:
		return v, nil
	}
	return def, errors.New("Type Not Match Float64")
}

func (c *GoConfig) String(key string, def string) (string, error) {
	value, ok := c.Item[key]
	if !ok {
		return def, nil
	}

	switch str := value.(type) {
	case string:
		return str, nil
	}
	return def, errors.New("Type Not Match String")
}

func (c *GoConfig) Bool(key string, def bool) (bool, error) {
	value, ok := c.Item[key]
	if !ok {
		return def, nil
	}
	switch v := value.(type) {
	case bool:
		return v, nil
	}
	return def, errors.New("Type Not Match Bool")
}

func (c *GoConfig) Array(key string) ([]interface{}, error) {
	value, ok := c.Item[key]
	if !ok {
		return nil, nil
	}
	switch v := value.(type) {
	case []interface{}:
		return v, nil
	}
	return nil, errors.New("Type Not Match Array")
}

func (c *GoConfig) Map(key string) (map[string]interface{}, error) {
	value, ok := c.Item[key]
	if !ok {
		return nil, nil
	}
	switch v := value.(type) {
	case map[string]interface{}:
		return v, nil
	}
	return nil, errors.New("Type Not Match Map")
}
