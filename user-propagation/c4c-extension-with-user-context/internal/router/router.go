package router

import (
	"net/http"

	"github.com/SAP-samples/kyma-runtime-extension-samples/user-propagation/c4c-extension-with-user-context/internal/handlers"
	"github.com/gorilla/mux"
)

func New() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", handlers.NewDispatcher().CreateTask).Methods(http.MethodPost)
	return router
}
