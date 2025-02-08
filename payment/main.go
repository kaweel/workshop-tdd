package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kaweel/workshop-tdd/payment/clock"
	"github.com/kaweel/workshop-tdd/payment/handler"
	"github.com/kaweel/workshop-tdd/payment/messaging"
	"github.com/kaweel/workshop-tdd/payment/service"
	"github.com/kaweel/workshop-tdd/payment/storage"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	wait := time.Duration(time.Second * 15)

	connStr := "sqlserver://sa:SuperStrong@Passw0rd@localhost:1433?database=tdd-workshop&encrypt=disable&parseTime=True&time_zone=UTC"
	db, err := gorm.Open(sqlserver.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect MSSQL: %v", err)
	}

	orderStorage := storage.NewOrderStorage(db)
	paymentTranasctionStorage := storage.NewPaymentTranasctionStorage(db)
	kafkaProducer := messaging.NewKafkaProducer()
	clock := clock.NewClock()
	paymentService := service.NewService(orderStorage, paymentTranasctionStorage, kafkaProducer, clock)
	handlerPayment := handler.NewHandler(paymentService)

	r := mux.NewRouter()
	r.HandleFunc("/payment", handlerPayment.Payment()).GetMethods()
	// Add your routes as needed

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Println("Starting server on port 8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
