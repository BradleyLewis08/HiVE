package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	k8sclient "github.com/BradleyLewis08/HiVE/internal/kubernetes"
	k8sProvisioner "github.com/BradleyLewis08/HiVE/internal/provisioner"
	"github.com/gorilla/mux"
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
	provisioner := k8sProvisioner.NewProvisioner(client)
	if clientInitErr != nil {
		return nil, clientInitErr
	}

	return &Server{k8sClient: client, k8sProvisioner: provisioner}, nil
}


func main() {
	server, err := NewServer()

	if err != nil {
		panic(err)
	}
	BASE_URL := "/api/v1"
	r := mux.NewRouter()

	// Add the route for the server
	r.HandleFunc(BASE_URL+"/server", GetServer).Methods("GET")
	r.HandleFunc(BASE_URL+"/environment", server.createEnvironment).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func (s *Server) createEnvironment(w http.ResponseWriter, r *http.Request) {
	err := s.k8sProvisioner.CreateEnvironment();
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetServer(w http.ResponseWriter, r *http.Request) {
	env := Environment{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(env)
}
