FROM python:3.12-slim

RUN apt-get update && apt-get install procps -y

WORKDIR /app

COPY service1-python/ ./
RUN pip install -r requirements.txt

CMD ["python", "main.py"]