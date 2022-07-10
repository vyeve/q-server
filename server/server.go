package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"github.com/vyeve/q-server/models"
	"github.com/vyeve/q-server/repository"
	"github.com/vyeve/q-server/utils/logger"
	"github.com/vyeve/q-server/utils/validator"

	"go.uber.org/fx"
)

type Server interface {
	Init()
}

type serverImpl struct {
	logger    logger.Logger
	validator validator.ValidatorJSON
	repo      repository.Repository
	server    *http.Server
	limit     chan struct{}
}

func New(params Params) Server {
	srv := &serverImpl{
		logger:    params.Logger,
		validator: params.Validator,
		repo:      params.Repo,
	}
	port, found := syscall.Getenv(EnvServerPort)
	if !found {
		port = strconv.Itoa(defaultPort)
	}
	srv.server = &http.Server{
		Addr:    ":" + port,
		Handler: srv,
	}
	limit, err := strconv.Atoi(os.Getenv(EnvRequestsLimit))
	if err != nil {
		limit = defaultRequestLimit
	}
	srv.limit = make(chan struct{}, limit)
	params.LifeCycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				srv.logger.Infof("start to listen on port %s", port)
				go func() {
					if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						srv.logger.Fatalf("start server err: %v", err)
					}
				}()
				return nil
			},
			OnStop: func(context.Context) error {
				srv.logger.Infof("server stopped.")
				defer close(srv.limit) // need to close channel after stop HTTP server
				return srv.server.Close()
			},
		},
	)
	return srv
}

// Init method is needed to invoke server on start
func (s *serverImpl) Init() {}

// ServeHTTP implements http.Handler interface
func (s *serverImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Debugf("request %q. method: %s", r.URL.Path, r.Method)
	s.limit <- struct{}{}
	defer func() {
		<-s.limit
	}()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	switch r.URL.Path {
	case transferEndpoint:
		s.handleUploadTransfers(w, r)
	case uploadEndpoint:
		s.handleUploadFile(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// handleUploadTransfers handles /transfers endpoint
func (s *serverImpl) handleUploadTransfers(w http.ResponseWriter, r *http.Request) {
	s.uploadTransfers(r.Context(), w, r.Body)
}

// handleUploadFile handles /upload endpoint
func (s *serverImpl) handleUploadFile(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile(fileKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.uploadTransfers(r.Context(), w, f)
}

// uploadTransfers common method to upload transfers
func (s *serverImpl) uploadTransfers(ctx context.Context, w http.ResponseWriter, body io.ReadCloser) {
	defer body.Close() // nolint: errcheck
	p, err := ioutil.ReadAll(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = s.validator.Validate(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	receipt := new(models.Receipt)
	err = json.Unmarshal(p, receipt)
	switch err {
	case nil:
	case models.ErrIncorrectAmount, models.ErrUnsupportedFiat:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.repo.UploadTransfers(ctx, receipt)
	switch err {
	case nil:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "successfully uploaded transfers") // nolint: errcheck
	case repository.ErrInsufficientFunds, repository.ErrUnknownOrganization:
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	case repository.ErrNoTransfers:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
