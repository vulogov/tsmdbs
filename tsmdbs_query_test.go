package tsmdbs

import (
  "fmt"
  "time"
  "testing"
  "github.com/stretchr/testify/assert"
)


func TestQNSMDBS1(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  val, err := db.Query(`42`)
  assert.Equal(t, err, nil)
  assert.Equal(t, val, float64(42))
  db.Close()
}

func TestQNSMDBS2(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  val, err := db.Query(`41 + 1`)
  assert.Equal(t, err, nil)
  assert.Equal(t, val, float64(42))
  db.Close()
}

func TestQNSMDBS3(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  _, err = db.Query(`db.query() `)
  assert.NotEqual(t, err, nil)
  db.Close()
}

func TestQNSMDBS4(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Store(nil, time.Now(), "testhost", "testkey", 42, []string{"abc", "cde"}, map[string]interface{}{"Hello":"world"})
  assert.Equal(t, err, nil)
  fmt.Println("D=", res)
  _, err = db.Query(`query(end(start(db.testhost.testkey, date("2022-01-01 23:59:59")), now()))`)
  assert.Equal(t, err, nil)
  db.Close()
}

func TestQNSMDBS5(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Store(nil, time.Now(), "testhost", "testkey", 42, []string{"abc", "cde"}, map[string]interface{}{"Hello":"world"})
  assert.Equal(t, err, nil)
  fmt.Println("D=", res)
  out, err := db.Query(`float(query(db.testhost.testkey))`)
  assert.Equal(t, err, nil)
  assert.Equal(t, len(out.([]float64)), 1)
  db.Close()
}

func TestQNSMDBS6(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Store(nil, time.Now(), "testhost", "testkey", 42, []string{"abc", "cde"}, map[string]interface{}{"Hello":"world"})
  res, err = db.Store(nil, time.Now(), "testhost", "testkey", 21, []string{"abc", "cde"}, map[string]interface{}{"Hello":"world"})
  assert.Equal(t, err, nil)
  fmt.Println("D=", res)
  out, err := db.Query(`mean(float(query(db.testhost.testkey)))`)
  assert.Equal(t, err, nil)
  assert.Equal(t, out, 31.5)
  db.Close()
}
