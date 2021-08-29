# HTTP-Forward
Forward http request

# Usage
main.go

```go
package main

import (
    forward "github.com/Lyciums/HTTP-Forward"
)

func main() {
    forward.StartForwardService("6666")
}

```

```shell
go run main.go
```

# Build
```shell
go build main.go
```

# Deployd
```shell
nohup main > forward.log 2>&1 &
```
