# rsv, a very simple CSV marshaller/unmarshaller

## Background
You ever have a CSV file that you need to decode into a struct? Its really annoying to do. 
This library will hopefully make it less annoying, as it aims to plug into existing CSV parsing pipelines, but hopefully
reduce the amount of boilerplate required to write. This means that most of the heavy lifting is left to the encoding/csv package.

This library *currently* supports Unmarshalling and Marshalling via an `idx` struct tag, which specifies what index in the
row this struct tag should be marshalled and unmarshalled into.

## Usage

### Reading a CSV
```go
type Employee struct {
    Id        string `idx:"0"`
    FirstName string `idx:"1"`
    LastName  string `idx:"2"`
    Salary    int64  `idx:"3"` 
}
....
row := []string{"E1234", "John", "Smith", "40000"}
employee := Employee{}
err := rsv.UnmarshalRow(row, &employee)
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

