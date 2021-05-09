package rotas

import (
	"campanha-golang/src/controller"
	"net/http"
)

var rotasCampanha = []Rota {
	{
		URI: "/campanha",
		Metodo: http.MethodPost,
		Funcao: controller.PostCampanha,
	},
	{
		URI: "/campanha",
		Metodo: http.MethodGet,
		Funcao: controller.GetCampanhas,
	},
	{
		URI: "/campanha/{idCampanha}",
		Metodo: http.MethodPut,
		Funcao: controller.PutCampanha,
	},
	{
		URI: "/campanha/{idCampanha}",
		Metodo: http.MethodDelete,
		Funcao: controller.DeleteCampanha,
	},
	{
		URI: "/campanha/{idCampanha}",
		Metodo: http.MethodGet,
		Funcao: controller.GetCampanhaById,
	},
	{
		URI: "/time/{idTimeCoracao}",
		Metodo: http.MethodGet,
		Funcao: controller.GetCampanhaByTimeCoracao,
	},
}