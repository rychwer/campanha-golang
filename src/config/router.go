package config

import (
	"campanha-golang/src/config/rotas"
	"github.com/gorilla/mux"
)

func GerarRouter() *mux.Router {
	router := mux.NewRouter()
	return rotas.Configurar(router)
}