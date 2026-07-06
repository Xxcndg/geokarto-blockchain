package main

import (
	"fmt"
	"net/http"
	"os"
)

var (
	KoinBeredarGlobal = 1295000000.0
	IndeksBlok        = 0
	MaksimalSupply    = 1300000000.0
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		sisaKuota := MaksimalSupply - KoinBeredarGlobal
		fmt.Fprintf(w, "%f:%d", sisaKuota, IndeksBlok)
	})

	http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
		if KoinBeredarGlobal < MaksimalSupply {
			KoinBeredarGlobal += 1.0
			IndeksBlok++
			w.Write([]byte("MINED_ACCEPTED"))
		} else {
			w.Write([]byte("SUPPLY_EXHAUSTED"))
		}
	})

	fmt.Println("Server Geokarto Global Cloud Aktif di Port " + port)
	http.ListenAndServe(":"+port, nil)
}
EOF