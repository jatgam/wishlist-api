FROM node:10.17.0-alpine as builder
ENV NODE_ENV development
WORKDIR /app
ADD package.json package-lock.json ./
RUN npm ci
ADD . .
RUN npm run tsc && \
    npm run test

FROM node:10.17.0-alpine
ENV NODE_ENV production
WORKDIR /app
COPY --from=builder /app/package.json /app/package-lock.json ./
RUN npm ci
COPY --from=builder /app/build .
EXPOSE 4000
ENTRYPOINT [ "node", "./server.js" ]
