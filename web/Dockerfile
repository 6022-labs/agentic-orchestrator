FROM node:22-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

ARG VITE_APP_KIT_PROJECT_ID
ARG VITE_ENVIRONMENT
ARG VITE_IPFS_GATEWAY

ENV VITE_APP_KIT_PROJECT_ID=$VITE_APP_KIT_PROJECT_ID
ENV VITE_ENVIRONMENT=$VITE_ENVIRONMENT
ENV VITE_IPFS_GATEWAY=$VITE_IPFS_GATEWAY

RUN npm run build

FROM nginx:1.29-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /var/www/html

EXPOSE 80
