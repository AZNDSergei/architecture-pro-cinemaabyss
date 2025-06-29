package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

/* ---------- конфигурация ---------- */

var (
	monolithURL = env("MONOLITH_URL", "http://monolith:8099")
	moviesURL   = env("MOVIE_URL", "http://movie-service:8081")
	listenAddr  = ":" + env("PORT", "8080")
)

var (
	gradualMigration = strings.EqualFold(env("GRADUAL_MIGRATION", "false"), "true")
	migrationPercent = mustInt(env("MOVIES_MIGRATION_PERCENT", "0"))
)

/* ---------- утилиты ---------- */

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func mustInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || n > 100 {
		log.Fatalf("MOVIES_MIGRATION_PERCENT должен быть 0–100 (сейчас %q)", s)
	}
	return n
}

func newProxy(target string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("некорректный URL %q: %v", target, err)
	}
	p := httputil.NewSingleHostReverseProxy(u)
	p.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("❌ Proxy error to %s: %v", target, err)
		http.Error(w, "Proxy error", http.StatusBadGateway)
	}
	return p
}

/* ---------- main ---------- */

func main() {
	proxyMonolith := newProxy(monolithURL)
	proxyMovies := newProxy(moviesURL)

	mux := http.NewServeMux()

	// liveness / readiness
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// основной хендлер
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		if strings.HasPrefix(r.URL.Path, "/api/movies") {
			if chooseMovies() {
				log.Printf("→ forwarding to movie-service (%s)", moviesURL)
				proxyMovies.ServeHTTP(w, r)
				return
			}
			log.Printf("→ forwarding to monolith (%s)", monolithURL)
			proxyMonolith.ServeHTTP(w, r)
			return
		}

		log.Printf("→ forwarding to monolith (%s) (fallback)", monolithURL)
		proxyMonolith.ServeHTTP(w, r)
	})

	log.Printf("proxy listening on %s (gradual=%v, percent=%d)", listenAddr, gradualMigration, migrationPercent)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("proxy stopped: %v", err)
	}
}

/* ---------- выбор цели для /api/movies ---------- */

func chooseMovies() bool {
	if !gradualMigration || migrationPercent == 0 {
		return false // всё в монолит
	}
	if migrationPercent == 100 {
		return true // всё в movies
	}
	return rand.Intn(100) < migrationPercent
}
