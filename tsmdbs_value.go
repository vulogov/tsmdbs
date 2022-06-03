package tsmdbs

import (
  "time"
  "errors"
  "github.com/jinzhu/now"
  "github.com/Jeffail/gabs/v2"
  "github.com/aidarkhanov/nanoid/v2"
)

func (ts *TSMDBS) Store(mtype interface{}, date interface{}, host string, key string, value interface{}, rel []string, attr map[string]interface{}) (string, error) {
  data := gabs.New()
  switch mtype.(type) {
  case string:
    data.Set(mtype.(string), "type")
  default:
    data.Set("metric", "type")
  }
  switch date.(type) {
  case string:
    data.Set(date.(string), "timestamp")
  case float64:
    _t := time.Unix(0, int64(date.(float64)) * int64(time.Millisecond))
    data.Set(_t, "timestamp")
  case int64:
    _t := time.Unix(0, date.(int64) * int64(time.Millisecond))
    data.Set(_t, "timestamp")
  case time.Time:
    data.Set(date.(time.Time), "timestamp")
  }
  data.Set(host, "host")
  data.Set(key, "key")
  data.Set(value, "value")
  data.Array("relation")
  if rel != nil {
    for _, r := range rel {
      data.ArrayAppend(r, "relation")
    }
  }
  if attr != nil {
    for k, v := range attr {
      data.Set(v, "attributes", k)
    }
  }
  return ts.Value(data)
}

func (ts *TSMDBS) Json(jdata []byte) (string, error) {
  data, err := gabs.ParseJSON(jdata)
  if err != nil {
    return "", err
  }
  return ts.Value(data)
}

func (ts *TSMDBS) Value(data *gabs.Container) (string, error) {
  var id string
  var stamp int64

  err := ts.db.Ping()
  if err != nil {
    return "", err
  }
  host := data.Search("host").Data()
  if host == nil {
    return "", errors.New("Field 'host' is missing in JSON")
  }
  key := data.Search("key").Data()
  if key == nil {
    return "", errors.New("Field 'key' is missing in JSON")
  }
  value := data.Search("value").Data()
  if value == nil {
    return "", errors.New("Field 'value' is missing in JSON")
  }
  timestamp := data.Search("timestamp").Data()
  if timestamp == nil {
    stamp, err = ts.Now()
    if err != nil {
      return "", err
    }
  } else {
    switch timestamp.(type) {
    case string:
      ctime, err := now.Parse(timestamp.(string))
      if err != nil {
        return "", err
      }
      stamp, err  = ts.Time(ctime)
      if err != nil {
        return "", err
      }
    case float64:
      _t := time.Unix(0, int64(timestamp.(float64)) * int64(time.Millisecond))
      stamp, err  = ts.Time(_t)
      if err != nil {
        return "", err
      }
    case int64:
      _t := time.Unix(0, timestamp.(int64) * int64(time.Millisecond))
      stamp, err  = ts.Time(_t)
      if err != nil {
        return "", err
      }
    case time.Time:
      stamp, err  = ts.Time(timestamp.(time.Time))
      if err != nil {
        return "", err
      }
    default:
      return "", errors.New("timestamp value not of the supported types")
    }
  }
  host_id, err := ts.Host(host.(string))
  if err != nil {
    return "", err
  }
  key_id, err  := ts.Key(key.(string))
  if err != nil {
    return "", err
  }
  id, err = nanoid.New()
  if err != nil {
    return "", err
  }
  tx, err := ts.db.Begin()
  if err != nil {
    return "", err
  }
  tx.Exec("INSERT INTO DATA(ID, TS, HOST, KEY, VALUE) VALUES(?, ?, ?, ?, ?)", id, stamp, host_id, key_id, data.Bytes())
  err = tx.Commit()
  if err != nil {
    return "", err
  }
  if data.Exists("relation") {
    for _, r := range data.S("relation").Children() {
      rid, err := ts.Relation(r.Data().(string))
      if err != nil {
        return "", err
      }
      tx, err := ts.db.Begin()
      if err != nil {
        return "", err
      }
      tx.Exec("INSERT INTO DREL(REL, DATA) VALUES(?, ?)", rid, id)
      err = tx.Commit()
      if err != nil {
        continue
      }
    }
  }
  return id, nil
}
