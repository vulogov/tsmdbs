package tsmdbs

import (
  "fmt"
  "time"
  "testing"
  "github.com/jinzhu/now"
  "github.com/stretchr/testify/assert"
)

func TestTSNSMDBS1(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Now()
  assert.Equal(t, err, nil)
  fmt.Println("TS index", res)
  db.Close()
}

func TestTSNSMDBS2(t *testing.T) {
  db, err := TS(DBNAME_ON_DISK)
  start, _ := now.Parse("01/31/2022")
  assert.Equal(t, err, nil)
  res, err := db.Now()
  res, err = db.Time(start)
  assert.Equal(t, err, nil)
  fmt.Println("TS index", res)
  out, err := db.Range(start, time.Now())
  assert.Equal(t, err, nil)
  assert.Equal(t, len(out), 2)
  db.Close()
}
