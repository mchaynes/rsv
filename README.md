# rsv, a very simple CSV marshaller/unmarshaller

## Background
You ever have a CSV file that you need to decode into a struct? Its really annoying to do. rsv isn't trying to re-invent the
wheel, it is meant to to be used with the `encoding/csv` package in the stdlib. Instead of manually building structs
from the output of `csv.Read()`, this library allows you to specify the index of the row via a struct tag `idx`. This clearly
is a fairly brittle way of dealing with CSVs, however in practice CSV files tend to be very static, so this is frequently the simplest way to read in a CSV file.

rsv also doesn't require header columns to be present. Very large data sets tend to be split into many smaller CSV files, and
frequently don't have headers in the split files. rsv is designed with this in mind. 

### Parsing
rsv deals with the parsing the `[]string` into all number types, such as `float64`, `int64`, `uint`, etc. 
If any of the values in the row fail to parse into their respective types, an `ErrFailedToParse` will be returned


## Usage

### Reading a CSV
```go
package main
type Employee struct {
    Id        string `idx:"0"`
    FirstName string `idx:"1"`
    LastName  string `idx:"2"`
    Salary    *int64  `idx:"3,omitempty"` 
}
....
row := []string{"E1234", "John", "Smith", "40000"}
employee := Employee{}
rsv.UnmarshalRow(row, &employee)
row := []string{"E1234", "Jane", "Smit", ""}
// third index is empty (""), therefore salary will be `nil` since omitempty was specified
rsv.UnmarshalRow(row, &employee)
```


Full Example:
```go
package main
import (
    "encoding/csv"
    "io"
    "os"

    "github.com/mchaynes/rsv"
)

type Employee struct {
    Id        string `idx:"0"`
    FirstName string `idx:"1"`
    LastName  string `idx:"2"`
    Salary    int64  `idx:"3"` 
}

func main(){
    f, err := os.Open("employee.csv")
    if err != nil {
        panic(err)
    }
    r := csv.NewReader(f)
    for {
        row, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }
        employee := Employee{}
        err = rsv.UnmarshalRow(row, &employee)
        if err != nil {
            panic(err)
        }
        // do some work
    }
}
```

### Writing to a CSV
```go
employee := Employee{
    Id: "E1234",
    FirstName: "John",
    LastName: "Smith",
    Salary: 40000,
}
row, err := rsv.MarshalRow(employee)
fmt.Println(row)
// output: []interface{"E1234", "John", "Smith", 40000}
```

## Comparisons

| Library                           | Description                                                           |
|---------------------------------- |-----------------------------------------------------------------------|
| https://github.com/gocarina/gocsv | Very good csv parsing library, but requires header row to unmarshal   |

