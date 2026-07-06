package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Struktur data blok kriptografi murni
type Block struct {
	Index        int
	Timestamp    string
	WalletTarget string
	PrevHash     string // KUNCI UTAMA: Menyimpan sidik jari blok sebelum dirinya
	Hash         string // KUNCI UTAMA: Sidik jari unik dari blok ini sendiri
}

var (
	KoinBeredarGlobal = 1295000000.0
	MaksimalSupply    = 1300000000.0
	Blockchain        []Block // Tempat menyimpan rantai blok di dalam memori cloud
)

// Algoritma Matematika untuk menghitung Hash murni berdasarkan data transaksi
func HitungHashBlok(b Block) string {
	record := strconv.Itoa(b.Index) + b.Timestamp + b.WalletTarget + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Fungsi pencetakan Blok Pertama di Dunia (Genesis Block) sebagai jangkar rantai
func BuatBlokGenesis() {
	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().Format("02-01-2006 15:04:05 WIB"),
		WalletTarget: "SYSTEM_GENESIS_GATE",
		PrevHash:     "0000000000000000000000000000000000000000000000000000000000000000", // Kosong karena tidak ada blok sebelumnya
	}
	genesis.Hash = HitungHashBlok(genesis)
	Blockchain = append(Blockchain, genesis)
}

func main() {
	// Jalankan blok genesis sesaat setelah server cloud menyala
	BuatBlokGenesis()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// API Jalur 1: Melihat seluruh mata rantai blok beserta sisa kuota koin secara transparan
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		sisaKuota := MaksimalSupply - KoinBeredarGlobal
		
		fmt.Fprintf(w, "=== GEOKARTO CRYPTO CHAIN REGISTER ===\n")
		fmt.Fprintf(w, "SISA PASOKAN KOIN GLOBAL: %.2f GEK\n", sisaKuota)
		fmt.Fprintf(w, "TOTAL BLOK TERENKRIPSI: %d\n\n", len(Blockchain))
		
		for _, b := range Blockchain {
			fmt.Fprintf(w, "[BLOK #%d]\n", b.Index)
			fmt.Fprintf(w, " Waktu Kunci : %s\n", b.Timestamp)
			fmt.Fprintf(w, " Penambang   : %s\n", b.WalletTarget)
			fmt.Fprintf(w, " PREV HASH   : %s\n", b.PrevHash)
			fmt.Fprintf(w, " HASH BLOK   : %s\n", b.Hash)
			fmt.Fprintf(w, "------------------------------------------------------------\n")
		}
	})

	// API Jalur 2: Menambang 1 koin baru dengan sistem rantai matematika mengunci
	http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
		if KoinBeredarGlobal >= MaksimalSupply {
			w.Write([]byte("REJECTED: SUPPLY_EXHAUSTED"))
			return
		}

		// Mengambil blok terakhir yang ada di ujung rantai
		blokTerakhir := Blockchain[len(Blockchain)-1]
		
		// Membuat blok baru yang WAJIB mengikat Hash dari blok terakhir tersebut
		blokBaru := Block{
			Index:        blokTerakhir.Index + 1,
			Timestamp:    time.Now().Format("02-01-2006 15:04:05 WIB"),
			WalletTarget: "GK_PEER_NODE_" + strconv.Itoa(rand.Intn(9000)+1000),
			PrevHash:     blokTerakhir.Hash, // <--- IKATAN RANTAI MATEMATIKA MUTLAK
		}
		blokBaru.Hash = HitungHashBlok(blokBaru)

		// Suntikkan blok baru yang sudah sah ke dalam pangkalan rantai global
		Blockchain = append(Blockchain, blokBaru)
		KoinBeredarGlobal += 1.0

		w.Write([]byte(fmt.Sprintf("MINED_ACCEPTED: BLOK #%d LOCKED WITH HASH %s", blokBaru.Index, blokBaru.Hash[:8])))
	})

	fmt.Println("Server Kriptografi Rantai Blok Geokarto Aktif di Port " + port)
	http.ListenAndServe(":"+port, nil)
}
