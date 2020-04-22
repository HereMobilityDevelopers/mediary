# mediary
Add interceptors to http.Client and you will be able to
* Dump request and/or response to a `Log`
* Alter your requests before they are sent or responses before they are returned
* Add Tracing information using `Opentracing`/`Jaeger`
* Send metrics to `statsd`

All this and more while using a "regular" `http.Client` 

## Usage
```go
var client *http.Client
client = mediary.Init().AddInterceptors(your interceptor).Build()
client.Get("https://golang.org")
```

## Dump Example
```go
client := mediary.Init().AddInterceptors(dumpInterceptor).Build()
client.Get("https://golang.org")

func dumpInterceptor(req *http.Request, handler mediary.Handler) (*http.Response, error) {
	if bytes, err := httputil.DumpRequestOut(req, true); err == nil {
		fmt.Printf("%s", bytes)

		//GET / HTTP/1.1
		//Host: golang.org
		//User-Agent: Go-http-client/1.1
		//Accept-Encoding: gzip
	}
	return handler(req)
}
```

## Interceptor
**Interceptor** is a function 

<code>type Interceptor func(*http.Request, Handler) (*http.Response, error)</code>

**Handler** is just an alias to <code>http.Roundtripper</code> function called `RoundTrip`

<code>type Handler func(*http.Request) (*http.Response, error) </code>

## Multiple interceptors

It's possible to chain interceptors

```go
client := mediary.Init().
    AddInterceptors(First Interceptor, Second Interceptor).
    AddInterceptors(Third Interceptor).
    Build()

```
This is how it actually works
- First Intereptor: *Request*
    - Second Interceptor: *Request*
        - Third Interceptor: *Request*
            - **Handler** <-- Actual http call
        - Third Interceptor: *Response*
    - Second Interceptor: *Response*
- First Interceptor: *Response*

## Using custom client/transport

If you already have a pre-configured `http.Client` you can use it as well
```go
yourClient := &http.Client{}
yourClientWithInterceptor := mediary.Init().
    WithPreconfiguredClient(yourClient).
    AddInterceptors(your interceptor).
    Build()
```

Read [this](https://github.com/HereMobilityDevelopers/mediary/wiki/Reasoning) for more information
