## pterodactyl-go

This is an implementation of the Pterodactyl project's API in Go

### How to test?

Follow the instructions found at https://pterodactyl.io/panel/1.0/getting_started.html to setup Pterodactyl on your server, local machine or docker containers. Once it is running and you're able to login as an admin you can generate an API key to be used. After that the initialization looks like this

```go
package main

import (
	pterodactyl "github.com/tehzwen/pterodactyl-go/pkg"
)

var (
	API_KEY         = ""
	PTERO_HOST_ADDR = "" // ie: http://localhost
)

func main() {
    ptero, err := pterodactyl.NewApplicationClient(PTERO_HOST_ADDR, pterodactyl.WithApiKey(API_KEY))
    if err != nil {
        fatal(err)
    }

    // run whatever api operations you'd like from here
}
```

