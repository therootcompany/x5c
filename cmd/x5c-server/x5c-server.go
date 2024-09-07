package main

import (
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/therootcompany/x5c"
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

func main() {
	var webRoot string
	var port string

	flag.StringVar(&webRoot, "web-root", "", "Optional root directory for serving files")
	flag.StringVar(&port, "port", "8080", "Port to run the server on")
	flag.Parse()

	http.HandleFunc("GET /api/x509", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		certParam := r.URL.Query().Get("cert")
		if r.URL.Path == "/" && certParam != "" {
			redirectToHash(w, r)
			return
		}

		if webRoot != "" {
			filePath := filepath.Join(webRoot, r.URL.Path)
			if _, err := os.Stat(filePath); err == nil {
				http.ServeFile(w, r, filePath)
				return
			}
		}

		fs := http.FS(static.FS)
		http.FileServer(fs).ServeHTTP(w, r)
	})

	fmt.Printf("Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
