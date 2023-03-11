package api

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type contextKey string

const (
	UserIDContextKey contextKey = "user-id"
)

var _ http.ResponseWriter = gzipWriter{}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipResponseHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncodingHeader := r.Header.Get("Accept-Encoding")
		if !strings.Contains(acceptEncodingHeader, "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, err := io.WriteString(w, err.Error())
			if err != nil {
				http.Error(w, "Could not compress response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				log.Printf("error closing writer: %v\n", err)
			}
		}(gz)

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

var _ io.ReadCloser = gzipRequestBody{}

type gzipRequestBody struct {
	ReadCloser io.ReadCloser
}

func (g gzipRequestBody) Read(p []byte) (n int, err error) {
	return g.ReadCloser.Read(p)
}

func (g gzipRequestBody) Close() error {
	return g.ReadCloser.Close()
}

func gzipRequestHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncodingHeader := r.Header.Get("Content-Encoding")
		if !strings.Contains(contentEncodingHeader, "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		log.Println("Encoded request is received")

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(gz *gzip.Reader) {
			err := gz.Close()
			if err != nil {
				log.Printf("error closing reader: %v\n", err)
			}
		}(gz)
		r.Body = gzipRequestBody{ReadCloser: gz}
		next.ServeHTTP(w, r)
	})
}

func checkIfExistsCookieAndGenerateIfNot(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("my-cookie")
		if err != nil {
			cookie = &http.Cookie{Name: "my-cookie", Value: strconv.FormatInt(time.Now().Unix(), 10)}
			http.SetCookie(w, cookie)
		}
		ctx := context.WithValue(r.Context(), UserIDContextKey, cookie.Value)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
