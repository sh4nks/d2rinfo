package server

import (
	"d2rinfo/config"
	"d2rinfo/controller"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"

	otter "github.com/maypok86/otter/v2"
)

const (
	D2EMU_TZ_API     = "https://d2emu.com/api/v1/tz"
	D2EMU_DCLONE_API = "https://d2emu.com/api/v1/dclone"
)

// UserAgentMiddleware checks for a specific User-Agent header.
func UserAgentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "D2RLoader" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type Server struct {
	Router *chi.Mux
	Config *config.Config
}

func New(cfg *config.Config) *Server {
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID, middleware.Logger, middleware.Recoverer,
		UserAgentMiddleware,
	)
	router.Use(httprate.LimitByIP(cfg.RateLimit, time.Minute))
	cache := otter.Must(&otter.Options[string, any]{
		InitialCapacity:  100,
		ExpiryCalculator: FixedIntervalExpiry{},
	})

	// register our controllers here
	d2rCtrl := controller.NewD2RInfoController(cfg, cache)

	router.Get("/api/d2rinfo", d2rCtrl.GetD2RInfoData)
	return &Server{
		Config: cfg,
		Router: router,
	}
}

func (srv *Server) StartServer() {
	addr := fmt.Sprintf("%s:%d", srv.Config.Host, srv.Config.Port)
	log.Printf("Running on http://%s/ (Press CTRL+C to quit)", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Router))
}
