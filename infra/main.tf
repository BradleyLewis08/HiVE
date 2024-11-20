provider "aws" {
  region = var.region
}

data "aws_availability_zones" "available" {
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

locals {
  cluster_name    = "hive-${var.environment}-${random_string.suffix.result}"
  tags = {
    Environment = var.environment
    Project     = "HiVE"
    ManagedBy   = "Terraform"
    CostCenter  = var.cost_center
  }
}

resource "random_string" "suffix" {
  length  = 6
  special = false
  upper   = false
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.8.1"

  name = "hive-${var.environment}-vpc"
  cidr = var.vpc_cidr

  azs             = slice(data.aws_availability_zones.available.names, 0, var.az_count)
  private_subnets = [for i in range(var.az_count) : cidrsubnet(var.vpc_cidr, 4, i)]
  public_subnets  = [for i in range(var.az_count) : cidrsubnet(var.vpc_cidr, 4, i + var.az_count)]

  enable_nat_gateway     = true
  single_nat_gateway     = var.environment != "production" # Use single NAT gateway for non-prod
  enable_dns_hostnames   = true
  enable_dns_support     = true
  
  # Add flow logs for network monitoring
  enable_flow_log                      = true
  create_flow_log_cloudwatch_log_group = true
  create_flow_log_cloudwatch_iam_role  = true

  public_subnet_tags = {
    "kubernetes.io/role/elb"                      = 1
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb"             = 1
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
  }

  tags = local.tags
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "20.8.5"

  cluster_name                = local.cluster_name
  cluster_version            = var.kubernetes_version
  cluster_enabled_log_types  = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  # Enable OIDC provider for service accounts
  enable_irsa = true
  
  # Enable cluster encryption
  cluster_encryption_config = {
    provider_key_arn = aws_kms_key.eks.arn
    resources        = ["secrets"]
  }

  # Enable public access to the API server
  cluster_endpoint_public_access = true
  
  # Optionally restrict access to specific IPs
  cluster_endpoint_public_access_cidrs = ["0.0.0.0/0"]  # Be more restrictive in production

  # Add necessary security group rules
  node_security_group_additional_rules = {
    ingress_self_all = {
      description = "Node to node all ports/protocols"
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      type        = "ingress"
      self        = true
    }
    egress_all = {
      description = "Node all egress"
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      type        = "egress"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  # Enable necessary networking add-ons
  cluster_addons = {
    vpc-cni = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    coredns = {
      most_recent = true
    }
    aws-ebs-csi-driver = {
      service_account_role_arn = module.irsa-ebs-csi.iam_role_arn
    }
  }

  # Node groups configuration
  eks_managed_node_groups = {
    system = {
      name           = "system-ng"
      instance_types = var.system_instance_types
      capacity_type  = "ON_DEMAND"

      min_size     = var.system_node_min
      max_size     = var.system_node_max
      desired_size = var.system_node_desired

      # Enable node group autoscaling
      enable_monitoring = true
      
      labels = {
        NodeGroupType = "system"
      }
      
      # Add taints to ensure system workloads only
      taints = [
        {
          key    = "dedicated"
          value  = "system"
          effect = "NO_SCHEDULE"
        }
      ]
    }

    workload = {
      name           = "workload-ng"
      instance_types = var.workload_instance_types
      capacity_type  = "SPOT" # Use spot instances for student workloads
      
      min_size     = var.workload_node_min
      max_size     = var.workload_node_max
      desired_size = var.workload_node_desired

      labels = {
        NodeGroupType = "workload"
      }
    }
  }

  tags = local.tags
}

# Create KMS key for cluster encryption
resource "aws_kms_key" "eks" {
  description             = "EKS Cluster Encryption Key"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = local.tags
}

# EBS CSI driver configuration
module "irsa-ebs-csi" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version = "5.39.0"

  create_role                   = true
  role_name                     = "AmazonEKSEBSCSIRole-${local.cluster_name}"
  provider_url                  = module.eks.oidc_provider
  role_policy_arns             = [data.aws_iam_policy.ebs_csi_policy.arn]
  oidc_fully_qualified_subjects = ["system:serviceaccount:kube-system:ebs-csi-controller-sa"]
}

data "aws_iam_policy" "ebs_csi_policy" {
  arn = "arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
}

# Add VPC endpoints for EKS
module "vpc_endpoints" {
  source  = "terraform-aws-modules/vpc/aws//modules/vpc-endpoints"
  version = "5.0.0"

  vpc_id             = module.vpc.vpc_id
  security_group_ids = [module.eks.cluster_security_group_id]

  endpoints = {
    s3 = {
      service         = "s3"
      service_type    = "Gateway"
      route_table_ids = module.vpc.private_route_table_ids
    }
    ec2 = {
      service             = "ec2"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
    }
    ecr_api = {
      service             = "ecr.api"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
    }
    ecr_dkr = {
      service             = "ecr.dkr"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
    }
    sts = {
      service             = "sts"
      private_dns_enabled = true
      subnet_ids          = module.vpc.private_subnets
    }
  }

  tags = local.tags
}