package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// TODO ? API spec related, show null or not
// TODO ? API spec related, show 0 for certain types
// TODO Bithday issue ns.Time.IsZero, Birthday & different time zones? 9am Saving & returning to front-end
// TODO Do we really need Sqlx & its Scan?

type NullInt64 struct {
	Int64 int64
	Valid bool
}

func (n *NullInt64) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}

	if v, ok := value.(int64); ok {
		n.Int64, n.Valid = v, false
		return nil
	}

	return nil
}

func (n NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

func (ns NullInt64) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Int64)
	}
	return []byte("null"), nil
}

// ZeroOmitInt64 positive int > 0, MarshalJSON ignore null & 0
type ZeroOmitInt64 struct {
	Int64 int64
	Valid bool
}

func (n *ZeroOmitInt64) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}

	if v, ok := value.(int64); ok && v > 0 {
		n.Int64, n.Valid = v, true
		return nil
	}

	return nil
}

func (n ZeroOmitInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

func (ns ZeroOmitInt64) MarshalJSON() ([]byte, error) {
	if ns.Valid && ns.Int64 > 0 {
		return json.Marshal(ns.Int64)
	}
	return []byte("null"), nil
}
func (ns *ZeroOmitInt64) UnmarshalJSON(b []byte) error {
	// TODO test edge scenario
	var s int64
	var ss string
	if err := json.Unmarshal(b, &s); err != nil {
		if err := json.Unmarshal(b, &ss); err != nil {
			ns.Int64, ns.Valid = 0, true
			return nil
		}
		i, err := strconv.ParseInt(ss, 10, 64)
		if err != nil {
			ns.Int64, ns.Valid = 0, true
			return nil
		}
		s = i
	}
	ns.Int64 = s
	ns.Valid = true
	return nil
}


// NullBool is an alias for sql.NullBool data type
type NullBool sql.NullBool

// Scan implements the Scanner interface for NullBool
func (nb *NullBool) Scan(value interface{}) error {
	var b sql.NullBool
	if err := b.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nb = NullBool{b.Bool, false}
	} else {
		*nb = NullBool{b.Bool, true}
	}

	return nil
}


// NullString MarshalJSON ignores ""
type NullString struct {
	String string
	Valid  bool
}

func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	// TODO review and check perfomance with PGX
	ns.String = fmt.Sprintf("%v", value)
	return nil
}

func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid || len(ns.String) == 0 {
		return "", nil
	}
	return ns.String, nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid && ns.String != "" {
		return json.Marshal(ns.String)
	}
	return []byte("null"), nil
}

func (ns *NullString) UnmarshalJSON(b []byte) error {

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ns.String = s
	ns.Valid = true
	return nil
}


func SetStr(s string) *NullString {
	return &NullString{String: s, Valid:true}
}
func SetBool(b bool) *NullBool {
	return &NullBool{Bool: b, Valid:true}
}
func SetStrBool(s string) *NullBoolString {
	return &NullBoolString{String: s, Valid:true}
}

func SetDateStrTime(s string) *NullTime {
	//layout := "2018-10-01T15:04:05"
	str := s+"T00:00:00.001Z"
	if len(s)==4 {
		str = s+"-01-01T00:00:00.001Z"
	}

	fmt.Println("SetStrTime:", str)

	t, _ := time.Parse(time.RFC3339, str)

	return &NullTime{Time: t, Valid:true}
}
func SetStrFromInt61(i int64) *NullString {
	s := strconv.FormatInt(i, 10)
	return &NullString{String: s, Valid:true}
}
func SetInt64(i int64) *ZeroOmitInt64 {
	return &ZeroOmitInt64{Int64: i, Valid:true}
}
func SetInt64FromStr(s string) *ZeroOmitInt64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", n, n)
	}
	return &ZeroOmitInt64{Int64: n, Valid:true}
}
func SetNow() *NullTime {
	return &NullTime{Time: time.Now(), Valid:true}
}


// NullString MarshalJSON ignores ""
type NullBoolString struct {
	String string
	Valid  bool
}

func (ns *NullBoolString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "0", false
		return nil
	}
	ns.Valid = true
	ns.String = "1"
	return nil
}

func (ns NullBoolString) Value() (driver.Value, error) {
	if !ns.Valid || len(ns.String) == 0 {
		return "0", nil
	}
	return "1", nil
}

