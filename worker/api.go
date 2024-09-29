package worker

import (
	"fmt"
	"log"
	"net/http"
)

type Api struct {
	Address string
	Port    int
	Worker  *Worker
	Router  *http.ServeMux
}

type ErrResponse struct {
	HTTPStatusCode int
	Message        string
}

func (a *Api) initRouter() {
	a.Router = http.NewServeMux()
	// High-level (i.e., the genesis) handler:
	// You can't nest `Handlefunc`s inside `Handlefunc`s. There should be only one genesis Handlefunc.
	a.Router.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			a.GetTasksHandler(w, r)
		case "POST":
			a.StartTaskHandler(w, r)
		default:
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	})

	a.Router.HandleFunc("/tasks/{taskID}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			a.StopTaskHandler(w, r)
		default:
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (a *Api) ServeAPI() {
	a.initRouter()
	log.Printf("Listening on %s:%d", a.Address, a.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", a.Address, a.Port), a.Router)
}
