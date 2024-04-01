package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/LegendaryB/go-wake-agent/config"
	"github.com/LegendaryB/go-wake-agent/wol"

	"github.com/gorilla/mux"
)

type WakeRequest struct {
	Address string
}

func wakeDeviceHandler(allowedClients []*config.AllowedClient) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		authorizationHeaderValue := req.Header.Get("Authorization")

		if len(authorizationHeaderValue) == 0 {
			http.Error(res, "'Authorization' header missing.", http.StatusUnauthorized)
			return
		}

		clientNameHeaderValue := req.Header.Get("ClientName")

		if len(clientNameHeaderValue) == 0 {
			http.Error(res, "'ClientName' header missing.", http.StatusBadRequest)
			return
		}

		allowedClientIndex := slices.IndexFunc(allowedClients, func(ac *config.AllowedClient) bool {
			return ac.Name == clientNameHeaderValue
		})

		if allowedClientIndex == -1 {
			http.Error(res, "Unknown client, not registered on server side", http.StatusUnauthorized)
			return
		}

		allowedClient := allowedClients[allowedClientIndex]

		if authorizationHeaderValue != allowedClient.ApiToken {
			http.Error(res, "Invalid token", http.StatusUnauthorized)
			return
		}

		var wakeRequest WakeRequest

		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&wakeRequest)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		wol.SendWakeOnLANPacket(wakeRequest.Address)
	}
}

func main() {
	config, err := config.NewConfiguration()

	if err != nil {
		log.Panicf("Failed to load app configuration. %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/wakeDevice", wakeDeviceHandler(config.AllowedClients)).Methods("PATCH")

	listenAddr := fmt.Sprintf(":%d", config.Application.ListenPort)

	http.ListenAndServe(listenAddr, router)
}
