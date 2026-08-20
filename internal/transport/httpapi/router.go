package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"industrial-edge-protocol/internal/domain"
	"industrial-edge-protocol/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	app *service.Application
	log *slog.Logger
}

func NewRouter(app *service.Application, log *slog.Logger) http.Handler {
	s := &Server{app: app, log: log}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/readyz", s.health)
	mux.HandleFunc("/debug/pprof/", s.notFound)
	mux.HandleFunc("/api/v1/sites", s.sites)
	mux.HandleFunc("/api/v1/devices", s.devices)
	mux.HandleFunc("/api/v1/gateways", s.gateways)
	mux.HandleFunc("/api/v1/points", s.points)
	mux.HandleFunc("/api/v1/acquisition-plans", s.plans)
	mux.HandleFunc("/api/v1/rules", s.rules)
	mux.HandleFunc("/api/v1/alarms", s.alarms)
	mux.HandleFunc("/api/v1/commands", s.commands)
	mux.HandleFunc("/api/v1/telemetry/raw", s.telemetryRaw)
	mux.HandleFunc("/api/v1/telemetry/aggregates", s.telemetryAggregate)
	mux.HandleFunc("/api/v1/rollouts", s.rollouts)
	return logging(mux, log)
}
func logging(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Info("http request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	write(w, http.StatusOK, s.app.Health())
}
func (s *Server) notFound(w http.ResponseWriter, _ *http.Request) {
	writeError(w, http.StatusNotFound, "not found")
}
func (s *Server) sites(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		write(w, http.StatusOK, s.app.Repo.ListSites(r.Context()))
		return
	}
	if r.Method == http.MethodPost {
		var v domain.Site
		if decode(r, &v) != nil {
			writeError(w, 400, "invalid json")
			return
		}
		created, err := s.app.CreateSite(r.Context(), v)
		respond(w, created, err)
		return
	}
	writeError(w, 405, "method not allowed")
}
func (s *Server) devices(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		write(w, 200, s.app.Repo.ListDevices(r.Context()))
		return
	}
	if r.Method == http.MethodPost {
		var v domain.Device
		if decode(r, &v) != nil {
			writeError(w, 400, "invalid json")
			return
		}
		created, err := s.app.CreateDevice(r.Context(), v)
		respond(w, created, err)
		return
	}
	writeError(w, 405, "method not allowed")
}
func (s *Server) gateways(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		write(w, 200, s.app.Repo.ListGateways(r.Context()))
		return
	}
	if r.Method == http.MethodPost {
		var v domain.Gateway
		if decode(r, &v) != nil {
			writeError(w, 400, "invalid json")
			return
		}
		created, err := s.app.CreateGateway(r.Context(), v)
		respond(w, created, err)
		return
	}
	writeError(w, 405, "method not allowed")
}
func (s *Server) points(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		write(w, 200, s.app.Repo.ListPoints(r.Context()))
		return
	}
	if r.Method == http.MethodPost {
		var v domain.Point
		if decode(r, &v) != nil {
			writeError(w, 400, "invalid json")
			return
		}
		created, err := s.app.CreatePoint(r.Context(), v)
		respond(w, created, err)
		return
	}
	writeError(w, 405, "method not allowed")
}
func (s *Server) plans(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		write(w, 200, s.app.Repo.ListPlans(r.Context()))
		return
	}
	if r.Method == http.MethodPost {
		var v domain.AcquisitionPlan
		if decode(r, &v) != nil {
			writeError(w, 400, "invalid json")
			return
		}
		created, err := s.app.CreatePlan(r.Context(), v)
		respond(w, created, err)
		return
	}
	writeError(w, 405, "method not allowed")
}
func (s *Server) rules(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		write(w, 200, s.app.Repo.ListRules(r.Context()))
		return
	}
	if r.Method == http.MethodPost {
		var v domain.Rule
		if decode(r, &v) != nil {
			writeError(w, 400, "invalid json")
			return
		}
		created, err := s.app.CreateRule(r.Context(), v)
		respond(w, created, err)
		return
	}
	writeError(w, 405, "method not allowed")
}
func (s *Server) alarms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, 405, "method not allowed")
		return
	}
	write(w, 200, s.app.Repo.ListAlarms(r.Context()))
}
func (s *Server) commands(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		write(w, 200, s.app.Repo.ListCommands(r.Context()))
		return
	}
	if r.Method == http.MethodPost {
		var v domain.Command
		if decode(r, &v) != nil {
			writeError(w, 400, "invalid json")
			return
		}
		created, err := s.app.CreateCommand(r.Context(), v)
		respond(w, created, err)
		return
	}
	if r.Method == http.MethodPatch {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/commands/")
		if path == r.URL.Path {
			writeError(w, 400, "command id required")
			return
		}
		var body struct {
			Actor string `json:"actor"`
		}
		if decode(r, &body) != nil {
			writeError(w, 400, "invalid json")
			return
		}
		updated, err := s.app.ApproveCommand(r.Context(), path, body.Actor)
		respond(w, updated, err)
		return
	}
	writeError(w, 405, "method not allowed")
}
func (s *Server) telemetryRaw(w http.ResponseWriter, r *http.Request) {
	id := domain.ID(r.URL.Query().Get("point_id"))
	if id == "" {
		writeError(w, 400, "point_id required")
		return
	}
	from := parseTime(r.URL.Query().Get("from"), time.Now().Add(-time.Hour))
	to := parseTime(r.URL.Query().Get("to"), time.Now().Add(time.Second))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	write(w, 200, s.app.Telemetry.Query(id, from, to, limit))
}
func (s *Server) telemetryAggregate(w http.ResponseWriter, r *http.Request) {
	id := domain.ID(r.URL.Query().Get("point_id"))
	if id == "" {
		writeError(w, 400, "point_id required")
		return
	}
	from := parseTime(r.URL.Query().Get("from"), time.Now().Add(-time.Hour))
	to := parseTime(r.URL.Query().Get("to"), time.Now().Add(time.Second))
	write(w, 200, s.app.Telemetry.Aggregate(id, from, to))
}
func (s *Server) rollouts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		write(w, 200, s.app.Repo.ListRollouts(r.Context()))
		return
	}
	if r.Method == http.MethodPost {
		var v domain.Rollout
		if decode(r, &v) != nil {
			writeError(w, 400, "invalid json")
			return
		}
		v.ID = s.app.Repo.ID("rollout")
		v.State = "pending"
		v.CreatedAt = time.Now().UTC()
		s.app.Repo.PutRollout(r.Context(), v)
		write(w, 201, v)
		return
	}
	writeError(w, 405, "method not allowed")
}
func parseTime(value string, fallback time.Time) time.Time {
	if value == "" {
		return fallback
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return fallback
	}
	return parsed
}
func decode(r *http.Request, v any) error {
	defer r.Body.Close()
	d := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	d.DisallowUnknownFields()
	return d.Decode(v)
}
func respond(w http.ResponseWriter, v any, err error) {
	if err == nil {
		write(w, http.StatusCreated, v)
		return
	}
	status := http.StatusInternalServerError
	if errors.Is(err, domain.ErrInvalid) {
		status = http.StatusBadRequest
	}
	if errors.Is(err, domain.ErrConflict) {
		status = http.StatusConflict
	}
	if errors.Is(err, domain.ErrNotFound) {
		status = http.StatusNotFound
	}
	writeError(w, status, err.Error())
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": v, "request_id": time.Now().UTC().Format("20060102T150405.000000000Z")})
}
func writeError(w http.ResponseWriter, status int, message string) {
	write(w, status, map[string]string{"code": strconv.Itoa(status), "message": message})
}

var _ = context.Background
