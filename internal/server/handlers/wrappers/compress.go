package wrappers

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CompressWrapper(fn func(cxt *gin.Context)) func(cxt *gin.Context) {
	res := func(ctx *gin.Context) {
		decompressRequest(ctx)

		if !strings.Contains(ctx.GetHeader("Accept-Encoding"), "gzip") {
			fn(ctx)
			return
		}

		ctx.Writer = makeGzipWriter(ctx)
		fn(ctx)
	}

	return res
}

// An interface for a compress Writer (gzip/zlib/etc)
type WriteCloser interface {
	Write([]byte) (int, error)
	Close() error
}

type compWriter struct {
	gin.ResponseWriter
	nonCompressWriter io.Writer
	compressWriter    WriteCloser
}

func (w compWriter) Write(b []byte) (int, error) {
	cType := w.ResponseWriter.Header().Get("Content-Type")

	if strings.Contains(cType, "application/json") || strings.Contains(cType, "text/html") {
		w.ResponseWriter.Header().Add("Content-Encoding", "gzip")

		len, err := w.compressWriter.Write(b)
		if err != nil {
			return len, err
		}
		return len, w.compressWriter.Close()
	}
	return w.nonCompressWriter.Write(b)
}

func makeGzipWriter(ctx *gin.Context) compWriter {
	gz, err := gzip.NewWriterLevel(ctx.Writer, gzip.BestSpeed)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		ctx.Abort()
	}

	return compWriter{
		ResponseWriter:    ctx.Writer,
		nonCompressWriter: ctx.Writer,
		compressWriter:    gz,
	}
}

func decompressRequest(ctx *gin.Context) {
	if !strings.Contains(ctx.GetHeader("Content-Encoding"), "gzip") {
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer reader.Close()

	ctx.Request.Body = io.NopCloser(reader)
}
