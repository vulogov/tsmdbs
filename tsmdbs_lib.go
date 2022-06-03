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

func stddev(src []float64) float64 {
  return stat.StdDev(src, nil)
}

func stderr(src []float64) float64 {
  d := stat.StdDev(src, nil)
  n := float64(len(src))
  return stat.StdErr(d, n)
}

func variance(src []float64) float64 {
  return stat.Variance(src, nil)
}

func skew(src []float64) float64 {
  return stat.Skew(src, nil)
}

func mode(src []float64) float64 {
  m, _ := stat.Mode(src, nil)
  return m
}

func harmonicmean(src []float64) float64 {
  return stat.HarmonicMean(src, nil)
}

func geometricmean(src []float64) float64 {
  return stat.GeometricMean(src, nil)
}

func entropy(src []float64) float64 {
  return stat.Entropy(src)
}

func (ts *TSMDBS) ts_lib_config() error {
  ts.qctx["now"] = time.Now
  ts.qctx["query"] = Query
  ts.qctx["metric"] = Insert_Metric
  ts.qctx["start"] = Start
  ts.qctx["end"] = End
  ts.qctx["float"] = to_float
  ts.qctx["mean"] = mean
  ts.qctx["stddev"] = stddev
  ts.qctx["stderr"] = stderr
  ts.qctx["variance"] = variance
  ts.qctx["skew"] = skew
  ts.qctx["mode"] = mode
  ts.qctx["harmonicmean"] = harmonicmean
  ts.qctx["geometricmean"] = geometricmean
  ts.qctx["entropy"] = entropy
  return nil
}
