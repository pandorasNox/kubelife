package digitalocean

import (
	do "github.com/digitalocean/godo"
)

func test() {
	client := do.NewFromToken("my-digitalocean-api-token")
	_ = client
}
