# Unique ID Generator

Use 6 byte `time.UnixNano()` for submission time

## Install

    go get github.com/itrabbit/rid

## Usage

```go
fmt.Println(rid.New())
// -> krt54gkt2cakckbs0lm0
```

Get embedded info:

```go
id := rid.New()
id.Mid() // Hardware Address CRC ID
id.Pid() // Process Pid
id.Time()
id.Counter() 
```

