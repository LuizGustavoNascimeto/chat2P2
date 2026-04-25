/*
Descrição: Fornece logger singleton baseado em Zap e helpers de erro para projeto.
Autor: Luizg
Data de criação: 2026-04-25
Data de atualização: 2026-04-25
*/

package logger

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

var (
	instance *zap.Logger
	once     sync.Once
)

// Init inicializa logger global uma única vez.
// Entradas: ambiente lógico (development/production).
// Saída: nenhuma; panic em falha de inicialização.
func Init(env string) {
	once.Do(func() {
		var err error
		// if env == "production" {
		instance, err = zap.NewProduction()
		// } else {
		// 	instance, err = zap.NewDevelopment()
		// }
		if err != nil {
			panic(err)
		}
	})
}

// Get retorna instância global de logger já inicializada.
// Entradas: nenhuma.
// Saída: ponteiro para zap.Logger; panic se Init não foi chamado.
func Get() *zap.Logger {
	if instance == nil {
		panic("logger não inicializado, chame Init() primeiro")
	}
	return instance
}

// Sync descarrega buffers do logger ativo.
// Entradas: nenhuma.
// Saída: nenhuma.
func Sync() {
	if instance != nil {
		instance.Sync()
	}
}

// LogError cria erro, registra no logger e retorna valor para fluxo do chamador.
// Entradas: mensagem de erro.
// Saída: erro criado a partir da mensagem.
func LogError(message string) error {
	err := errors.New(message)
	Get().Error("datagram error", zap.Error(err))
	return err
}
