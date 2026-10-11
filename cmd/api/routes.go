package main

import (
	"fmt"
	"net/http"
)

func routes() {
	mux := http.NewServeMux()
	fmt.Print(mux)
}
