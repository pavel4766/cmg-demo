#!/bin/bash
set -x
export AWS_PROFILE=personal
export $(sops -d .env | xargs)

envsubst < ./backend-secrets.yaml.tpl > ./backend-secrets.yaml

kubectl delete deployment backend-deployment 
kubectl apply -f ./backend-secrets.yaml
kubectl apply -f ./backend-deployment.yaml
kubectl apply -f ./backend-lb.yaml

kubectl apply -f ./frontend-deployment.yaml
kubectl apply -f ./frontend-lb.yaml
