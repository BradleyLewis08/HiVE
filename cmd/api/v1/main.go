package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	k8sclient "github.com/BradleyLewis08/HiVE/internal/kubernetes"
	k8sProvisioner "github.com/BradleyLewis08/HiVE/internal/provisioner"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Server struct {
	k8sProvisioner *k8sProvisioner.Provisioner
}

func NewServer() (*Server, error) {
	client, clientInitErr := k8sclient.GetKubernetesClient()
	provisioner := k8sProvisioner.NewProvisioner(client)
	if clientInitErr != nil {
		return nil, clientInitErr
	}

	return &Server{k8sProvisioner: provisioner}, nil
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	server, err := NewServer()

	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Post("/environment", func(w http.ResponseWriter, r *http.Request) {
		server.createEnvironment(w, r) 
	})

	log.Println("Starting server on :8000")

	err = http.ListenAndServe(":8000", r)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

type EnvironmentRequest struct {
	CourseName string `json:"courseName"`
	NetIDs   []string `json:"netIDs"`
	Image   string   `json:"image"`
}

func transformCourseName(courseName string) string {
	lowerCase := strings.ToLower(courseName)
	transformed := strings.ReplaceAll(lowerCase, " ", "-")
	return transformed
}


func (s *Server) createEnvironment(w http.ResponseWriter, r *http.Request) {
	var envReq EnvironmentRequest

	if err := json.NewDecoder(r.Body).Decode(&envReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	envReq.CourseName = transformCourseName(envReq.CourseName)

	// Provision environment for each student (NetID)
	for _, netID := range envReq.NetIDs {
		fmt.Println("Creating environment for netID: ", netID)
		err := s.k8sProvisioner.ProvisionStudentEnvironment(
			envReq.CourseName,
			envReq.Image,
			netID,
		)
		if err != nil {
			http.Error(w, "Failed to create environment", http.StatusInternalServerError)
			return
		}
	}

	fmt.Printf("Created environments for all netIDs\n")

	// Create the router for the course
	courseServiceAddress, err := s.k8sProvisioner.ProvisionCourseRouter(envReq.CourseName, envReq.NetIDs)

	if err != nil {
		http.Error(w, "Failed to create router", http.StatusInternalServerError)
		return
	}

	response := struct {
		RouterAddress string `json:"routerAddress"` 	
	} {
		RouterAddress: courseServiceAddress,
	}

	fmt.Printf("Router created successfully at: %s\n", courseServiceAddress)

	// Get the IP of the NGINX service
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response);
}

