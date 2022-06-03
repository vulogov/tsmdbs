package tsmdbs

import (
  "time"
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

const VAL = `
{
  "host": "testhost",
  "key": "answer",
  "value": 42
}
`

const VAL2 = `
{
  "host": "testhost",
  "key": "answer",
  "value": 42,
  "relation": [
    "abc",
    "cde"
  ]
}
`

func TestVNSMDBS1(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Json([]byte(VAL))
  assert.Equal(t, err, nil)
  fmt.Println("VAL index", res)
  db.Close()
}

func TestVNSMDBS2(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Json([]byte(VAL2))
  assert.Equal(t, err, nil)
  fmt.Println("VAL index", res)
  db.Close()
}

func TestVNSMDBS3(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Store(nil, nil, "testhost", "testkey", 42, []string{"abc", "cde"}, map[string]interface{}{"Hello":"world"})
  assert.Equal(t, err, nil)
  fmt.Println("VAL index", res)
  db.Close()
}

func TestVNSMDBS4(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Store(nil, time.Now(), "testhost", "testkey", 42, []string{"abc", "cde"}, map[string]interface{}{"Hello":"world"})
  assert.Equal(t, err, nil)
  fmt.Println("VAL index", res)
  db.Close()
}
