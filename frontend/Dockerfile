FROM node:23-alpine3.20 AS build

ENV NODE_ENV=production 

WORKDIR /app
COPY . .
RUN npm install && npm run build


FROM node:23-alpine3.20

WORKDIR /app
COPY --from=build /app .

ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV SERVER_ENABLE_CORS=false
EXPOSE 4173

CMD ["npm","run", "preview","--","--config","/app/vite.config.js"]