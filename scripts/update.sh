#!/bin/bash

services=(cs101-ak2649-lb cs101-bel25-lb cs101-tlb76-lb)  # Add your service names here

for svc in "${services[@]}"
do
  kubectl patch svc $svc -p '{"spec": {"type": "ClusterIP"}}'
  echo "Updated $svc to ClusterIP"
done
