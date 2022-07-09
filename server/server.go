package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	params.LifeCycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				srv.logger.Infof("start to listen on port %s", port)
				go func() {
					if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						srv.logger.Fatalf("start server err: %v", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.server.Close()
			},
		},
	)
	return srv
}

func (s *serverImpl) Init() {}

func (s *serverImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	switch r.URL.Path {
	case transferEndpoint:
		s.handleUploadTransfers(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *serverImpl) handleUploadTransfers(w http.ResponseWriter, r *http.Request) {
	p, err := ioutil.ReadAll(r.Body)
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
	err = s.repo.UploadTransfers(r.Context(), receipt)
	switch err {
	case nil:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "successfully uploaded transfers")
	case repository.ErrInsufficientFunds, repository.ErrUnknownOrganization:
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	case repository.ErrNoTransfers:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}