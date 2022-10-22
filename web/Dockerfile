FROM node:16-alpine AS build

WORKDIR /app

RUN npm install -g @angular/cli

COPY . .

RUN npm install

RUN ng build

FROM nginx:alpine

COPY --from=build /app/dist/web /usr/share/nginx/html
