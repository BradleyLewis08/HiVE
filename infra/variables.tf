variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-2"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "development"
}

variable "cost_center" {
  description = "Cost center for billing"
  type        = string
  default     = "education"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "az_count" {
  description = "Number of AZs to use"
  type        = number
  default     = 2
}

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.29"
}

variable "system_instance_types" {
  description = "Instance types for system node group"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "workload_instance_types" {
  description = "Instance types for workload node group"
  type        = list(string)
  default     = ["t3.small", "t3.medium"]
}

variable "system_node_min" {
  description = "Minimum number of system nodes"
  type        = number
  default     = 1
}

variable "system_node_max" {
  description = "Maximum number of system nodes"
  type        = number
  default     = 2
}

variable "system_node_desired" {
  description = "Desired number of system nodes"
  type        = number
  default     = 1
}

variable "workload_node_min" {
  description = "Minimum number of workload nodes"
  type        = number
  default     = 1
}

variable "workload_node_max" {
  description = "Maximum number of workload nodes"
  type        = number
  default     = 10
}

variable "workload_node_desired" {
  description = "Desired number of workload nodes"
  type        = number
  default     = 2
}