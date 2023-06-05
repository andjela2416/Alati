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

	const (
		name = "config_service"
	)

	router.HandleFunc("/config/", CountCreateConfig(server.createConfigHandler)).Methods("POST")
	router.HandleFunc("/configs/", CountGetAllConfig(server.getAllHandler)).Methods("GET")

	/*router.HandleFunc("/config/{id}/", server.getConfigHandler).Methods("GET")
	router.HandleFunc("/config/{id}/", server.delConfigHandler).Methods("DELETE")*/
	router.HandleFunc("/config/{id}/{version}/", CountGetConfig(server.getConfigHandler)).Methods("GET")
	router.HandleFunc("/config/{id}/{version}/", CountDelConfig(server.delConfigHandler)).Methods("DELETE")
	router.HandleFunc("/config/{id}/{version}/{labels}", CountGetConfigByLabels(server.getPostByLabel)).Methods("GET")

	router.HandleFunc("/group/", CountCreateGroup(server.createGroupHandler)).Methods("POST")
	router.HandleFunc("/groups/", CountGetAllGroup(server.getAllGroupsHandler)).Methods("GET") //.
	/*router.HandleFunc("/group/{id}/", server.getGroupHandler).Methods("GET")
	router.HandleFunc("/group/{id}/", server.delGroupHandler).Methods("DELETE")*/
	router.HandleFunc("/group/{id}/{version}/", CountGetGroup(server.getGroupHandler)).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/", CountDelGroup(server.delGroupHandler)).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/{labels}", CountGetGroupByLabels(server.getGroupsByLabel)).Methods("GET")
	router.HandleFunc("/group/{groupId}/{g_version}/config/{id}/{c_version}/", CountDelGroupByLabels(server.delConfigFromGroupHandler)).Methods("DELETE")
	router.HandleFunc("/group/{g_id}/{g_version}/config/{c_id}/{c_version}/", CountAddConfigToGroup(server.addConfigToGroup)).Methods("PUT")
	router.HandleFunc("/swagger.yaml", SwaggerHits(server.swaggerHandler)).Methods("GET")

	// s c r a p e m e t r i c s f rom s e r v i c e , show UI on l o c a l h o s t : 9 0 9 0
	router.Path("/metrics").Handler(metricsHandler())
	// SwaggerUI
	optionsDevelopers := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := middleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

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
