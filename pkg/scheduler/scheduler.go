package scheduler

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/FxIvan/grcp-golang/pkg/common"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CommandRequest struct {
	Command    string `json:"command"`
	ScheduleAt string `json:"schedule_at"`
}

type Task struct {
	Id          string
	Command     string
	ScheduledAt pgtype.Timestamp
	PickedAt    pgtype.Timestamp
	StartedAt   pgtype.Timestamp
	CompletedAt pgtype.Timestamp
	FailedAt    pgtype.Timestamp
}

type SchedulerServer struct {
	serverPort         string
	dbConnectionString string
	dbPool             *pgxpool.Pool
	ctx                context.Context
	cancel             context.CancelFunc
	httpServer         *http.Server
}

/*
 pgtype.Timestamp por que postgresql solo acepta formato UTC
 pgxpool permite manejar varias conexiones a la base de datos, pgxpool.Pool permite ejecutar varias gorutinas de manera concurrente sin crear una nueva conexión a la base de datos
*/

func (s *SchedulerServer) Start() error {
	var err error

	//nos conectamos a la DB que dentro tiene hasta 5 intentos
	s.dbPool, err = common.ConnectToDatabase(s.ctx, s.dbConnectionString)
	if err != nil {
		return err
	}

	//Definimos rutas
	//http.HandleFunc("/schedule", )
	//http.HandleFunc("/status/",)

	s.httpServer = &http.Server{
		Addr: s.serverPort,
	}

	log.Printf("Starting scheduler server on %s\n", s.serverPort)

	//Iniciamos el servidor HTTP goroutine
	//Es decir que se ejecute varias veces independientemente de otra ejecucion
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s\n", err)
		}
	}()

	return s.awaitShutdown()
}

// Funcion de interrupcion para detener el servidor de manera segura
// Es una señal SIGTERM
func (s *SchedulerServer) awaitShutdown() error {
	//
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	return s.Stop()
}
