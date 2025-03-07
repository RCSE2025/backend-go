package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type DateOnly struct {
	time.Time
}

func (d DateOnly) Value() (driver.Value, error) {
	return d.Format("2006-01-02"), nil // Store only the date part
}

func (date *DateOnly) Scan(value interface{}) error {
	scanned, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprint("Failed to scan DateOnly value:", value))
	}
	*date = DateOnly{scanned}
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.Format("2006-01-02"))), nil
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(`"2006-01-02"`, string(data))
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}
