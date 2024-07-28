package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/chadsmith12/minhttp"
)

func main() {
    directory := flag.String("directory", "/tmp", "The directory to retrieve the static files from")
    flag.Parse()

    minhttp.MapGet("/", func(hr minhttp.HttpRequest, w minhttp.ResponseWriter) error {
	minhttp.WriteOk(w)
	return nil
    })

    minhttp.MapGet("/echo/{str}", func(hr minhttp.HttpRequest, w minhttp.ResponseWriter) error {
	echo := hr.Params["str"]
	minhttp.WriteText(w, echo)	
	return nil
    })

    minhttp.MapGet("/user-agent", func(hr minhttp.HttpRequest, w minhttp.ResponseWriter) error {
	minhttp.WriteText(w, hr.Headers.UserAgent())

	return nil
    })

    minhttp.MapGet("/files/{fileName}", func(hr minhttp.HttpRequest, w minhttp.ResponseWriter) error {
	filePath := path.Join(*directory, hr.Params["fileName"])
	text, err := os.ReadFile(filePath)
	if err != nil {
	    minhttp.WriteNotFound(w)
	    return err
	}
	
	minhttp.WriteOctetStream(w, text)
	return nil
    })

    minhttp.MapPost("/files/{fileName}", func(hr minhttp.HttpRequest, w minhttp.ResponseWriter) error {
	bodyBytes, err := io.ReadAll(hr.Body)
	if err != nil {
	    minhttp.WriteInernalServerError(w, err.Error())
	    return err
	}
	
	filePath := path.Join(*directory, hr.Params["fileName"])
	err = os.WriteFile(filePath, bodyBytes, 0666)
	if err != nil {
	    minhttp.WriteInernalServerError(w, err.Error())
	    return err
	}
	
	minhttp.WriteCreated(w)
	return nil
    })
    err := minhttp.ListenAndServe("0.0.0.0:4221")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
