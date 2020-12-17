package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/FaranushKarimov/crud/cmd/app/middleware"
	"github.com/FaranushKarimov/crud/pkg/customers"
	"github.com/FaranushKarimov/crud/pkg/customers/security"
	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	security     *security.Service
}

// NewServer ...
func NewServer(mux *mux.Router, customersSvc *customers.Service, security *security.Service) *Server {
	return &Server{
		mux:          mux,
		customersSvc: customersSvc,
		security:     security,
	}
}

// ServerHTTP ...
func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
}

// Init ...
func (s *Server) Init() {
	s.mux.Use(middleware.Logger)
	auth := middleware.Basic(s.security.Auth)
	s.mux.Use(auth)
	chMd := middleware.CheckHeader("Content-Type", "application/json")

	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods("GET")
	s.mux.Handle("/customers", chMd(http.HandlerFunc(s.handleSaveCustomer))).Methods("POST")
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods("GET")
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods("GET")
	s.mux.HandleFunc("/customers/{id}", s.handleRemoveCustomerByID).Methods("DELETE")
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockCustomerByID).Methods("POST")
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnblockCustomerByID).Methods("DELETE")
}

func (s *Server) handleGetCustomerByID(rw http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.ByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		if errors.Is(err, customers.ErrNotFound) {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetAllCustomers(rw http.ResponseWriter, r *http.Request) {
	items, err := s.customersSvc.All(r.Context())
	if err != nil {
		log.Print(err)
		if errors.Is(err, customers.ErrNotFound) {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetAllActiveCustomers(rw http.ResponseWriter, r *http.Request) {
	items, err := s.customersSvc.AllActive(r.Context())
	if err != nil {
		log.Print(err)
		if errors.Is(err, customers.ErrNotFound) {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleSaveCustomer(rw http.ResponseWriter, r *http.Request) {
	var itemToSave *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&itemToSave); err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.Save(r.Context(), itemToSave)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleRemoveCustomerByID(rw http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.RemoveByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		if errors.Is(err, customers.ErrNotFound) {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleBlockCustomerByID(rw http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.BlockByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		if errors.Is(err, customers.ErrNotFound) {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleUnblockCustomerByID(rw http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.UnblockByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		if errors.Is(err, customers.ErrNotFound) {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(data)
	if err != nil {
		log.Print(err)
	}
}
