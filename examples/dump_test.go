package examples

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/HereMobilityDevelopers/mediary"
)

func TestDumpInterceptor(t *testing.T) {
	client := mediary.Init().AddInterceptors(dumpInterceptor).Build()
	client.Get("https://golang.org")

}

func dumpInterceptor(req *http.Request, handler mediary.Handler) (*http.Response, error) {
	if bytes, err := httputil.DumpRequestOut(req, true); err == nil {
		fmt.Printf("%s", bytes)

		//GET / HTTP/1.1
		//Host: golang.org
		//User-Agent: Go-http-client/1.1
		//Accept-Encoding: gzip
	}
	resp, err := handler(req)
	// if err == nil {
	// 	if bytes, dumpError := httputil.DumpResponse(resp, true); dumpError == nil {
	// 		fmt.Printf("%s\n\n", bytes)
	// 	}
	// }
	return resp, err
}
