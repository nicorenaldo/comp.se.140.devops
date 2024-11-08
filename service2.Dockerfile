FROM python:3.12-slim

RUN apt-get update && apt-get install procps -y

WORKDIR /app

COPY service2-python/ ./

CMD ["python", "main.py"]