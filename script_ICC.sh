!/usr/bin/env bash
echo ""
echo "********************************************************************************"
echo "          Welcome to service chaining with NetworkServiceMesh                   "
echo "********************************************************************************"
echo ""

make helm-delete-vpn

make helm-delete-nsm-monitoring

make helm-delete-nsm

echo "******************go generate**********************"
go generate ./...
echo "*****************go build**************************"
go build ./...


make k8s-nsmgr-save

sleep 2m
echo "******** installation helm ******"
make helm-install-nsm

echo "******** installing nsm-monitoring *****"

 make helm-install-nsm-monitoring


echo "****** installing vpn*********************"

 make helm-install-vpn



echo "*********checking************* "

#make k8s-check





