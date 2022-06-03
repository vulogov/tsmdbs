package tsmdbs

import (
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestHKNSMDBS1(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Host("testhost")
  assert.Equal(t, err, nil)
  fmt.Println("HOST index", res)
  db.Close()
}

func TestHKNSMDBS2(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Key("testkey")
  assert.Equal(t, err, nil)
  fmt.Println("KEY index", res)
  db.Close()
}

func TestHKNSMDBS3(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Key("answser")
  res, err := db.Key("testkey")
  assert.Equal(t, err, nil)
  fmt.Println("KEY index", res)
  db.Close()
}

func TestHKNSMDBS4(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Relation("therelation")
  assert.Equal(t, err, nil)
  fmt.Println("REL index", res)
  db.Close()
}
