FROM --platform=$BUILDPLATFORM node:alpine3.19 AS build

WORKDIR /app

COPY ../packages.json .
COPY ../yarn.lock .
COPY ../client ./client
COPY ../packages ./packages

RUN yarn install --frozen-lockfile
RUN yarn workspace @vertex-center/client build

FROM nginx:alpine3.18-slim

COPY --from=build /app/dist /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY docker/entrypoint.sh /

RUN ["chmod", "+x", "/entrypoint.sh"]
ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
