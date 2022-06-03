package tsmdbs

import (
  "time"
  "github.com/jinzhu/now"
)

func (ts *TSMDBS) Now() (int64, error) {
  stamp := time.Now()
  return ts.Time(stamp)
}

func (ts *TSMDBS) Time(stamp time.Time) (int64, error) {
  var id int64
  var found bool

  err := ts.db.Ping()
  if err != nil {
    return 0, err
  }
  found = false
  ts1 := now.With(stamp).BeginningOfMinute().UnixNano() / int64(time.Millisecond)
  ts2 := now.With(stamp).EndOfMinute().UnixNano() / int64(time.Millisecond)
  rows, err := ts.db.Query("SELECT ID FROM TS WHERE START=? AND END=?", ts1, ts2)
  if err != nil {
    return 0, err
  }
  for rows.Next() {
    err = rows.Scan(&id)
    found = true
  }
  rows.Close()
  if found {
    return id, nil
  } else {
    tx, err := ts.db.Begin()
    if err != nil {
      return 0, err
    }
    tx.Exec("INSERT INTO TS(START, END) VALUES(?, ?)", ts1, ts2)
    err = tx.Commit()
    if err != nil {
      return 0, err
    }
  }
  return ts.Time(stamp)
}

func (ts *TSMDBS) Range(start time.Time, end time.Time) ([]int64, error) {
  var id int64
  out := make([]int64, 0)
  ts1 := now.With(start).BeginningOfMinute().UnixNano() / int64(time.Millisecond)
  ts2 := now.With(end).EndOfMinute().UnixNano() / int64(time.Millisecond)
  rows, err := ts.db.Query("SELECT ID FROM TS WHERE START>=? AND END<=?", ts1, ts2)
  if err != nil {
    return nil, err
  }
  for rows.Next() {
    err = rows.Scan(&id)
    out = append(out, id)
  }
  rows.Close()
  return out, nil
}
