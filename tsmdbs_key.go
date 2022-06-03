package tsmdbs

import (

)


func (ts *TSMDBS) Key(key string) (int64, error) {
  var id int64
  var found bool

  err := ts.db.Ping()
  if err != nil {
    return 0, err
  }
  rows, err := ts.db.Query("SELECT ID FROM KEY WHERE NAME=?", key)
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
    tx.Exec("INSERT INTO KEY(NAME) VALUES(?)", key)
    err = tx.Commit()
    if err != nil {
      return 0, err
    }
  }
  return ts.Key(key)
}
