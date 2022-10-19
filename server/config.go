package main

import "github.com/namsral/flag"

type config struct {
	bindAddr string
	moneroPayAddr string
	callbackAddr string
	amount uint64
	postgresCS string
}

var Conf config

func loadConfig() {
	flag.StringVar(&Conf.bindAddr, "bind", ":5000", "Bind address:port")
	flag.StringVar(&Conf.moneroPayAddr, "moneropay-addr", "http://moneropay:5000/receive", "MoneroPay endpoint.")
	flag.StringVar(&Conf.callbackAddr, "callback-addr", "http://shady:1337", "Callback base URL.")
	flag.Uint64Var(&Conf.amount, "amount", 1000000000000, "Monero amount to request (in piconeroj)")
	flag.StringVar(&Conf.postgresCS, "postgresql", "postgresql://postgres:changeMePlease@localhost:5432/shady",
	    "PostgreSQL connection string")
	flag.Parse()
}
