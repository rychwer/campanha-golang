package controller

import (
	"strconv"
	"github.com/gorilla/mux"
	"encoding/json"
	"campanha-golang/src/model"
	"io/ioutil"
	"campanha-golang/src/repository"
	"campanha-golang/src/handle"
	"campanha-golang/src/banco"
	"net/http"
)

func PostCampanha(w http.ResponseWriter, r *http.Request) {
	bodyRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		handle.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var campanha model.Campanha

	if erro = json.Unmarshal(bodyRequest, &campanha); erro != nil {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = campanha.Validar(); erro != nil {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	}

	campanha.FormatarData()

	db, erro := banco.ConectarBanco()
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryCampanha(db)

	hasCampanhaNome, erro := repository.VerificaCampanhaPorNome(campanha.NomeCampanha)
	if erro != nil && !hasCampanhaNome {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	if hasCampanhaNome {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	}
	
	campanhasVigentes, erro := repository.RecuperarTodasCampanhasVigentes()
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if len(campanhasVigentes) > 0 {
		campanhasVigentesParaAtualizar := campanha.AtualizaDataVigenciaCampanha(campanhasVigentes)
		if erro = repository.AtualizarTodasAsCampanhas(campanhasVigentesParaAtualizar); erro != nil {
			handle.Erro(w, http.StatusInternalServerError, erro)
			return
		}
	}
	
	campanha.ID, erro = repository.CriarCampanha(campanha)
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	
	handle.JSON(w, http.StatusCreated, campanha)
}

func GetCampanhas(w http.ResponseWriter, request *http.Request) {
	db, erro := banco.ConectarBanco()
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryCampanha(db)
	campanhas, erro := repository.RecuperaTodasCampanhas()
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	handle.JSON(w, http.StatusOK, campanhas)
}

func PutCampanha(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	campanhaID, erro := strconv.ParseUint(parametros["idCampanha"], 10, 64)
	if erro != nil {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	}

	bodyRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		handle.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var campanha model.Campanha

	if erro = json.Unmarshal(bodyRequest, &campanha); erro != nil {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro := campanha.Validar(); erro != nil {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.ConectarBanco()
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryCampanha(db)
	
	if erro := repository.AtualizarCampanha(campanhaID, campanha); erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	handle.JSON(w, http.StatusNoContent, nil)
}

func DeleteCampanha(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	campanhaID, erro := strconv.ParseUint(parametros["idCampanha"], 10, 64)
	if erro != nil {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.ConectarBanco()
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryCampanha(db)
	if erro := repository.DeletarCampanha(campanhaID); erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	handle.JSON(w, http.StatusNoContent, nil)
}

func GetCampanhaById(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	campanhaID, erro := strconv.ParseUint(parametros["idCampanha"], 10, 64)
	if erro != nil {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := banco.ConectarBanco()
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryCampanha(db)
	campanha, erro := repository.RecuperarCampanhaPorID(campanhaID)
	if erro != nil && campanha.ID == 0 {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	} else if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	handle.JSON(w, http.StatusOK, campanha)
}

func GetCampanhaByTimeCoracao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	idTimeCoracao, erro := strconv.ParseUint(parametros["idTimeCoracao"], 10, 64)
	if erro != nil {
		handle.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := banco.ConectarBanco()
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repository.NewRepositoryCampanha(db)
	campanha, erro := repository.RecuperarCampanhaPorTimeCoracao(idTimeCoracao)
	if erro != nil {
		handle.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	handle.JSON(w, http.StatusOK, campanha)
}