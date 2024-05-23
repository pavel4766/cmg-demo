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
        resources:
          limits:
            cpu: 500m
          requests:
            cpu: 200m  
        command: ["sh", "-c"]
        args:
        - |
          sed -i "s/__bur__/$BACKEND_URL/g" /usr/share/nginx/html/env.js;
          nginx -g 'daemon off;'

apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: cmg-backend-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: cmg-backend-deployment
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 60

