package tsmdbs

import (
  "time"
  "gonum.org/v1/gonum/stat"
)

func to_float(src []interface{}) []float64 {
  out := make([]float64, 0)
  for n := 0; n < len(src); n++ {
    switch v := src[n].(type) {
    case float64:
      out = append(out, v)
    case int64:
      out = append(out, float64(v))
    }
  }
  return out
}

func mean(src []float64) float64 {
  return stat.Mean(src, nil)
}

func (ts *TSMDBS) ts_lib_config() error {
  ts.qctx["now"] = time.Now
  ts.qctx["query"] = Query
  ts.qctx["metric"] = Insert_Metric
  ts.qctx["start"] = Start
  ts.qctx["end"] = End
  ts.qctx["float"] = to_float
  ts.qctx["mean"] = mean
  return nil
}
