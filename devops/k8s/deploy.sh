#!/bin/bash
set -x
export AWS_PROFILE=personal
export $(sops -d .env | xargs)

envsubst < ./backend-secrets.yaml.tpl > ./backend-secrets.yaml

kubectl delete deployment cmg-backend-deployment 
kubectl apply -f ./backend-secrets.yaml
kubectl apply -f ./cmg-backend-deployment.yaml
kubectl apply -f ./cmg-backend-lb.yaml
