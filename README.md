# receiver [![Build Status](https://secure.travis-ci.org/Gujarats/receiver.png)](http://travis-ci.org/Gujarats/receiver)

Get the request value into given struct with validation

let's say you have incoming request and wanted to parse the `Form` from `http.Request`.
You can do that by using this library.

## Installation

``` go 
go get github.com/Gujarats/receiver
```


## Usage
First of all you need to create the struct for storing the data lets say you have some `struct ` like this : 

```go
// Note : the latitude,longitude,name,distance is the key name to get the value name 
// second argument is to tell wheter the request is required or not
type Sample struct {
	Lat      float64 `request:"latitude,required"`
	Lon      float64 `request:"longitude,required"`
	Name     string  `request:"name,required"`
	Distance int64   `request:"distance,optional"`
}


... some code here

var sample Sample
// r is comes from r *http.Request
err := receiver.SetData(&sample,r)
if err != nil {
    log.Fatal(err)
}

... some code here

```
