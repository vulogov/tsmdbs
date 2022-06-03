package tsmdbs

import (
  "fmt"
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
  rel     []string
}

type tsmdbsLabelMetrics struct {
  started bool
  ts      *TSMDBS
  rel     []string
}

func (q *tsmdbsLabelMetrics) SelectGVal(ctx context.Context, key string) (interface{}, error) {
  var out *tsmdbsLabelMetrics

  if q.started {
    out = q
  } else {
    out = new(tsmdbsLabelMetrics)
    out.ts = q.ts
    out.started = true
  }

  if out.started && key == "query" {
    return out.Query, nil
  } else if out.started && key == "sample" {
    return out.Sample, nil
  }

  out.rel = append(out.rel, key)
  return out, nil
}

func (q *tsmdbsLabelMetrics) Sample() ([]float64, error) {
  res, err := q.Query()
  if err != nil {
    return nil, err
  }
  return to_float(res.([]interface{})), nil
}

func (q *tsmdbsLabelMetrics) Query() (interface{}, error) {
  var data []byte

  if len(q.rel) == 0 {
    return nil, errors.New("relation context not set")
  }
  _query := "select distinct value from data,drel,relation where data.id = drel.data and drel.rel = relation.id and relation.name in %v"
  labels := ""
  is_started := false
  for _, v := range q.rel {
    if ! is_started {
      labels += fmt.Sprintf("( '%v' ", v)
      is_started = true
    } else {
      labels += fmt.Sprintf(", '%v' ", v)
    }
  }
  labels += " )"
  query := fmt.Sprintf(_query, labels)
  out := make([]interface{}, 0)
  rows, err := q.ts.db.Query(query)
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
  return out, nil
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
  } else if ! out.isinterval && key == "metric" {
    return out.Insert, nil
  } else if ! out.isinterval && key == "sample" {
    return out.Sample, nil
  }

  if out.host == 0 {
    out.host, err = q.ts.Host(key)
    out._host = key
  } else if out.key == 0 {
    out.key, err = q.ts.Key(key)
    out._key = key
  } else {
    out.rel = append(out.rel, key)
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

func (q *tsmdbsQueryMetrics) Sample() (interface{}, error) {
  q, err := q.All()
  if err != nil {
    return nil, err
  }
  res, err := Query(q)
  if err != nil {
    return nil, err
  }
  return to_float(res), nil
}

func (q *tsmdbsQueryMetrics) Insert(value interface{}) (interface{}, error) {
  return Insert_Metric(q, value)
}

func Insert_Metric(q *tsmdbsQueryMetrics, value interface{}) (interface{}, error) {
  return Insert(q, "metric", value)
}

func Insert(q *tsmdbsQueryMetrics, mtype string, value interface{}) (interface{}, error) {
  if ! q.isinterval {
    q.tsstart = time.Now()
    q.tsend      = time.Now()
    q.isinterval = true
  }
  if q.host == 0 || q.key == 0 {
    return nil, errors.New("Context for insert not present")
  }
  out, err := q.ts.Store(mtype, q.tsstart, q._host, q._key, value, q.rel, map[string]interface{}{})
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
  ts.qctx["labels"] = &tsmdbsLabelMetrics{ts: ts, started: false}
	value, err := lang.Evaluate(
    query,
    ts.qctx,
  )
  if err != nil {
    return nil, err
  }
  return value, nil
}
