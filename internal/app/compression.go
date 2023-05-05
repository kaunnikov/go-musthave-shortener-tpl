package app

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

var successCompressionContentType = [2]string{"application/json", "text/html"}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CustomCompression(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isGoodContentType := false
		contentType := r.Header.Get("Content-Type")

		// Проверям, ожидает ли клиент, что сервер будет сжимать данные gzip
		isNeedCompression := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		// Проверяем тип контента
		for _, c := range successCompressionContentType {
			if contentType == c {
				isGoodContentType = true
			}
		}

		// Если условия для сжатия не выполнены - отдаём ответ
		if !isNeedCompression || !isGoodContentType {
			h.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			Sugar.Errorf("Error gzip compression: %s", err)
			return
		}
		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				Sugar.Errorf("Error gz close: %s", err)
			}
		}(gz)

		w.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
