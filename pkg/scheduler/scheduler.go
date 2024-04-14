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
	//Creamos un canal con el fin de captar el cierre de sesion
	stop := make(chan os.Signal, 1)
	//Configuramos el canal stop cuando capte una interrupcion como por ejemplo CTRL + C y SIGTERM para que termine el proceso limpiamente
	//stop es el canal de notififaciones que esperamos del sistema y las otra dos opciones son las señales que esperamos
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	//Esperamos a la señal que sera a traves del canal stop, que puede ser Interrupt o SIGTERM
	<-stop
	//Se ejecuta esto una vez que llegue la señal de stop
	return s.Stop()
}

func (s *SchedulerServer) Stop() error{
	//Cerramos conexion de la DB
	s.dbPool.Close()
	if s.httpServer != nil{
		//Para dar una espera de 10 segundo para que se complete la operaciones antes de cerrar los proceso
		//En caso de que no se completase las operaciones, entonces se cancelaran y le diran al cliente que la operacion ah sido cancelada
		ctx , cancel := context.WithTimeout(context.Background() ,10*time.Second)
		//Este se ejecuto despues del return, por que defer indicare que se ejecute a los ultimo de todo
		defer cancel()
		//Cerramos la conexion del servidor
		return s.httpServer.Shutdown(ctx)
	}
	log.Println("Scheduler server and database pool stopped")
	return nil
}


