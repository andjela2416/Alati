package main

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//docker run -p 8000:8000 gorest

/*type Service struct {
	data map[string]*Config `json:"data"`
}

type Config struct {
	entries map[string]string `json:"entries"`
}*/

func main() {
	fmt.Println("Hello world")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	server := configServer{
		data:      map[string]*Config{},
		groupData: map[string]*Group{},
	}
	router.HandleFunc("/config/", server.createConfigHandler).Methods("POST")
	router.HandleFunc("/configs/", server.getAllHandler).Methods("GET")
	router.HandleFunc("/config/{id}/", server.getConfigHandler).Methods("GET")
	router.HandleFunc("/config/{id}/", server.delConfigHandler).Methods("DELETE")

	router.HandleFunc("/group/", server.createGroupHandler).Methods("POST")
	router.HandleFunc("/group/{groupId}/config{id}/", server.addConfigToGroup).Methods("PUT")
	router.HandleFunc("/swagger.yaml", server.swaggerHandler).Methods("GET")

	// SwaggerUI
	optionsDevelopers := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := middleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

	// ReDoc
	// optionsShared := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	// sharedDocumentationHandler := middleware.Redoc(optionsShared, nil)
	// router.Handle("/docs", sharedDocumentationHandler)

	// start server

	// start server

	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
