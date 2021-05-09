package repository

import (
	"errors"
	"strings"
	"fmt"
	"time"
	"campanha-golang/src/model"
	"database/sql"
)

const dataFormatAnoMesDia = "2006-01-02"
const dataFormatDiaMesAno = "02/01/2006"

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

	resultado, erro := statement.Exec(campanha.NomeCampanha, campanha.IdTimeCoracao, campanha.DataVigenciaBanco.Format(dataFormatAnoMesDia))

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
	linhas, erro := repository.db.Query("select * from campanha as c where c.dataVigencia >= ? ORDER by c.dataVigencia ASC", time.Now().Format(dataFormatAnoMesDia))

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
		campanha.DataVigencia = dataFormatada.Format(dataFormatDiaMesAno)

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

	dataFormatada, _ := time.Parse(model.LayoutData, campanha.DataVigencia)

	if _, erro = statement.Exec(campanha.NomeCampanha, campanha.IdTimeCoracao, dataFormatada.Format(dataFormatAnoMesDia), campanhaID); erro != nil {
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

	if campanha.ID == 0 {
		return campanha, errors.New("nenhuma campanha encontrada") 
	}

	dataFormatada, _ := time.Parse(time.RFC3339, campanha.DataVigencia)
	campanha.DataVigencia = dataFormatada.Format(dataFormatDiaMesAno)
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

func (repository Campanha) VerificaCampanhaPorNome(nomeCampanha string) (bool, error) {
	linhas, erro := repository.db.Query("select count(campanha.idCampanha) from campanha where nomeCampanha like ?", nomeCampanha)
	if erro != nil {
		return false, erro
	}
	defer linhas.Close()

	var countCampanhaMesmoNome uint64

	for linhas.Next() {
		if erro = linhas.Scan(&countCampanhaMesmoNome); erro != nil {
			return false, erro
		}
	}

	if countCampanhaMesmoNome > 0 {
		return true, errors.New("já existe uma campanha com o mesmo nome")
	}
	
	return false, nil
}