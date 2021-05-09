package model

import (
	"time"
	"errors"
)

const LayoutData = "02/01/2006"

type Campanha struct {
	ID uint64 `json:"idCampanha,omitempty"`
	NomeCampanha string `json:"nomeCampanha"`
	IdTimeCoracao uint64 `json:"timeCoracao,omitempty"`
	DataVigencia string `json:"dataVigencia"`
	DataVigenciaBanco time.Time `json:"-"`
}

func (campanha *Campanha) Validar() error {
	if campanha.NomeCampanha == "" {
		return errors.New("o nome da campanha deve ser obrigatório")
	}
	if campanha.DataVigencia == "" {
		return errors.New("a campanha deve conter uma data de vigência")
	}
	if _, erro := time.Parse(LayoutData, campanha.DataVigencia); erro != nil {
		return errors.New("a campanha deve conter uma data de vigência no formato DD/MM/YYYY")
	} 
	
	return nil
}

func (campanha *Campanha) FormatarData() {
	valor, _ := time.Parse(LayoutData, campanha.DataVigencia)
	campanha.DataVigenciaBanco = valor
}

func (campanha *Campanha) AtualizaDataVigenciaCampanha(campanhasVigentes []Campanha) []Campanha{
	for i := 0; i < len(campanhasVigentes); i++ {
        dataVigenciaAtualizada := converteStringToTime(campanhasVigentes[i].DataVigencia)
		dataAtualizada := dataVigenciaAtualizada.AddDate(0,0,1)
		campanhasVigentes[i].DataVigenciaBanco = dataAtualizada
		campanhasVigentes[i].DataVigencia = dataAtualizada.Format("2006-01-02")
    }

	var campanhasDataVigenciaAtualizadas []Campanha
	dataBanco, _ := time.Parse("02/01/2006", campanha.DataVigencia)
	campanha.DataVigencia = dataBanco.Format("2006-01-02")
    campanhasDataVigenciaAtualizadas = append(campanhasDataVigenciaAtualizadas, *campanha)

	for i := 0; i < len(campanhasVigentes); i++ {
        for j := 0; j < len(campanhasDataVigenciaAtualizadas); j++ {
            if campanhasVigentes[i].DataVigencia == campanhasDataVigenciaAtualizadas[j].DataVigencia {
				dataAtualizada := campanhasVigentes[i].DataVigenciaBanco.AddDate(0,0,1)
                campanhasVigentes[i].DataVigenciaBanco = dataAtualizada
				campanhasVigentes[i].DataVigencia = dataAtualizada.Format("2006-01-02")
            }
        }
        campanhasDataVigenciaAtualizadas = append(campanhasDataVigenciaAtualizadas, campanhasVigentes[i])
    }

	campanha.DataVigencia = dataBanco.Format("02/01/2006")

	return removeIndex(campanhasDataVigenciaAtualizadas, 0)
}

func converteStringToTime(stringTime string) time.Time {
	valor, _ := time.Parse(time.RFC3339, stringTime)
	return valor
}

func removeIndex(s []Campanha, index int) []Campanha {
	return append(s[:index], s[index+1:]...)
}
