package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgtype"
)

const (
	ISO8601        = "2006-01-02T15:04:05"
	FormatDDMMYYYY = "02.01.2006"
)

type DBTime time.Time

func (c *DBTime) Scan(value interface{}) error {
	var t pgtype.Timestamp
	if err := t.Scan(value); err != nil {
		return fmt.Errorf("couldn't scan the timestamp: %w", err)
	}

	*c = DBTime(t.Time)

	return nil
}

func (c *DBTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "" {
		return
	}

	t, err := time.Parse(ISO8601, s)
	if err != nil {
		return fmt.Errorf("time parse error , %w", err)
	}

	*c = DBTime(t)

	return nil
}

func (c *DBTime) MarshalJSON() ([]byte, error) {
	if c == nil || time.Time(*c).IsZero() {
		return nil, nil
	}

	return []byte(fmt.Sprintf("%q", time.Time(*c).Format(ISO8601))), nil
}

type PlannedShippingDate time.Time
