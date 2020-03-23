package mediary

import (
	"net/http"
)

type customRoundTripper struct {
	inner             http.RoundTripper
	unitedInterceptor Interceptor
}

func prepareCustomRoundTripper(actual http.RoundTripper, interceptors ...Interceptor) http.RoundTripper {
	return &customRoundTripper{
		inner:             actual,
		unitedInterceptor: uniteInterceptors(interceptors),
	}
}

func (crt *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return crt.unitedInterceptor(req, crt.inner.RoundTrip)
}

func uniteInterceptors(interceptors []Interceptor) Interceptor {
	if len(interceptors) == 0 {
		return func(req *http.Request, handler Handler) (*http.Response, error) {
			// That's why we needed an alias to http.RoundTripper.RoundTrip
			return handler(req)
		}
	}

	return func(req *http.Request, handler Handler) (*http.Response, error) {
		tailhandler := func(innerReq *http.Request) (*http.Response, error) {
			unitedInterceptor := uniteInterceptors(interceptors[1:])
			return unitedInterceptor(req, handler)
		}
		headInterceptor := interceptors[0]
		return headInterceptor(req, tailhandler)
	}
}
