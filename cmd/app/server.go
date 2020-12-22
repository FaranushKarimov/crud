package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/FaranushKarimov/crud/pkg/customers"
	"github.com/FaranushKarimov/crud/pkg/security"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Server ...
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	securitySvc  *security.Service
}

// NewServer ...
func NewServer(mux *mux.Router, customersSvc *customers.Service, securitySvc *security.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc, securitySvc: securitySvc}
}

// ServeHTTP ...
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

const (
	// GET ...
	GET = "GET"
	// POST ...
	POST = "POST"
	// DELETE ...
	DELETE = "DELETE"
)

// Init ...
func (s *Server) Init() {
	// s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)

	// s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomers)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)

	// s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)

	// s.mux.HandleFunc("/customers.save", s.handleSaveCustomer)
	// s.mux.HandleFunc("/customers", s.handleSaveCustomer).Methods(POST)

	// s.mux.HandleFunc("/customers.removeById", s.handleRemoveByID)
	s.mux.HandleFunc("/customers/{id}", s.handleRemoveByID).Methods(DELETE)

	// s.mux.HandleFunc("/customers.blockById", s.handleBlockByID)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)

	// s.mux.HandleFunc("/customers.unblockById", s.handleUnblockByID)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnblockByID).Methods(DELETE)

	// s.mux.Use(middleware.Basic(s.securitySvc.Auth))

	s.mux.HandleFunc("/api/customers", s.handleSaveCustomer).Methods(POST)
	s.mux.HandleFunc("/api/customers/token", s.handleCreateToken).Methods(POST)
	s.mux.HandleFunc("/api/customers/token/validate", s.handleValidateToken).Methods(POST)
}

func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.ByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	item, err := s.customersSvc.All(r.Context())
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {

	item, err := s.customersSvc.AllActive(r.Context())
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) handleSaveCustomer(w http.ResponseWriter, r *http.Request) {
	var customer *customers.Customer

	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	customer.Password = string(hash)

	item, err := s.customersSvc.Save(r.Context(), customer)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleRemoveByID(w http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		return
	}

	item, err := s.customersSvc.RemoveByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		return
	}

	item, err := s.customersSvc.BlockByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) handleUnblockByID(w http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		return
	}

	item, err := s.customersSvc.UnblockByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) handleCreateToken(w http.ResponseWriter, r *http.Request) {
	var customer *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, err := s.securitySvc.TokenForCustomer(r.Context(), customer.Login, customer.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return

	}

	response := map[string]interface{}{"status": "ok", "token": token}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	var customer *struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := s.securitySvc.AuthenticateCustomer(r.Context(), customer.Token)
	if err != nil {
		// status := http.StatusInternalServerError
		status := http.StatusNotFound
		errText := "not found"

		if err == security.ErrNoSuchUser {
			status = http.StatusNotFound
			errText = "not found"
		}
		if err == security.ErrExpireToken {
			status = http.StatusBadRequest
			errText = "expired"
		}

		response := map[string]interface{}{"status": "fail", "reason": errText}
		data, err := json.Marshal(response)
		if err != nil {
			http.Error(w, http.StatusText(status), status)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		// w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		if err != nil {
			log.Print(err)
			return
		}
		return
	}

	response := make(map[string]interface{})
	response["status"] = "ok"
	response["customerId"] = id

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}

}
