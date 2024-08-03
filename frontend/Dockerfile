FROM node:22-alpine3.19 as build

ENV NODE_ENV=production 

WORKDIR /app
COPY . .
RUN npm install --force && npm run build


FROM node:22-alpine3.19

WORKDIR /app
COPY --from=build /app .

ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=4173
ENV SERVER_ENABLE_CORS=false
EXPOSE 4173

CMD ["npm","run", "preview","--","--config","/app/vite.config.js"]