func (ns NullBoolString) MarshalJSON() ([]byte, error) {
	if ns.Valid && ns.String != "" {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

func (ns *NullBoolString) UnmarshalJSON(b []byte) error {
	var s string
	var bb bool
	var i int64
	if err := json.Unmarshal(b, &s); err != nil {
		if err := json.Unmarshal(b, &bb); err != nil {
			if err := json.Unmarshal(b, &i); err != nil {
				return err
			}
			if i>0 {
				s = "1"
			} else {
				s = "0"
			}
		}
		if bb {
			s = "1"
		} else {
			s = "0"
		}
	}

	ns.String = s
	ns.Valid = true
	return nil
}





// NullTime MarshalJSON ignores ZeroTime
type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	if _, ok := value.(time.Time); ok && !value.(time.Time).IsZero() {
		nt.Time, nt.Valid = value.(time.Time), true
		return nil
	}
	t := new(time.Time) // zero time
	nt.Time, nt.Valid = *t, false
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (ns NullTime) MarshalJSON() ([]byte, error) {
	if ns.Valid && !ns.Time.IsZero() {
		//return json.Marshal(ns.Time.UnixNano() / 1e6)
		return json.Marshal(ns.Time.Format("2006-01-02T15:04:05.000Z"))
		//time.RFC3339
	}
	return []byte("null"), nil
}

func (ns *NullTime) UnmarshalJSON(b []byte) error {


	if bytes.Equal(b, []byte("null")) || string(b) == "" || string(b) == `""` {
		ns.Time = time.Time{}
		ns.Valid = false
		return nil
	}

	var s time.Time
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ns.Time = s
	ns.Valid = true

	return nil
}

// JSONText is a json.RawMessage, which is a []byte underneath.
// Value() validates the json format in the source, and returns an error if
// the json is not valid.  Scan does no validation.  JSONText additionally
// implements `Unmarshal`, which unmarshals the json within to an interface{}
type JSONText json.RawMessage

var emptyJSON = JSONText("{}")

// MarshalJSON returns the *j as the JSON encoding of j.
func (j JSONText) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return emptyJSON, nil
	}
	return j, nil
}

// UnmarshalJSON sets *j to a copy of data
func (j *JSONText) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSONText: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// Value returns j as a value.  This does a validating unmarshal into another
// RawMessage.  If j is invalid json, it returns an error.
func (j JSONText) Value() (driver.Value, error) {
	var m json.RawMessage
	var err = j.Unmarshal(&m)
	if err != nil {
		return []byte{}, err
	}
	return []byte(j), nil
}

// Scan stores the src in *j.  No validation is done.
func (j *JSONText) Scan(src interface{}) error {
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		if len(t) == 0 {
			source = emptyJSON
		} else {
			source = t
		}
	case nil:
		*j = emptyJSON
	default:
		return errors.New("Incompatible type for JSONText")
	}
	*j = JSONText(append((*j)[0:0], source...))
	return nil
}

// Unmarshal unmarshal's the json in j to v, as in json.Unmarshal.
func (j *JSONText) Unmarshal(v interface{}) error {
	if len(*j) == 0 {
		*j = emptyJSON
	}
	return json.Unmarshal([]byte(*j), v)
}

// String supports pretty printing for JSONText types.
func (j JSONText) String() string {
	return string(j)
}

// NullJSONText represents a JSONText that may be null.
// NullJSONText implements the scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullJSONText struct {
	JSONText
	Valid bool // Valid is true if JSONText is not NULL
}

// Scan implements the Scanner interface.
func (n *NullJSONText) Scan(value interface{}) error {
	if value == nil {
		n.JSONText, n.Valid = emptyJSON, false
		return nil
	}
	n.Valid = true
	return n.JSONText.Scan(value)
}

// Value implements the driver Valuer interface.
func (n NullJSONText) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.JSONText.Value()
}

func SqlPrepare(x interface{}) interface{} {
	t := reflect.ValueOf(x).Type().String()

	val := reflect.ValueOf(x)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var innerValue reflect.Value

	switch t {
	case "*int64":
		return x
	case "*nullable.NullString":
		f := val.Interface().(NullString)
		innerValue = reflect.ValueOf(f.String)
		return innerValue.String()
	case "*nullable.NullBoolString":
		f := val.Interface().(NullBoolString)
		innerValue = reflect.ValueOf(f.String)
		return innerValue.String()
	case "*nullable.NullTime":
		f := val.Interface().(NullTime)
		innerValue = reflect.ValueOf(f.Time)
		return f.Time
	case "*nullable.ZeroOmitInt64":
		f := val.Interface().(ZeroOmitInt64)
		innerValue = reflect.ValueOf(f.Int64)
		return innerValue.Int()
	case "*nullable.NullInt64":
		f := val.Interface().(NullInt64)
		innerValue = reflect.ValueOf(f.Int64)
		return innerValue.Int()
	case "*nullable.NullJSONText":
		f := val.Interface().(NullJSONText)
		innerValue = reflect.ValueOf(f.JSONText)
		return innerValue.Bytes()
	}
	return "" // TODO review return
}
