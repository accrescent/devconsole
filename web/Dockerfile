FROM node:19-alpine AS build

WORKDIR /app

RUN npm install -g @angular/cli

COPY package.json .

RUN npm install

COPY . .

RUN ng build --configuration development

FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf

COPY --from=build /app/dist/web /usr/share/nginx/html
