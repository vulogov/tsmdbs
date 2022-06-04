
[![Go](https://github.com/vulogov/tsmdbs/actions/workflows/go.yml/badge.svg)](https://github.com/vulogov/tsmdbs/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/vulogov/tsmdbs/branch/main/graph/badge.svg?token=KGM8LR6KIQ)](https://codecov.io/gh/vulogov/tsmdbs)

# TSMDBS, just another time-series database.

TSMDBS is a embeddable time-series database module designed for storage and analysis of telemetry data. TSMDBS is a package for Golang and the database intended to be embedded in Golang application. There is no server or external app that provides access to TSMDBS functionality.

## What TSMDBS is ...

TSMDBS is a database for storing and querying metrics. It provides Golang API and query language that you can use through Golang API. TSMDBS designed to be ether in-momery or on-disk. TSMDBS uses SQLITE as a database backend, so if you are using it as "on-the-disk" database, you can use any SQLITE tools to access the data. TSMDBS is good enough to serve as a temporary storage for telemetry, that you are going to further process or analyze.

## What TSMDBS is not ...

Due to the fact that TSMDBS primarily designed to be used as "in-memory" temporary storage for the telemetry data, TSMDBS is not replacement for the production, long term, telemetry storage, especially in large quantities. TSMDBS not designed to provide per-second/millisecond/nanosecond timestamp accuracy. While internal timestamp representation is in milliseconds, all telemetry data allocated in "per-minute" buckets. TSMDBS uses JSON for internal telemetry representation. So, there are benefits and disadvantages in that design choice.

## Why shall you give TSMDBS a try ?

If you are creating an instrumentations that are working with time-based data and particularly telemetry data, you are aware that in order to do any meaningful manipulations with telemetry, you have to query multiple samples and store it somewhere. Storing and retrieving of time-series data could become a chore. TSMDBS gives you a bit of assistance, providing you following features:

- Simple, clean, thread-safe API for Golang.
- Ability to store metrics of different data types.
- Ability to add attributes to each metric, making it more meaningful.
- Ability to add "relations" to telemetry. You can store your metric with attached labels and then retrieve you metric using thouse lables.
- You can use ether API calls or TSMDBS Query language (TSQL) to store and retrieve the data.
- TSQL gives you rich logic, arithmetic, query and statistical functions capabilities. In one line of TSQL code you can sample the telemetry and perform statistical computation with it.

## How to use TSMDBS

TSMDBS is a embeddable database for Golang applications and not designed to be used outside of custom Golang applications.

### Installation

```Golang
go get github.com/vulogov/tsmdbs
```

### Use TSMDBS Golang API

#### How to create in-memory TSMDBS database ?

```Golang
import "github.com/vulogov/tsmdbs"

db, err := tsmdbs.InMemory()
```

#### How to create on-disk NSMDBS database ?

```Golang
import "github.com/vulogov/tsmdbs"

db, err := tsmdbs.TS("{file name}")

```

#### How to store telemetry in TSMDBS ?

```Golang
import "os"
import "fmt"
import "github.com/vulogov/tsmdbs"

db, err := tsmdbs.InMemory()
if err != nil {
  fmt.Println(err)
  os.Exit(1)
}
id, err := db.Store(mtype, time.Now(), "hostname", "metric.key", theValue, []string{}, map[string]interface{})
```

Parameters for the Store() functions are:

- Metric type. The string. Default is "metric".
- Timestamp for the metric. You can pass time.Time, number of milliseconds or string with common string representation of date and time. If you pass nil, time.Now() will be used.
- Hostname as string.
- Telemetry key, as string.
- Some telemetry value
- List of the strings with labels for that telemetry.
- Map with attributes for that telemetry.

#### How to close TSMDBS instance ?

```Golang
import "os"
import "fmt"
import "github.com/vulogov/tsmdbs"

db, err := tsmdbs.InMemory()
if err != nil {
  fmt.Println(err)
  os.Exit(1)
}
db.Close()
```

### How to use TSQL

#### Golang API function for serving a TSQL queries

```Golang
import "os"
import "fmt"
import "github.com/vulogov/tsmdbs"

db, err := tsmdbs.InMemory()
if err != nil {
  fmt.Println(err)
  os.Exit(1)
}
result, err := db.Query("{TSQL query}")
if err != nil {
  fmt.Println(err)
  os.Exit(1)
}
```

#### Basic format and functions of TSQL

###### Numbers. For example

```
42
```
##### Strings. Like thouse

```
"Hello world"
```

##### Arithmetic operations: +, -, *, /

```
40 + 2
( 3 + 2 ) * 10
```

##### Logical expressions: <, >, <=, >=, !=, ==, ||, &&

```
1 < 2
```
##### Date declaration

```
date("01/01/1970")
```
##### Telemetry query

```
db.{hostname}.{telemetry key}(.function())
```

Functions are optional and here is the list of functions.

| Function name | Parameters | Description |
|---------------|-----------|-----------|
| time()        | time.Time datetime | Define begin and end intervals as begin is passed as parameter and end is time.Now() |
| all() | None | Set start and end of the query interval as "begin of times" till "now" |
| query() | None | Perform telemetry query for all elements of the specified metric. Returns list of interface{} items |
| sample() | None | Works like a query(), but returns the list of float64 values. |
| insert() | Value | Store telemetry in TSDB |

##### Telemetry-related functions

| Function name | Parameters | Description |
|---------------|-----------|-----------|
| query() | Query context returned by db. construct | Perform query for the context |
| start() | Query context, start timestamp | Set the start timestamp for the context |
| end() | Query context, end timestamp | Set the end timestamp for the context |

##### How to define relations for insert() command ?

```
db.{hostname}.{telemetry key}.{label}.{label}(.optional function())
```

##### Relation query

```
labels.{label}.{label}.query() or sample()
```

#### TSQL examples

##### Storing metrics

```
db.testhost.answer.application1.application2.insert(42)
```

In this example, we are storing telemetry for the host = "testhost", telemetry key = "answer", and value = 42, while attaching this metric to labels "application1" and "application2"

##### Querying all metrics

```
db.testhost.answer.sample()
```
Sampling telemetry stored for host "testhost" and telemetry key "answer"

##### Time-based query of the metrics

```
query(
    end(
      start(
        db.testhost.testkey, date("2022-01-01 23:59:59")
      ),
    now())
)
```
Querying host="testhost" and telemetry key = "testkey", for start date set to "2022-01-01 23:59:59" and end date to now()

#### Useful embedded functions

| Function name | Parameters | Description |
|---------------|-----------|-----------|
| mean() | []float64 | Statistical Mean |
| stddev() | []float64 | Standard deviation |
| stderr() | []float64 | Standard error |
| variance() | []float64 | Statistical variance |
| skew() | []float64 | Skewnes of the sample |
| mode() | []float64 | Mode of the sample |
| harmonicmean() | []float64 | Harmonic Mean computation |
| geometricmean() | []float64 | Geometric Mean computation |
| entropy() | []float64 | Statistical entropy of the sample |

Example:

```
stddev(db.testhost.testkey.sample())
```

Query telemetry from host = "testhost" with telemetry key "testkey", then sampling it and then compute standard deviation for the sample.


### Show me the code !

Every bit of the TSMDBS source code could be located on public GitHub repository: [https://github.com/vulogov/tsmdbs](https://github.com/vulogov/tsmdbs)
