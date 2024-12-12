# Provisioning service

The provisioner is a Golang HTTP service that is meant to be deployed to the same k8s cluster
that HiVE is hosted on. It interacts with k8s to provision and manage development environments.

### Repo structure:

    .
    ├── cmd            			# Houses main Golang server logic
    ├── frontend                # Deployment templates for environments
    ├── infra             		# Terraform config files for cluster
    ├── internal             	# Internal business logic
    ├── scripts             	# Useful scripts for starting server
    ├── services             	# Service templates
    └── README.md

To start an instance of the provisioner, ensure you have a EKS cluster running on AWS (typically through terraform - see `infra` for the template used for the senior thesis artifact).

Then, you may run `bash scripts/start.sh` which sets up the provisoiner on your cluster.
