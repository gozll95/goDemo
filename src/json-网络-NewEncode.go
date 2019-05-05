package main

import (
	"encoding/json"
	"net/http"
)

var s store

type bigTask struct {
	ID     int
	Status string
	Data   map[string]string
}

type store struct {
	task []bigTask
}

var bigTasks []bigTask

func init() {

	bigTasks = []bigTask{
		{0, "wait", map[string]string{"data": "http://localhost/0/data"}},
		{1, "done", map[string]string{"data": "http://localhost/1/data"}},
		{2, "dead", map[string]string{"data": "http://localhost/2/data"}},
	}
}
func main() {
	http.HandleFunc("/big_task", bigtask)

	http.ListenAndServe(":10004", nil)
}

func bigtask(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	if e := enc.Encode(bigTasks); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)

	}

}

/*
curl http://localhost:10004/big_task

[{"ID":0,"Status":"wait","Data":{"data":"http://localhost/0/data"}},{"ID":1,"Status":"done","Data":{"data":"http://localhost/1/data"}},{"ID":2,"Status":"dead","Data":{"data":"http://localhost/2/data"}}]
*/
