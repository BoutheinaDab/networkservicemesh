#!/bin/bash

kubectl="kubectl -n ${NSM_NAMESPACE}"

## First argument is the number of NS

## Clean old traces files

rm "Traces/Iperf/$3/$2_NS_$1NSE/*"

rm "Traces/CPU_Stats/$1NSE_$2NS/*"

#  iper all the things!
EXIT_VAL=0
i=0
echo "===== >>>>> PROCESSING iperf between pods in NS  <<<<< ==========="
for nsc in $(${kubectl} get pods -o=name | grep -E "vpn-gateway-nsc" | sed 's@.*/@@'); do
    echo "===== >>>>> PROCESSING ${nsc}  <<<<< ==========="
    #if [[ ${nsc} == vppagent-* ]]; then

        LogvpnFileName="vpn-gateway-nsc_client_vpn_$i.txt"
        ${kubectl} logs "${nsc}" -c iperf-vpn-client >  "Traces/Iperf/$3/$2_NS_$1NSE/${LogvpnFileName}"
        LogFirewallFileName="vpn-gateway-nsc_client_firewall_$i.txt"
        ${kubectl} logs "${nsc}" -c iperf-vpn-firewall-client >  "Traces/Iperf/$3/$2_NS_$1NSE/${LogFirewallFileName}"
        i=$((i + 1))
done
j=0
for nse in $(${kubectl} get pods -o=name | grep -E "vpn-gateway-nse" | sed 's@.*/@@'); do
    echo "===== >>>>> PROCESSING ${nse}  <<<<< ==========="
 
        LogvpnFileName="vpn-gateway-nse_server_$j.txt"
        ${kubectl} logs "${nse}" -c iperf-vpn-server >  "Traces/Iperf/$3/$2_NS_$1NSE/${LogvpnFileName}"
        j=$((j + 1))
done

k=0
for nse in $(${kubectl} get pods -o=name | grep -E "vppagent-firewall" | sed 's@.*/@@'); do
    echo "===== >>>>> PROCESSING ${nse}  <<<<< ==========="
 
    LogFirewallFileName="vppagent-firewall_server_$k.txt"
    ${kubectl} logs "${nse}" -c iperf-firewall-server >  "Traces/Iperf/$3/$2_NS_$1NSE/${LogFirewallFileName}"
    LogClientFileName="vppagent-firewall_client_$k.txt"
    ${kubectl} logs "${nse}" -c iperf-firewall-client >  "Traces/Iperf/$3/$2_NS_$1NSE/${LogClientFileName}"

    k=$((k + 1))
done

exit ${EXIT_VAL}
