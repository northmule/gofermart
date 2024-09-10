package middlewarehandler

import (
	"github.com/northmule/gophermart/internal/app/services/compressor"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"net/http"
	"strings"
)

// Ожидаемые типы
var expectedContentTypes = map[string]bool{
	"application/json":   true,
	"text/html":          true,
	"application/x-gzip": true,
}

func MiddlewareGzipCompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		requestContentType := request.Header.Get("Content-type")
		if _, ok := expectedContentTypes[requestContentType]; !ok {
			logger.LogSugar.Infof("Ответ сервера не сжимается, для Content-type не поддерживается сжатие")
			next.ServeHTTP(response, request)
			return
		}

		modifiedResponse := response
		// Клиент поддерживает gzip, сжимаем данные для него
		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			gzipWriter := compressor.NewGzipWriter(response)
			modifiedResponse = gzipWriter
			defer gzipWriter.Close()
		}

		// Пришли сжатые данные, надо распаковать
		if strings.Contains(request.Header.Get("Content-Encoding"), "gzip") {
			gzipReader, err := compressor.NewGzipReader(request.Body)
			if err != nil {
				logger.LogSugar.Error("Ошибка compressor.NewGzipReader", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			request.Body = gzipReader
			defer gzipReader.Close()
		}

		next.ServeHTTP(modifiedResponse, request)
	})
}
