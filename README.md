# Simple Unique ID Generator

Use 6 byte time.UnixNano() for submission time. Just use the XOR operation of the counter value

## Install

    go get github.com/itrabbit/rid

## Usage

```go
id := rid.New()

println(id.String())
// -> krt54gkt2cakckbs0lm0
```

Get embedded info:

```go
id.Mid() // Hardware Address CRC ID
id.Pid() // Process Pid
id.Time()
id.Counter() 
```

