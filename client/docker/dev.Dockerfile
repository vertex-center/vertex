FROM node:alpine3.19

WORKDIR /app

COPY package.json yarn.lock tsconfig.base.json ./

# Because of the workspace, we need to copy all package.json files to keep the same lockfile.
COPY client/package.json ./client/
COPY docs/package.json ./docs/
COPY packages/components/package.json ./packages/components/

RUN corepack enable
RUN yarn install --frozen-lockfile

EXPOSE 5173

CMD ["yarn", "workspace", "@vertex-center/client", "dev", "--host"]
