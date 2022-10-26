FROM node:16-alpine AS build

WORKDIR /app

RUN npm install -g @angular/cli

COPY . .

RUN npm install

CMD ng serve --host 0.0.0.0 --disable-host-check

EXPOSE 4200
