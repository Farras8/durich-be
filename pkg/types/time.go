package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// NullTime is a wrapper around time.Time that implements JSON marshalling with consistent format
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the sql.Scanner interface
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}
	nt.Valid = true
	switch v := value.(type) {
	case time.Time:
		nt.Time = v
		return nil
	case []byte:
		return nt.Time.UnmarshalText(v)
	case string:
		return nt.Time.UnmarshalText([]byte(v))
	default:
		return fmt.Errorf("cannot scan type %T into NullTime: %v", value, value)
	}
}

// Value implements the driver.Valuer interface
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON implements the json.Marshaler interface
// Returns RFC3339 format like created_at and updated_at
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, nt.Time.Format(time.RFC3339))), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}

	// Remove quotes
	str := string(data)
	if len(str) > 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	nt.Time = t
	nt.Valid = true
	return nil
}
