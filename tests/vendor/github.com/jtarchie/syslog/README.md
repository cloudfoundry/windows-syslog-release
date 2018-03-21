# syslog

## Usage

The primary usage the library is with the golang package. The Ruby library was only a proof of concept.


### Golang

```sh
go get github.com/jtarchie/syslog
```

```go

import "github.com/jtarchie/syslog/pkg/log"

// ... some code far away ...
log, err := syslog.Parse(line)
if err != nil {
  println("Error", err)
}
println("Message", string(log.Message())) 
```
