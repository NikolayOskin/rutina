// +build ignore

package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/neonxp/rutina/v3"
)

func main() {
	// New instance with builtin context
	r := rutina.New(rutina.ListenOsSignals(os.Interrupt, os.Kill))

	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	// Starting http server and listen connections
	r.Go(func(ctx context.Context) error {
		if err := srv.ListenAndServe(); err != nil {
			return err
		}
		log.Println("Server stopped")
		return nil
	})

	// Gracefully stopping server when context canceled
	r.Go(func(ctx context.Context) error {
		<-ctx.Done()
		log.Println("Stopping server...")
		return srv.Shutdown(ctx)
	})

	if err := r.Wait(); err != nil {
		log.Fatal(err)
	}
	log.Println("All routines successfully stopped")
}
