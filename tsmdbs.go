package tsmdbs

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

type TSMDBS struct {
  db        *sql.DB
  name      string
  qctx      map[string]interface{}
}

func TS(name string) (*TSMDBS, error) {
  var err error

  res := new(TSMDBS)
  res.db, err = sql.Open("sqlite3", name)
  if err != nil {
    return nil, err
  }
  res.name    = name
  res.qctx    = make(map[string]interface{})
  err = res.Recreate()
  res.ts_lib_config()
  return res, err
}

func (ts *TSMDBS) Recreate() error {
  var err error

  _, err = ts.db.Exec(DBSQL)
  if err != nil {
    return err
  }
  return nil
}

func (ts *TSMDBS) Close() {
  ts.db.Close()
}
