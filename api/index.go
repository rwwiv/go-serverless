package api

import (
	"encoding/json"
	"net/http"
)

type Foo struct {
	Hello string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	foo := Foo{Hello: "World"}
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(foo)
	w.Write(res)
}
