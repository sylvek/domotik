FROM python:2.7-alpine
COPY domotik /domotik

WORKDIR /domotik

RUN ["pip", "install", "pymongo"]

ENTRYPOINT ["python", "/domotik/weekly.py"]
