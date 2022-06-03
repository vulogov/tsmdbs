package tsmdbs

import (
  "fmt"
  "testing"
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
