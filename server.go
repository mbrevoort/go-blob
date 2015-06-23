package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type Data struct {
	sync.Mutex
	Store map[string][]byte
	Types map[string]string
}

const TENMB = 1048576 * 10

func main() {
	data := Data{
		Store: map[string][]byte{},
		Types: map[string]string{},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		key := r.URL.Path
		fmt.Printf("%s %s\n", r.Method, key)

		switch r.Method {
		case "GET":
			data.Lock()
			payload, ok := data.Store[key]

			if !ok {
				w.WriteHeader(http.StatusNotFound)
				data.Unlock()
				return
			}
			contentType, ok := data.Types[key]
			if ok {
				w.Header().Set("Content-type", contentType)
			}
			data.Unlock()

			w.WriteHeader(200)
			w.Write(payload)
		case "PUT":
			payload, err := ioutil.ReadAll(io.LimitReader(r.Body, TENMB))
			if err != nil {
				respond(w, http.StatusInternalServerError, err.Error())
				return
			}
			contentType := r.Header.Get("Content-Type")

			data.Lock()
			data.Store[key] = payload
			data.Types[key] = contentType
			data.Unlock()
			w.WriteHeader(http.StatusCreated)
		case "DELETE":
			fmt.Printf("DELETE %s\n", key)
			data.Lock()
			delete(data.Store, key)
			delete(data.Types, key)
			data.Unlock()
			w.WriteHeader(http.StatusOK)
		default:
			respond(w, http.StatusNotFound, key+" not found")
		}
	})

	listen := fmt.Sprintf(":%d", 3000)
	fmt.Println("Listening " + listen)
	http.ListenAndServe(listen, nil)
}

func respond(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}
