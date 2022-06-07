FROM maven:3-jdk-11-slim AS MAVEN_TOOL_CHAIN
COPY pom.xml /tmp/
COPY src /tmp/src/
WORKDIR /tmp/
RUN mvn package

FROM gcr.io/distroless/java17
COPY --from=MAVEN_TOOL_CHAIN /tmp/target/app.jar /app.jar
WORKDIR /

ENV MQTT_HOST=mosquitto
ENV BACKUP_PATH=/domotik/backup.json

ENTRYPOINT ["java", "-jar", "app.jar"]