package main

import (
	"log"
	"net/http"
	"os"

	"github.com/zetrus/signature-service/internal/api"
	"github.com/zetrus/signature-service/internal/repository"
	"github.com/zetrus/signature-service/internal/service"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}
	repo := repository.NewMemory()
	svc := service.New(repo)
	h := api.NewHandler(svc)
	mux := http.NewServeMux()
	h.Register(mux)
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
