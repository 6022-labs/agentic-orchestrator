FROM node:26-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

ARG VITE_APP_KIT_PROJECT_ID
ARG VITE_ENVIRONMENT
ARG VITE_GATEWAY_API_BASE_URL

ENV VITE_APP_KIT_PROJECT_ID=$VITE_APP_KIT_PROJECT_ID
ENV VITE_ENVIRONMENT=$VITE_ENVIRONMENT
ENV VITE_GATEWAY_API_BASE_URL=$VITE_GATEWAY_API_BASE_URL

RUN npm run build

FROM nginx:1.31-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /var/www/html

EXPOSE 80
