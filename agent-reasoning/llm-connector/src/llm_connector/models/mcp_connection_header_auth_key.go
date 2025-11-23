package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// MCPConnectionHeaderAuthKey holds the schema definition for the MCPConnectionHeaderAuthKey.
type MCPConnectionHeaderAuthKey string

func (t *MCPConnectionHeaderAuthKey) String() string {
	return string(*t)
}

type MCPConnectionHeaderAuthKeys []MCPConnectionHeaderAuthKey

func (t *MCPConnectionHeaderAuthKeys) FromArray(values []string) {
	*t = make(MCPConnectionHeaderAuthKeys, len(values))
	for i, v := range values {
		(*t)[i] = MCPConnectionHeaderAuthKey(v)
	}
}

func (t MCPConnectionHeaderAuthKeys) Value() (driver.Value, error) {
	// Convert to []string for storage
	strs := make([]string, len(t))
	for i, v := range t {
		strs[i] = string(v)
	}
	// Store as JSON array
	return json.Marshal(strs)
}

func (t *MCPConnectionHeaderAuthKeys) Scan(value interface{}) error {
	if value == nil {
		*t = MCPConnectionHeaderAuthKeys{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var strs []string
	if err := json.Unmarshal(b, &strs); err != nil {
		return err
	}
	t.FromArray(strs)
	return nil
}
