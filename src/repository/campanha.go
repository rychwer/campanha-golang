package repository

import (
	"strings"
	"fmt"
	"time"
	"campanha-golang/src/model"
	"database/sql"
)

type Campanha struct {
	db *sql.DB
}

func NewRepositoryCampanha(db *sql.DB) *Campanha {
	return &Campanha{db}
}

func (repository Campanha) CriarCampanha(campanha model.Campanha) (uint64, error) {
	statement, erro := repository.db.Prepare("insert into campanha(nomeCampanha, idTimeCoracao, dataVigencia) values (?, ?, ?)",)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(campanha.NomeCampanha, campanha.IdTimeCoracao, campanha.DataVigenciaBanco.Format("2006-01-02"))

	if erro != nil {
		return 0, erro
	}

	ultimoIdInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIdInserido), nil
}

func (repository Campanha) RecuperarTodasCampanhasVigentes() ([]model.Campanha, error) {
	linhas, erro := repository.db.Query("select * from campanha as c where c.dataVigencia >= ? ORDER by c.dataVigencia ASC", time.Now().Format("2006-01-02"))

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	var campanhas []model.Campanha

	for linhas.Next() {
		var campanha model.Campanha

		if erro = linhas.Scan(&campanha.ID, &campanha.NomeCampanha, &campanha.IdTimeCoracao, &campanha.DataVigencia); erro != nil {
			return nil, erro
		}
		
		campanhas = append(campanhas, campanha)
	}

	return campanhas, nil
}

func (repository Campanha) RecuperaTodasCampanhas() ([]model.Campanha, error) {

	linhas, erro := repository.db.Query("select * from campanha",)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var campanhas []model.Campanha

	for linhas.Next() {
		var campanha model.Campanha

		if erro = linhas.Scan(&campanha.ID, &campanha.NomeCampanha, &campanha.IdTimeCoracao, &campanha.DataVigencia); erro != nil {
			return nil, erro
		}

		dataFormatada, _ := time.Parse(time.RFC3339, campanha.DataVigencia)
		campanha.DataVigencia = dataFormatada.Format("02/01/2006")

		campanhas = append(campanhas, campanha)
	}

	return campanhas, nil
}

func (repository Campanha) AtualizarCampanha(campanhaID uint64, campanha model.Campanha) error {
	statement, erro := repository.db.Prepare("update campanha set nomeCampanha = ?, idTimeCoracao = ?, dataVigencia = ? where idCampanha = ?",)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(campanha.NomeCampanha, campanha.IdTimeCoracao, campanha.DataVigencia, campanhaID); erro != nil {
		return erro
	}

	return nil
}

func (repository Campanha) DeletarCampanha(campanhaID uint64) error {
	statement, erro := repository.db.Prepare("delete from campanha where idCampanha = ?",)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(campanhaID); erro != nil {
		return erro
	}

	return nil
}

func (repository Campanha) RecuperarCampanhaPorID(campanhaID uint64) (model.Campanha, error) {
	linhas, erro := repository.db.Query("select * from campanha where idCampanha = ?", campanhaID)
	if erro != nil {
		return model.Campanha{}, erro
	}
	defer linhas.Close()

	var campanha model.Campanha

	for linhas.Next() {
		if erro = linhas.Scan(&campanha.ID, &campanha.NomeCampanha, &campanha.IdTimeCoracao, &campanha.DataVigencia); erro != nil {
			return model.Campanha{}, erro
		}
	}

	return campanha, nil
}

func (repository Campanha) RecuperarCampanhaPorTimeCoracao(idTimeCoracao uint64) (model.Campanha, error) {
	linhas, erro := repository.db.Query("select * from campanha where idTimeCoracao = ?", idTimeCoracao)
	if erro != nil {
		return model.Campanha{}, erro
	}
	defer linhas.Close()

	var campanha model.Campanha

	for linhas.Next() {
		if erro = linhas.Scan(&campanha.ID, &campanha.NomeCampanha, &campanha.IdTimeCoracao, &campanha.DataVigencia); erro != nil {
			return model.Campanha{}, erro
		}
	}

	return campanha, nil
}

func (repository Campanha) AtualizarTodasAsCampanhas(campanhas []model.Campanha) error {

	valueStrings := make([]string, 0, len(campanhas))
	valueStringsIn := make([]string, 0, len(campanhas))
    valueArgs := make([]interface{}, 0, len(campanhas) * 2)

    for index, campanha := range campanhas {
		if index == len(campanhas) - 1 {
			valueStringsIn = append(valueStringsIn, fmt.Sprintf("%d", campanha.ID))
		} else {
			valueStringsIn = append(valueStringsIn, fmt.Sprintf("%d,", campanha.ID))
		}
        valueStrings = append(valueStrings, "WHEN idCampanha = ? THEN ? ")
        valueArgs = append(valueArgs, campanha.ID)
        valueArgs = append(valueArgs, campanha.DataVigencia)
    }

    stmt := fmt.Sprintf("UPDATE campanha SET dataVigencia = CASE %s %s (%s)", strings.Join(valueStrings, ""), "ELSE 0001-01-01 END WHERE idCampanha in", strings.Join(valueStringsIn, ""))

    _, erro := repository.db.Exec(stmt, valueArgs...)

	if erro != nil {
		return erro
	}

    return nil
}