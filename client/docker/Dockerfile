FROM --platform=$BUILDPLATFORM node:alpine3.19 AS build

WORKDIR /app

COPY . .

RUN corepack enable
RUN yarn install --frozen-lockfile
RUN yarn workspace @vertex-center/client build

FROM nginx:alpine3.18-slim

COPY --from=build /app/client/dist /usr/share/nginx/html
COPY client/docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY client/docker/entrypoint.sh /

RUN ["chmod", "+x", "/entrypoint.sh"]
ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
