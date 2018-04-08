FROM node:9
COPY app /app

EXPOSE 3000
WORKDIR /app
RUN ["npm", "install"]
ENTRYPOINT ["node", "/app/index.js"]
