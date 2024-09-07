package main

import (
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/therootcompany/x5c"
	"github.com/therootcompany/x5c/internal"
	"github.com/therootcompany/x5c/static"
)

type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(JSONError{
		Error: message,
		Code:  statusCode,
	})
}

func redirectToHash(w http.ResponseWriter, r *http.Request) {
	certParam := r.URL.Query().Get("cert")
	query := url.Values{"cert": {certParam}}
	search := query.Encode()
	redirectURL := "/#?" + search
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func decodeX509(w http.ResponseWriter, r *http.Request) {
	certParam := r.URL.Query().Get("cert")
	if certParam == "" {
		msg := "missing '?cert=<hex|base64|pem>"
		jsonError(w, msg, http.StatusBadRequest)
		return
	}

	certBytes, err := x5c.MagicDecodeCertString(certParam)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		msg := fmt.Sprintf("successfully parsed string as bytes, but failed to decode certificate: %v", err)
		jsonError(w, msg, http.StatusBadRequest)
		return
	}

	info := x5c.Summarize(cert)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func main() {
	var port string
	overlayFS := &internal.OverlayFS{}
	var rateLimit int

	flag.StringVar(&overlayFS.WebRoot, "web-root", "./public/", "serve from the given directory")
	flag.BoolVar(&overlayFS.WebRootOnly, "web-root-only", false, "do not serve the embedded web root")
	flag.IntVar(&rateLimit, "rate-limit", 50, "number of requests per IP per 5 minutes (0 for unlimited)")
	flag.StringVar(&port, "port", "8080", "bind and listen for http on this port")
	flag.Parse()

	overlayFS.LocalFS = http.Dir(overlayFS.WebRoot)
	overlayFS.EmbedFS = http.FS(static.FS)

	decoderHandler := decodeX509
	if rateLimit > 0 {
		decoderHandler = internal.RateLimitMiddleware(decoderHandler)
	}
	http.HandleFunc("GET /api/x509", decoderHandler)

	fileServer := http.FileServer(overlayFS)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		certParam := r.URL.Query().Get("cert")
		if r.URL.Path == "/" && certParam != "" {
			redirectToHash(w, r)
			return
		}

		fileServer.ServeHTTP(w, r)
	})

	fmt.Printf("Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
