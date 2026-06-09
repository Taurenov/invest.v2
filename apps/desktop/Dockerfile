FROM node:22-alpine
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
EXPOSE 1420
# Браузер на Windows ходит на API через localhost:8080 (проброшен из контейнера api)
ENV VITE_API_URL=http://127.0.0.1:8080
ENV VITE_API_TOKEN=dev-token
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "1420"]
