package config

import "reflect"

// GooglePlay is the generated configuration struct for the connector.
type GooglePlay struct {
	ServiceAccountKeyPath string `mapstructure:"service-account-key-path"`
	DeveloperID           string `mapstructure:"developer-id"`
}

func (c *GooglePlay) findFieldByTag(tagValue string) (any, bool) {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("mapstructure")
		if tag == tagValue {
			return v.Field(i).Interface(), true
		}
	}
	return nil, false
}

func (c *GooglePlay) GetStringSlice(fieldName string) []string {
	v, ok := c.findFieldByTag(fieldName)
	if !ok {
		return []string{}
	}
	t, ok := v.([]string)
	if !ok {
		panic("wrong type")
	}
	return t
}

func (c *GooglePlay) GetString(fieldName string) string {
	v, ok := c.findFieldByTag(fieldName)
	if !ok {
		return ""
	}
	t, ok := v.(string)
	if !ok {
		panic("wrong type")
	}
	return t
}

func (c *GooglePlay) GetInt(fieldName string) int {
	v, ok := c.findFieldByTag(fieldName)
	if !ok {
		return 0
	}
	t, ok := v.(int)
	if !ok {
		panic("wrong type")
	}
	return t
}

func (c *GooglePlay) GetBool(fieldName string) bool {
	v, ok := c.findFieldByTag(fieldName)
	if !ok {
		return false
	}
	t, ok := v.(bool)
	if !ok {
		panic("wrong type")
	}
	return t
}

func (c *GooglePlay) GetStringMap(fieldName string) map[string]any {
	v, ok := c.findFieldByTag(fieldName)
	if !ok {
		return map[string]any{}
	}
	t, ok := v.(map[string]any)
	if !ok {
		panic("wrong type")
	}
	return t
}
