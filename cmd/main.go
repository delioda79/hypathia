package main

import (
	"fmt"
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"
	phttp "github.com/beatlabs/patron/sync/http"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	version = "dev"
	name    = "hypatia"
)

func init() {
	err := patron.Setup(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Debugf("no .env file exists: %v", err)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {

	r := phttp.NewRouteRaw("/", "GET", Index, false)

	srv, err := patron.New(
		name,
		version,
		patron.Routes([]phttp.Route{r}),
	)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to run service %v", err)
	}

	return nil
}

func Index(w http.ResponseWriter, req *http.Request) {
	f, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		fmt.Println("error reading file")
	}
	w.Write(f)
}
