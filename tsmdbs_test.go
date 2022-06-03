package tsmdbs

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

const DBNAME=":memory:"
const DBNAME_ON_DISK="test.db"


func TestNSMDBS1(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Close()
}
