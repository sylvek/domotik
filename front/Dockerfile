FROM ghcr.io/sylvek/domotik-wallpapers:latest AS WALLPAPER

FROM library/node:14-alpine
COPY domotik /domotik
COPY --from=WALLPAPER /img /domotik/public/img
COPY --from=WALLPAPER /data.js /domotik/public/javascripts/data.js

WORKDIR /domotik
RUN apk --no-cache add git
RUN npm install
RUN node_modules/bower/bin/bower install --allow-root

ENV DB_PATH=/database
EXPOSE 3000
ENTRYPOINT ["npm", "start"]
