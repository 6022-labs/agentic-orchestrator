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

COPY nginx.conf /etc/nginx/templates/default.conf.template
COPY --chmod=755 10-gateway-base-url.envsh /docker-entrypoint.d/
# Fills the resolver from resolv.conf and limits envsubst to the template's own variables.
ENV NGINX_ENTRYPOINT_LOCAL_RESOLVERS=1 \
    NGINX_ENVSUBST_FILTER=^(GATEWAY__BASE_URL|NGINX_LOCAL_RESOLVERS)$
COPY --from=build /app/dist /var/www/html

EXPOSE 80
