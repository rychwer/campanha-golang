package rotas

import (
	"campanha-golang/src/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type Rota struct {
	URI string
	Metodo string
	Funcao func(http.ResponseWriter, *http.Request)
}

func Configurar(router *mux.Router) *mux.Router {
	rotas := rotasCampanha
	for _, rota := range rotas {
		router.HandleFunc(rota.URI, rota.Funcao).Methods(rota.Metodo)
	}

	router.Use(middleware.Logger)

	return router
}