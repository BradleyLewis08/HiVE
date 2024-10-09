package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/BradleyLewis08/HiVE/internal/imager"
	k8sclient "github.com/BradleyLewis08/HiVE/internal/kubernetes"
	k8sProvisioner "github.com/BradleyLewis08/HiVE/internal/provisioner"
	"github.com/BradleyLewis08/HiVE/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	k8sClient *kubernetes.Clientset
	k8sProvisioner *k8sProvisioner.Provisioner
}

type Environment struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

func NewServer() (*Server, error) {
	client, clientInitErr := k8sclient.GetKubernetesClient()
	imager := imager.NewImager()
	provisioner := k8sProvisioner.NewProvisioner(client, imager)
	if clientInitErr != nil {
		return nil, clientInitErr
	}

	return &Server{k8sClient: client, k8sProvisioner: provisioner}, nil
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

	log.Println("Starting server on :8080")

	http.ListenAndServe(":8080", r)
}

type EnvironmentRequest struct {
	CourseName string `json:"courseName"`
	Image string `json:"image"`
	Capacity  int    `json:"capacity"`
	NetIDs   []string `json:"netIDs"`
}


func (s *Server) createEnvironment(w http.ResponseWriter, r *http.Request) {
	var envReq EnvironmentRequest

	if err := json.NewDecoder(r.Body).Decode(&envReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Create environment for each netID
	for _, netID := range envReq.NetIDs {
		fmt.Println("Creating environment for netID: ", netID)
		err := s.k8sProvisioner.CreateEnvironment(
			envReq.Capacity,
			envReq.CourseName,
			envReq.Image,
			netID,
		)
		if err != nil {
			http.Error(w, "Failed to create environment", http.StatusInternalServerError)
			return
		}
	}

	// Create NGINX reverse proxy
	routes := utils.ConstructReverseProxyRoutes(
		envReq.NetIDs,
		envReq.CourseName,
	)

	if err := s.k8sProvisioner.CreateReverseProxy(routes); err != nil {
		http.Error(w, "Failed to create reverse proxy", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Environment created successfully\n")

	// Get the IP of the NGINX service
	w.WriteHeader(http.StatusCreated)
}

