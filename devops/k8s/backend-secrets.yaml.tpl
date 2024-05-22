apiVersion: v1
kind: Secret
metadata:
  name: cmg-secret
type: Opaque
stringData:
  DB_USER: $DB_USER
  DB_PASSWORD: $DB_PASSWORD
  DB_HOST: $DB_HOST
    
