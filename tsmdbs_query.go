package tsmdbs

import (
  "time"
  "context"
  "errors"
  "github.com/jinzhu/now"
  "github.com/Jeffail/gabs/v2"
  "github.com/PaesslerAG/gval"
)

type tsmdbsQueryMetrics struct {
  started bool
  ts      *TSMDBS
  host    int64
  _host   string
  key     int64
  _key    string
  isinterval bool
  tsstart time.Time
  tsend   time.Time
}

func (q *tsmdbsQueryMetrics) SelectGVal(ctx context.Context, key string) (interface{}, error) {
  var out *tsmdbsQueryMetrics
  var err error

  if q.started {
    out = q
  } else {
    out = new(tsmdbsQueryMetrics)
    out.ts = q.ts
    out.started = true
    out.isinterval = false
  }

  if ! out.isinterval && key == "time" {
    return out.Time, nil
  } else if ! out.isinterval && key == "all" {
    return out.All, nil
  } else if ! out.isinterval && key == "query" {
    return out.Query, nil
  } else if ! out.isinterval && key == "insert" {
    return out.Insert, nil
  }

  if out.host == 0 {
    out.host, err = q.ts.Host(key)
    out._host = key
  } else if out.key == 0 {
    out.key, err = q.ts.Key(key)
    out._key = key
  }

  if err != nil {
    return nil, err
  }
  return out, nil
}

func (q *tsmdbsQueryMetrics) Time(t time.Time) (*tsmdbsQueryMetrics, error) {
  q.isinterval = true
  q.tsstart    = t
  q.tsend      = time.Now()
  return q, nil
}

func (q *tsmdbsQueryMetrics) All() (*tsmdbsQueryMetrics, error) {
  q.isinterval = true
  start, err := now.Parse("1/1/1970")
  if err != nil {
    return nil, err
  }
  q.tsstart    = start
  q.tsend      = time.Now()
  return q, nil
}

func (q *tsmdbsQueryMetrics) Query() (interface{}, error) {
  q, err := q.All()
  if err != nil {
    return nil, err
  }
  return Query(q)
}

func (q *tsmdbsQueryMetrics) Insert(value interface{}) (interface{}, error) {
  return Insert(q, value)
}

func Insert(q *tsmdbsQueryMetrics, value interface{}) (interface{}, error) {
  if ! q.isinterval {
    q.tsstart = time.Now()
    q.tsend      = time.Now()
    q.isinterval = true
  }
  if q.host == 0 || q.key == 0 {
    return nil, errors.New("Context for insert not present")
  }
  out, err := q.ts.Store(nil, q.tsstart, q._host, q._key, value, []string{}, map[string]interface{}{})
  if err != nil {
    return nil, err
  }
  return out, nil
}

func Query(q *tsmdbsQueryMetrics) ([]interface{}, error) {
  var data []byte

  if q.host == 0 || q.key == 0 || ! q.started {
    return nil, errors.New("Context for query is not defined")
  }
  out := make([]interface{}, 0)
  if ! q.isinterval {
    rows, err := q.ts.db.Query("SELECT VALUE FROM DATA WHERE HOST=? AND KEY=?", q.host, q.key)
    if err != nil {
      return nil, err
    }
    for rows.Next() {
      err = rows.Scan(&data)
      if err != nil {
        return nil, err
      }
      jdata, err := gabs.ParseJSON(data)
      if err != nil {
        return nil, err
      }
      v := jdata.Search("value").Data()
      if v != nil {
        out = append(out, v)
      }
    }
    rows.Close()
  } else {
    interval, err := q.ts.Range(q.tsstart, q.tsend)
    if err != nil {
      return nil, err
    }
    for _, i := range interval {
      rows, err := q.ts.db.Query("SELECT VALUE FROM DATA WHERE HOST=? AND KEY=? AND TS=?", q.host, q.key, i)
      if err != nil {
        return nil, err
      }
      for rows.Next() {
        err = rows.Scan(&data)
        if err != nil {
          return nil, err
        }
        jdata, err := gabs.ParseJSON(data)
        if err != nil {
          return nil, err
        }
        v := jdata.Search("value").Data()
        if v != nil {
          out = append(out, v)
        }
      }
      rows.Close()
    }
  }
  return out, nil
}

func Start(q *tsmdbsQueryMetrics, t time.Time) (interface{}, error) {
  q.tsstart = t
  q.isinterval = true
  return q, nil
}

func End(q *tsmdbsQueryMetrics, t time.Time) (interface{}, error) {
  q.tsend = t
  q.isinterval = true
  return q, nil
}

func (ts *TSMDBS) Query(query string) (interface{}, error) {
  var err error

  lang := gval.Full()
  ts.qctx["db"] = &tsmdbsQueryMetrics{ts: ts, started: false, host: 0, key: 0}
	value, err := lang.Evaluate(
    query,
    ts.qctx,
  )
  if err != nil {
    return nil, err
  }
  return value, nil
}
