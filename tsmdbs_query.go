package tsmdbs

import (
  "time"
  "context"
  "errors"
  "github.com/Jeffail/gabs/v2"
  "github.com/PaesslerAG/gval"
)

type tsmdbsQueryMetrics struct {
  started bool
  ts      *TSMDBS
  host    int64
  key     int64
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

  if out.host == 0 {
    out.host, err = q.ts.Host(key)
  } else if out.key == 0 {
    out.key, err = q.ts.Key(key)
  }

  if err != nil {
    return nil, err
  }
  return out, nil
}

func Query(q *tsmdbsQueryMetrics) (interface{}, error) {
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
  ts.qctx["now"] = time.Now
  ts.qctx["query"] = Query
  ts.qctx["start"] = Start
  ts.qctx["end"] = End
	value, err := lang.Evaluate(
    query,
    ts.qctx,
  )
  if err != nil {
    return nil, err
  }
  return value, nil
}
