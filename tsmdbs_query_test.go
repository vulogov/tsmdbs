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

func TestQNSMDBS7(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Store(nil, time.Now(), "testhost", "testkey", 42, []string{"abc", "cde"}, map[string]interface{}{"Hello":"world"})
  assert.Equal(t, err, nil)
  fmt.Println("D=", res)
  out, err := db.Query(`float(db.testhost.testkey.query())`)
  assert.Equal(t, err, nil)
  assert.Equal(t, len(out.([]float64)), 1)
  db.Close()
}

func TestQNSMDBS8(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  res, err := db.Query(`db.testhost.testkey.insert(42)`)
  assert.Equal(t, err, nil)
  fmt.Println("D=", res)
  out, err := db.Query(`float(db.testhost.testkey.query())`)
  assert.Equal(t, err, nil)
  assert.Equal(t, len(out.([]float64)), 1)
  db.Close()
}

func TestQNSMDBS9(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.insert(42)`)
  db.Query(`db.testhost.testkey.insert(41)`)
  db.Query(`db.testhost.testkey.insert(40)`)
  db.Query(`db.testhost.testkey.insert(30)`)
  db.Query(`db.testhost.testkey.insert(42)`)
  out, err := db.Query(`db.testhost.testkey.sample()`)
  assert.Equal(t, err, nil)
  assert.Equal(t, len(out.([]float64)), 5)
  db.Close()
}

func TestQNSMDBS10(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.insert(42)`)
  db.Query(`db.testhost.testkey.insert(41)`)
  db.Query(`db.testhost.testkey.insert(40)`)
  db.Query(`db.testhost.testkey.insert(30)`)
  db.Query(`db.testhost.testkey.insert(42)`)
  out, err := db.Query(`stddev(db.testhost.testkey.sample())`)
  assert.Equal(t, err, nil)
  fmt.Println("stddev=",out)
  assert.Equal(t, out, 5.0990195135927845)
  db.Close()
}

func TestQNSMDBS11(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.insert(42)`)
  db.Query(`db.testhost.testkey.insert(41)`)
  db.Query(`db.testhost.testkey.insert(40)`)
  db.Query(`db.testhost.testkey.insert(30)`)
  db.Query(`db.testhost.testkey.insert(42)`)
  out, err := db.Query(`stderr(db.testhost.testkey.sample())`)
  assert.Equal(t, err, nil)
  fmt.Println("stderr=",out)
  assert.Equal(t, out, 2.2803508501982757)
  db.Close()
}

func TestQNSMDBS12(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.insert(42)`)
  db.Query(`db.testhost.testkey.insert(41)`)
  db.Query(`db.testhost.testkey.insert(40)`)
  db.Query(`db.testhost.testkey.insert(30)`)
  db.Query(`db.testhost.testkey.insert(42)`)
  out, err := db.Query(`variance(db.testhost.testkey.sample())`)
  assert.Equal(t, err, nil)
  fmt.Println("variance=",out)
  assert.Equal(t, out, float64(26))
  db.Close()
}

func TestQNSMDBS13(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.insert(42)`)
  db.Query(`db.testhost.testkey.insert(41)`)
  db.Query(`db.testhost.testkey.insert(40)`)
  db.Query(`db.testhost.testkey.insert(30)`)
  db.Query(`db.testhost.testkey.insert(42)`)
  out, err := db.Query(`skew(db.testhost.testkey.sample())`)
  assert.Equal(t, err, nil)
  fmt.Println("skew=",out)
  assert.Equal(t, out, -2.0931625961863882)
  db.Close()
}

func TestQNSMDBS14(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.insert(42)`)
  db.Query(`db.testhost.testkey.insert(41)`)
  db.Query(`db.testhost.testkey.insert(40)`)
  db.Query(`db.testhost.testkey.insert(30)`)
  db.Query(`db.testhost.testkey.insert(42)`)
  out, err := db.Query(`mode(db.testhost.testkey.sample())`)
  assert.Equal(t, err, nil)
  fmt.Println("mode=", out)
  assert.Equal(t, out, float64(42))
  db.Close()
}

func TestQNSMDBS15(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.insert(42)`)
  db.Query(`db.testhost.testkey.insert(41)`)
  db.Query(`db.testhost.testkey.insert(40)`)
  db.Query(`db.testhost.testkey.insert(30)`)
  db.Query(`db.testhost.testkey.insert(42)`)
  out, err := db.Query(`harmonicmean(db.testhost.testkey.sample())`)
  assert.Equal(t, err, nil)
  fmt.Println("harmonicmean=", out)
  assert.Equal(t, out, 38.36043662285587)
  db.Close()
}

func TestQNSMDBS16(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.insert(42)`)
  db.Query(`db.testhost.testkey.insert(41)`)
  db.Query(`db.testhost.testkey.insert(40)`)
  db.Query(`db.testhost.testkey.insert(30)`)
  db.Query(`db.testhost.testkey.insert(42)`)
  out, err := db.Query(`geometricmean(db.testhost.testkey.sample())`)
  assert.Equal(t, err, nil)
  fmt.Println("geometricmean=", out)
  assert.Equal(t, out, 38.698375708929646)
  db.Close()
}

func TestQNSMDBS17(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.insert(42)`)
  db.Query(`db.testhost.testkey.insert(41)`)
  db.Query(`db.testhost.testkey.insert(40)`)
  db.Query(`db.testhost.testkey.insert(30)`)
  db.Query(`db.testhost.testkey.insert(42)`)
  out, err := db.Query(`mode(db.testhost.testkey.sample())`)
  assert.Equal(t, err, nil)
  fmt.Println("mode=", out)
  assert.Equal(t, out, float64(42))
  db.Close()
}

func TestQNSMDBS18(t *testing.T) {
  db, err := TS(DBNAME)
  assert.Equal(t, err, nil)
  db.Query(`db.testhost.testkey.app1.insert(42)`)
  db.Query(`db.testhost.testkey.app1.insert(41)`)
  db.Query(`db.testhost.testkey.app2.insert(40)`)
  db.Query(`db.testhost.testkey.app2.insert(30)`)
  db.Query(`db.testhost.testkey.app1.insert(42)`)
  res, err := db.Query(`labels.app1.app2.sample()`)
  assert.Equal(t, err, nil)
  assert.Equal(t, len(res.([]float64)), 5)
  db.Close()
}
