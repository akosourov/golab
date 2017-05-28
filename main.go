package main

import (
	"flag"
	"github.com/akosourov/golab/api"
)

func main() {
	bindAddr := flag.String("bind_addr", ":8080", "Set bind address")
	lruSize := flag.Int("lru_size", 20, "Set size for LRU cache per driver")
	flag.Parse()
	a := api.New(*bindAddr, *lruSize)
	a.Start()
	a.WaitStop()
}
