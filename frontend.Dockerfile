# pull official base image
FROM node:current-alpine3.18
ARG GH_ACTOR=nduyphuong
ARG GH_REPO=gorya
LABEL org.opencontainers.image.source https://github.com/${GH_ACTOR}/${GH_REPO}
LABEL org.opencontainers.image.licenses MIT

ENV NODE_OPTIONS=--openssl-legacy-provider
# set working directory
WORKDIR /app

# install app dependencies
#copies package.json and package-lock.json to Docker environment
COPY client/package.json ./
RUN npm install

# Installs all node packages

# Copies everything over to Docker environment
COPY ./client ./

CMD yarn start
