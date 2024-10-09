#!/bin/bash

# List of netids
netids=("bel25" "ak276" "tlb76")

# Loop through each netid
for netid in "${netids[@]}"
do
    echo "Processing netid: $netid"
    
    # Delete the deployment
    kubectl delete deployment "hive-course-cs101-$netid"
    
    # Delete the service
    kubectl delete service "cs101-$netid-lb"
    
    echo "Completed operations for $netid"
    echo "------------------------"
done

kubectl delete configmap nginx-config
kubectl delete deployment nginx-reverse-proxy
kubectl delete service nginx-reverse-proxy-service

echo "All operations completed."
