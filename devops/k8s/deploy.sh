#!/bin/bash
set -x
export AWS_PROFILE=personal
export $(sops -d .env | xargs)

envsubst < ./backend-secrets.yaml.tpl > ./backend-secrets.yaml

kubectl delete deployment cmg-backend-deployment 
kubectl delete deployment cmg-frontend-deployment 

kubectl apply -f ./backend-secrets.yaml
kubectl apply -f ./backend-deployment.yaml
kubectl apply -f ./backend-lb.yaml

# get the loadbalancer hostname because it is dynamically assigned and not availabe via k8s service discovery 
export BACKEND_URL=$(kubectl get svc cmg-backend-service -o=json | jq '.status.loadBalancer.ingress[0].hostname' | tr -d '"')
envsubst < ./frontend-deployment.yaml.tpl > ./frontend-deployment.yaml
cat ./frontend-deployment.yaml

kubectl apply -f ./frontend-deployment.yaml
kubectl apply -f ./frontend-lb.yaml
