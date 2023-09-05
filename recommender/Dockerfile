FROM python:3.9-slim-bookworm as builder

WORKDIR /app

COPY ./scripts /app
COPY ./requirements.txt /app

# install dependencies
RUN pip install -r /app/requirements.txt
