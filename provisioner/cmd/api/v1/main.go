package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/BradleyLewis08/HiVE/internal/ingress"
	k8sclient "github.com/BradleyLewis08/HiVE/internal/kubernetes"
	k8sProvisioner "github.com/BradleyLewis08/HiVE/internal/provisioner"
	utils "github.com/BradleyLewis08/HiVE/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Server struct {
	k8sProvisioner *k8sProvisioner.Provisioner
	ingressManager *ingress.IngressManager
}

func NewServer() (*Server, error) {
	client, clientInitErr := k8sclient.GetKubernetesClient()
	provisioner := k8sProvisioner.NewProvisioner(client)
	ingressManager := ingress.NewIngressManager(client)
	if clientInitErr != nil {
		return nil, clientInitErr
	}

	return &Server{k8sProvisioner: provisioner, ingressManager: ingressManager}, nil
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

	server.ingressManager.ProvisionIngressController()

	r := chi.NewRouter()

	// Define routes

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Post("/environment", func(w http.ResponseWriter, r *http.Request) {
		server.createEnvironment(w, r) 
	})

	r.Post("/environment/delete", func(w http.ResponseWriter, r *http.Request) {
		server.deleteEnvironment(w, r)
	})

	log.Println("Starting server on :8000")

	err = http.ListenAndServe(":8000", r)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}


type EnvironmentProvisionRequest struct {
	CourseName string `json:"courseName"`
	AssignmentName string `json:"assignmentName"`
	NetIDs   []string `json:"netIDs"`
	Image   string   `json:"image"`
}

/* Creates an envvironment for this particular assignment and course, 
*  for each student in the request
*/
func (s *Server) createEnvironment(w http.ResponseWriter, r *http.Request) {
	var envReq EnvironmentProvisionRequest

	if err := json.NewDecoder(r.Body).Decode(&envReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	courseName := utils.LowerCaseAndStrip(envReq.CourseName)
	assignmentName := utils.LowerCaseAndStrip(envReq.AssignmentName)

	// Provision environment for each student (NetID)
	for _, netID := range envReq.NetIDs {
		err := s.k8sProvisioner.ProvisionStudentEnvironment(
			assignmentName,
			courseName,
			envReq.Image,
			netID,
		)
		if err != nil {
			http.Error(w, "Failed to create environment", http.StatusInternalServerError)
			return
		}

		// Add route to ingress controller
		err = s.ingressManager.AddRouteToIngress(
			assignmentName,
			courseName,
			netID,
		)

		if err != nil {
			http.Error(w, "Failed to add route to ingress controller", http.StatusInternalServerError)
			return
		}
	}

	fmt.Printf("Created environments for all netIDs\n")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

type EnvironmentDeleteRequest struct {
	AssignmentName string `json:"assignmentName"`
	CourseName string `json:"courseName"`
	NetID  string `json:"netIDs"`
}

func (s* Server) deleteEnvironment(w http.ResponseWriter, r* http.Request) {
	var envDeleteReq EnvironmentDeleteRequest

	if err := json.NewDecoder(r.Body).Decode(&envDeleteReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	courseName := utils.LowerCaseAndStrip(envDeleteReq.CourseName)
	assignmentName := utils.LowerCaseAndStrip(envDeleteReq.AssignmentName)

	err := s.k8sProvisioner.DeleteEnvironment(assignmentName, courseName, envDeleteReq.NetID)

	if err != nil {
		http.Error(w, "Failed to delete environment", http.StatusInternalServerError)
		return
	}

	// TODO: Remove route from ingress controller
	fmt.Println("Deleted environments for netID: ", envDeleteReq.NetID)
	w.WriteHeader(http.StatusNoContent)
}

