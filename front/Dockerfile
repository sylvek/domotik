FROM library/node:4
COPY domotik /domotik

EXPOSE 3000
WORKDIR /domotik

RUN ["npm", "install"]
RUN ["node_modules/bower/bin/bower", "install", "--allow-root"]

ENTRYPOINT ["npm", "start"]
