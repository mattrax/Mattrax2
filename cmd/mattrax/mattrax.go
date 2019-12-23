package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/api"
	"github.com/mattrax/Mattrax/internal/boltdb"
)

func main() {
	db, err := boltdb.Init()
	if err != nil {
		panic(err) // TODO
	}
	defer db.Close()

	userService, err := boltdb.NewUserService(db)
	if err != nil {
		panic(err) // TODO
	}

	policyService, err := boltdb.NewPolicyService(db)
	if err != nil {
		panic(err) // TODO
	}

	settingsService, err := boltdb.NewSettingsService(db)
	if err != nil {
		panic(err) // TODO
	}

	server := mattrax.Server{
		Config: mattrax.Config{
			Port:            8000,
			Domains:         []string{"mdm.otbeaumont.me", "enterpriseenrollment.otbeaumont.me"},
			CertFile:        "./certs/server.crt",
			KeyFile:         "./certs/server.key",
			DevelopmentMode: true,
		},
		UserService:     userService,
		PolicyService:   policyService,
		SettingsService: settingsService,
	}

	r := mux.NewRouter()

	api.InitAPI(server, r)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Mattrax Server. Created By Oscar Beaumont.")
	})

	// TODO: Gracefull Shutdown
	port := strconv.Itoa(server.Config.Port)
	log.Println("Listening on port " + port + "...")
	log.Fatal(http.ListenAndServeTLS(":"+port, server.Config.CertFile, server.Config.KeyFile, r))
}
