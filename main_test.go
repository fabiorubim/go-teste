package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guilhermeonrails/api-go-gin/controllers"
	"github.com/guilhermeonrails/api-go-gin/database"
	"github.com/guilhermeonrails/api-go-gin/models"
	"github.com/stretchr/testify/assert"
)

var ID int

func SetupDasRotasDeTeste() *gin.Engine {
	//Melhora a exibição dos testes
	gin.SetMode(gin.ReleaseMode)
	rotas := gin.Default()
	return rotas
}

// func TestFalhador(t *testing.T) {
// 	t.Fatalf("Teste falhou de propósito!")
// }

func TestVerificaStatusCodeDaSaudacaoComParametro(t *testing.T) {
	r := SetupDasRotasDeTeste()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/gui", nil)
	//NewRecorder implementa uma interface de response writer que armazena toda a resposta
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	// if resposta.Code != http.StatusOK {
	// 	t.Fatalf("Status erros: valor recebido recebido foi %d e o esperado era %d", resposta.Code, http.StatusOK)
	// }
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais!")
	mockDaResposta := `{"API diz:":"E ai gui, tudo beleza?"}`
	respostaBody, _ := ioutil.ReadAll(resposta.Body)
	assert.Equal(t, mockDaResposta, string(respostaBody))
}

func TestListandoTodosOsAlunosHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupDasRotasDeTeste()
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	req, _ := http.NewRequest("GET", "/alunos", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code)
	fmt.Println(resposta.Body)
}

func CriaAlunoMock() {
	aluno := models.Aluno{Nome: "Nome do Aluno Teste", CPF: "12345678910", RG: "12345679"}
	database.DB.Create(&aluno)
	ID = int(aluno.ID)

}

func DeletaAlunoMock() {
	var aluno models.Aluno
	database.DB.Delete(&aluno, ID)
}

func TestBuscaAlunoPorCPFHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupDasRotasDeTeste()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	req, _ := http.NewRequest("GET", "/alunos/cpf/12345678910", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code)
}

//Executar somente um teste: go test -run TestBuscaAlunoPorCPFHandler

func TestBuscaAlunoPorIDHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupDasRotasDeTeste()
	r.GET("/alunos/:id", controllers.BuscaAlunoPorID)
	pathDaBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", pathDaBusca, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	var alunoMock models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &alunoMock)
	fmt.Println(alunoMock.Nome)
	assert.Equal(t, "Nome do Aluno Teste", alunoMock.Nome, "Os nomes devem ser iguais!")
	assert.Equal(t, "12345678910", alunoMock.CPF)
	assert.Equal(t, "12345679", alunoMock.RG)
}

func TestDeletaAlunoHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	r := SetupDasRotasDeTeste()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	pathBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", pathBusca, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code)
}

func TestEditaUmAlunoHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupDasRotasDeTeste()
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	aluno := models.Aluno{Nome: "Nome do Aluno Teste", CPF: "47345678910", RG: "123456009"}
	valorJson, _ := json.Marshal(aluno)
	pathParaEditar := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", pathParaEditar, bytes.NewBuffer(valorJson))
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	var alunoMockAtualizado models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &alunoMockAtualizado)
	assert.Equal(t, "47345678910", alunoMockAtualizado.CPF)
	assert.Equal(t, "123456009", alunoMockAtualizado.RG)
	assert.Equal(t, "Nome do Aluno Teste", alunoMockAtualizado.Nome)
}
