# HTTP-Forward
Forward http request

# Usage
main.go

```
package main

import (
	forward "github.com/Lyciums/HTTP-Forward"
)

func main() {
	forward.StartForwardService("6666")
}

```

# Build
```
go build main.go
```

# Deployd
```
nohup main > forward.log 2>&1 &
```
