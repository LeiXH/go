package types

import (
	"database/sql"
	"fmt"
	"time"
)

const ctlayoutWithsec = "2006-01-02 15:04:05"

type SQLBool bool
type SQLTimestamp time.Time
type SQLNullInt sql.NullInt64
type SQLNullString sql.NullString

func (b *SQLBool) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("unexpected type for SQLBool: %T, %#v", src, src)
	}
	switch str[0] {
	case 0x0:
		*b = SQLBool(false)
	case 0x1:
		*b = SQLBool(true)
	}
	return nil
}

func (t *SQLTimestamp) Scan(src interface{}) error {
	_str, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("unexpected type for SQLTimestamp: %T, %#v", src, src)
	}
	str := string(_str)
	// log.Printf("get data %#v", string(str))
	v, err := time.Parse(ctlayoutWithsec, str)

	if err != nil {
		return err
	}
	*t = SQLTimestamp(v)
	return nil
}
func (t SQLTimestamp) MarshalJSON() ([]byte, error) {
	a := []byte(fmt.Sprintf(`"%s"`, time.Time(t).Format(ctlayoutWithsec)))
	return a, nil
}

func (i *SQLNullInt) Scan(value interface{}) error {
	si := sql.NullInt64{}
	if err := si.Scan(value); err == nil {
		i.Int64 = si.Int64
		i.Valid = si.Valid
	}
	return nil
}
func (i SQLNullInt) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return []byte(fmt.Sprintf("%d", i.Int64)), nil
	} else {
		return []byte("0"), nil
	}
}
func (i *SQLNullInt) SetValue(val int64) {
	i.Int64 = val
	i.Valid = true
}

func (s *SQLNullString) Scan(value interface{}) error {
	ss := sql.NullString{}
	if err := ss.Scan(value); err == nil {
		s.String = ss.String
		s.Valid = ss.Valid
	}
	return nil
}
func (s SQLNullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return []byte(fmt.Sprintf(`"%s"`, s.String)), nil
	} else {
		return []byte(`""`), nil
	}
}
func (s *SQLNullString) SetValue(val string) {
	s.String = val
	s.Valid = true
}
