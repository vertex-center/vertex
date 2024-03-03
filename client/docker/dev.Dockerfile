FROM node:alpine3.19

WORKDIR /app

COPY . .

RUN corepack enable
RUN yarn install --frozen-lockfile

EXPOSE 5173

CMD ["yarn", "workspace", "@vertex-center/client", "dev", "--host"]
