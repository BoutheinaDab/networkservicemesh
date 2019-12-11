#!/bin/bash

kubectl="kubectl -n ${NSM_NAMESPACE}"

## First argument is the number of NS

PodNumber=0
echo "===== >>>>> get the number of requests requesting the same NS  <<<<< ==========="

for nsc in $(${kubectl} get pods -o=name | grep -E "vpn-gateway-nsc" | sed 's@.*/@@'); do
    
    echo "===== >>>>> PROCESSING ${nsc}  <<<<< ==========="
    PodNumber=$((PodNumber + 1))
done
  
echo "pods number is ",${PodNumber}

return ${PodNumber}
