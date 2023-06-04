package main

import (
	"context"
	ps "example.com/mod/store"
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

	st, err := ps.New()
	if err != nil {
		log.Fatal(err)
	}

	server := configServer{
		store: st,
	}
	router.HandleFunc("/config/", server.createConfigHandler).Methods("POST")
	router.HandleFunc("/configs/", server.getAllHandler).Methods("GET")
	/*router.HandleFunc("/config/{id}/", server.getConfigHandler).Methods("GET")
	router.HandleFunc("/config/{id}/", server.delConfigHandler).Methods("DELETE")*/
	router.HandleFunc("/config/{id}/{version}/", server.getConfigHandler).Methods("GET")
	router.HandleFunc("/config/{id}/{version}/", server.delConfigHandler).Methods("DELETE")
	router.HandleFunc("/config/{id}/{version}/{labels}", server.getPostByLabel).Methods("GET")

	router.HandleFunc("/group/", server.createGroupHandler).Methods("POST")
	router.HandleFunc("/groups/", server.getAllGroupsHandler).Methods("GET")
	/*router.HandleFunc("/group/{id}/", server.getGroupHandler).Methods("GET")
	router.HandleFunc("/group/{id}/", server.delGroupHandler).Methods("DELETE")*/
	router.HandleFunc("/group/{id}/{version}/", server.getGroupHandler).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/", server.delGroupHandler).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/{labels}", server.getGroupsByLabel).Methods("GET")
	router.HandleFunc("/group/{groupId}/config/{id}/", server.delConfigFromGroupHandler).Methods("DELETE")
	router.HandleFunc("/group/{g_id}/config/{c_id}/", server.addConfigToGroup).Methods("PUT")
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
