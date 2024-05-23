apiVersion: apps/v1
kind: Deployment
metadata:
  name: cmg-frontend-deployment
  labels:
    app: cmg-frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: cmg-frontend
  template:
    metadata:
      labels:
        app: cmg-frontend
    spec:
      containers:
      - name: cmg-frontend 
        image: 730335406037.dkr.ecr.us-east-2.amazonaws.com/cmg-frontend:latest
        ports:
        - containerPort: 80
        command: ["sh", "-c"]
        args:
        - |
          sed -i "s/__bur__/$BACKEND_URL/g" /usr/share/nginx/html/env.js;
          nginx -g 'daemon off;'


