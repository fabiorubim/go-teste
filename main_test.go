package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guilhermeonrails/api-go-gin/controllers"
	"github.com/stretchr/testify/assert"
)

func SetupDasRotasDeTeste() *gin.Engine {
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
