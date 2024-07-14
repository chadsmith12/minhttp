package main

import (
	"fmt"
	"os"

	"github.com/chadsmith12/minhttp"
)

func main() {
    err := minhttp.ListenAndServe("0.0.0.0:4221")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
