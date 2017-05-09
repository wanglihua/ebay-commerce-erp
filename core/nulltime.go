package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/revel/revel"
)

// NullDate
//------------------------------------------------------------------------------------------------------------------

type NullDate struct {
	Time  time.Time
	Valid bool
}

func (t *NullDate) Scan(value interface{}) error {
	switch s := value.(type) {
	case time.Time:
		t.Time = s
		t.Valid = true
		return nil
	case string:
		//x, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
		x, err := Parse_time(s)
		if err != nil {
			return err
		}
		t.Time = x
		t.Valid = true
		return nil
	case []uint8:
		//x, err := time.ParseInLocation("2006-01-02 15:04:05", string(s), time.Local)
		x, err := Parse_time(string(s))
		if err != nil {
			return err
		}
		t.Time = x
		t.Valid = true
		return nil
	case nil:
		t.Valid = false
		return nil
	default:
		revel.ERROR.Printf("null: cannot scan type %T into NullLongTime: %v\n", value, value)
		return fmt.Errorf("null: cannot scan type %T into NullLongTime: %v\n", value, value)
	}

	return nil
}

func (t NullDate) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.Local().Format("2006-01-02 15:04:05"), nil
}

func MakeNullDate(t time.Time, valid bool) NullDate {
	return NullDate{
		Time:  t,
		Valid: valid,
	}
}

func NullDateFrom(t time.Time) NullDate {
	return MakeNullDate(t, true)
}

func NullDateFromPtr(t *time.Time) NullDate {
	if t == nil {
		return MakeNullDate(time.Time{}, false)
	}
	return MakeNullDate(*t, true)
}

func (t NullDate) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	//return t.Time.MarshalJSON()
	return []byte(`"` + t.Time.Local().Format("2006-01-02") + `"`), nil

}

func (t *NullDate) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		//	err = t.Time.UnmarshalJSON(data)
		t.UnmarshalText(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			revel.ERROR.Printf(`json: unmarshalling object into Go value of type NullDate requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively\n`, x["Time"], x["Valid"])
			return fmt.Errorf(`json: unmarshalling object into Go value of type NullDate requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively\n`, x["Time"], x["Valid"])
		}
		err = t.Time.UnmarshalText([]byte(ti))
		t.Valid = valid
		return err
	case nil:
		t.Valid = false
		return nil
	default:
		revel.ERROR.Printf("json: cannot unmarshal %v into Go value of type null.Time\n", reflect.TypeOf(v).Name())
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time\n", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err
}

func (t NullDate) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalText()
}

func (t *NullDate) UnmarshalText(text []byte) error {
	str := strings.Trim(string(text), `"`)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

func (t *NullDate) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

func (t NullDate) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// NullTime
//------------------------------------------------------------------------------------------------------------------

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (t *NullTime) Scan(value interface{}) error {
	switch s := value.(type) {
	case time.Time:
		t.Time = s
		t.Valid = true
		return nil
	case string:
		//x, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
		x, err := Parse_time(s)
		if err != nil {
			return err
		}
		t.Time = x
		t.Valid = true
		return nil
	case []uint8:
		//x, err := time.ParseInLocation("2006-01-02 15:04:05", string(s), time.Local)
		x, err := Parse_time(string(s))
		if err != nil {
			return err
		}
		t.Time = x
		t.Valid = true
		return nil
	case nil:
		t.Valid = false
		return nil
	default:
		revel.ERROR.Printf("null: cannot scan type %T into NullLongTime: %v\n", value, value)
		return fmt.Errorf("null: cannot scan type %T into NullLongTime: %v\n", value, value)
	}

	return nil
}

func (t NullTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

func MakeNullTime(t time.Time, valid bool) NullTime {
	return NullTime{
		Time:  t,
		Valid: valid,
	}
}

func NullTimeFrom(t time.Time) NullTime {
	return MakeNullTime(t, true)
}

func NullTimeFromPtr(t *time.Time) NullTime {
	if t == nil {
		return MakeNullTime(time.Time{}, false)
	}
	return MakeNullTime(*t, true)
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	//return t.Time.MarshalJSON()
	return []byte(`"` + t.Time.Local().Format("2006-01-02 15:04") + `"`), nil

}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		//	err = t.Time.UnmarshalJSON(data)
		t.UnmarshalText(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			revel.ERROR.Printf(`json: unmarshalling object into Go value of type NullTime requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively\n`, x["Time"], x["Valid"])
			return fmt.Errorf(`json: unmarshalling object into Go value of type NullTime requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively\n`, x["Time"], x["Valid"])
		}
		err = t.Time.UnmarshalText([]byte(ti))
		t.Valid = valid
		return err
	case nil:
		t.Valid = false
		return nil
	default:
		revel.ERROR.Printf("json: cannot unmarshal %v into Go value of type null.Time\n", reflect.TypeOf(v).Name())
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time\n", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err
}

func (t NullTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalText()
}

func (t *NullTime) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

func (t *NullTime) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

func (t NullTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// NullLongTime
//------------------------------------------------------------------------------------------------------------------

type NullLongTime struct {
	Time  time.Time
	Valid bool
}

func (t *NullLongTime) Scan(value interface{}) error {
	switch s := value.(type) {
	case time.Time:
		t.Time = s
		t.Valid = true
		return nil
	case string:
		//x, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
		x, err := Parse_time(s)
		if err != nil {
			return err
		}
		t.Time = x
		t.Valid = true
		return nil
	case []uint8:
		//x, err := time.ParseInLocation("2006-01-02 15:04:05", string(s), time.Local)
		x, err := Parse_time(string(s))
		if err != nil {
			return err
		}
		t.Time = x
		t.Valid = true
		return nil
	case nil:
		t.Valid = false
		return nil
	default:
		revel.ERROR.Printf("null: cannot scan type %T into NullLongTime: %v\n", value, value)
		return fmt.Errorf("null: cannot scan type %T into NullLongTime: %v\n", value, value)
	}

	return nil
}

func (t NullLongTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.Local().Format("2006-01-02 15:04:05"), nil
}

func MakeNullLongTime(t time.Time, valid bool) NullLongTime {
	return NullLongTime{
		Time:  t,
		Valid: valid,
	}
}

func NullLongTimeFrom(t time.Time) NullLongTime {
	return MakeNullLongTime(t, true)
}

func NullLongTimeFromPtr(t *time.Time) NullLongTime {
	if t == nil {
		return MakeNullLongTime(time.Time{}, false)
	}
	return MakeNullLongTime(*t, true)
}

func (t NullLongTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	//return t.Time.MarshalJSON()
	return []byte(`"` + t.Time.Local().Format("2006-01-02 15:04:05") + `"`), nil

}

func (t *NullLongTime) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		//	err = t.Time.UnmarshalJSON(data)
		t.UnmarshalText(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			revel.ERROR.Printf(`json: unmarshalling object into Go value of type NullLongTime requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively\n`, x["Time"], x["Valid"])
			return fmt.Errorf(`json: unmarshalling object into Go value of type NullLongTime requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively\n`, x["Time"], x["Valid"])
		}
		err = t.Time.UnmarshalText([]byte(ti))
		t.Valid = valid
		return err
	case nil:
		t.Valid = false
		return nil
	default:
		revel.ERROR.Printf("json: cannot unmarshal %v into Go value of type null.Time\n", reflect.TypeOf(v).Name())
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time\n", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err
}

func (t NullLongTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalText()
}

func (t *NullLongTime) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

func (t *NullLongTime) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

func (t NullLongTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
