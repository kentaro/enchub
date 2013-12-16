package main

import (
	"flag"
	"fmt"
	"github.com/djimenez/iconv-go"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var from = flag.String("from", "utf-8", "The encoding of original content")
var to = flag.String("to", "", "The encoding to be converted")
var backend = flag.String("backend", "", "The host to be proxied")
var port = flag.String("port", "8080", "The port proxy listens on")

func main() {
	flag.Parse()

	if *to == "" || *backend == "" {
		log.Fatalf("Both `to` and `backend` are required")
	}

	http.HandleFunc("/", handler(*backend))
	log.Fatalf("Server Error: %v", http.ListenAndServe(":"+*port, nil))
}

func handler(backend string) func(http.ResponseWriter, *http.Request) {
	client := new(http.Client)

	return func(writer http.ResponseWriter, request *http.Request) {
		proxyRequest := duplicateRequest(request, backend)
		response, err := client.Do(proxyRequest)

		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
		} else {
			for key, values := range response.Header {
				for i := range values {
					writer.Header().Add(key, values[i])
				}
			}

			for _, cookie := range response.Cookies() {
				cookie.Domain = request.Host
			}

			writer.WriteHeader(response.StatusCode)

			data, err := ioutil.ReadAll(response.Body)

			if err != nil {
				writer.Write([]byte(err.Error()))
			} else {
				writer.Write([]byte(data))
			}
		}
	}
}

func convertEncoding(content, from, to string) (result string, err error) {
	result, err = iconv.ConvertString(content, from, to)
	return
}

func replaceCharset(content, to string) (result string) {
	pattern := regexp.MustCompile(`((accept\-)?charset=)(['"])[^'"]+(['"])`)
	result = pattern.ReplaceAllString(content, `${1}${3}`+to+`${4}`)
	return
}

func duplicateRequest(request *http.Request, backend string) *http.Request {
	proxyRequest, err := http.NewRequest(
		request.Method,
		fmt.Sprintf("http://%s%s", backend, request.URL.String()),
		nil,
	)

	if err != nil {
		log.Fatalf("Request Error: %v", err)
	}

	proxyRequest.Host = backend
	proxyRequest.Proto = request.Proto
	proxyRequest.Body = request.Body

	for key, values := range request.Header {
		for i := range values {
			proxyRequest.Header.Add(key, values[i])
		}
	}

	for _, cookie := range proxyRequest.Cookies() {
		cookie.Domain = backend
	}

	return proxyRequest
}
