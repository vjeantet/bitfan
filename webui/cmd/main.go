package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vjeantet/bitfan/ui"
)

func main() {
	addr := "127.0.0.1:8080"
	if port := os.Getenv("PORT"); len(port) > 0 {
		fmt.Printf("Environment variable PORT=\"%s\"", port)
		addr = "127.0.0.1:" + port
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
	if err != nil {
		log.Fatal(err)
	}

	httpServerMux := http.NewServeMux()

	httpServerMux.Handle("/ui/", ui.Handler(dir, "ui"))
	http.ListenAndServe(addr, httpServerMux)
}
