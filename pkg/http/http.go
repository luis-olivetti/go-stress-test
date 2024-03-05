package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Run(engine *gin.Engine, addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// Inicializando o servidor em uma goroutine para que
	// não bloqueie o tratamento de encerramento gracioso abaixo
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Aguardando sinal de interrupção para encerrar o servidor com graciosidade
	// com um tempo limite de 5 segundos.
	quit := make(chan os.Signal, 1)
	// kill (sem parâmetro) envia syscall.SIGTERM por padrão
	// kill -2 é syscall.SIGINT. Geralmente Ctrl + C
	// kill -9 é syscall.SIGKILL, mas não pode ser capturado, então não precisa adicioná-lo
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// O contexto é usado para informar ao servidor que ele tem 5 segundos para finalizar
	// a solicitação que está atualmente tratando
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
