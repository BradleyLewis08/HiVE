#!/bin/bash

# List of netids
netids=("bel25" "ak2649" "tlb76")

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

kubectl delete configmap nginx-config-cs101
kubectl delete deployment nginx-reverse-proxy-cs101
kubectl delete service nginx-reverse-proxy-cs101

echo "All operations completed."
