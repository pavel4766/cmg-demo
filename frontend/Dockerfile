FROM node:22 as builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
EXPOSE 80

# The last line is not strictly necessary as it's the default command for the nginx image
CMD ["nginx", "-g", "daemon off;"]
