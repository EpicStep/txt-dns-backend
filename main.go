package main

import (
	"context"
	"fmt"
	"github.com/EpicStep/txt-dns-backend/pkg/server"
	dns "github.com/Focinfi/go-dns-resolver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const CloudFlareDNS = "1.1.1.1:53"

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, os.Kill)

	srv := server.New("localhost:8582", Routes())

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("server shutdown failed")
	}
}

func Routes() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("/lookup/txt", http.HandlerFunc(TXTLookup))

	return router
}

func TXTLookup(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	host := params.Get("host")

	if results, err := dns.Exchange(host, CloudFlareDNS, dns.TypeTXT); err == nil {
		if len(results) == 0 {
			fmt.Fprintf(w, "Host has no records")
		} else {
			for _, r := range results {
				fmt.Fprintf(w, "%s %s %s\n", r.Record, r.Type, r.Content)
			}
		}
	} else {
		fmt.Fprintf(w, "TXT lookup failed")
	}
}
