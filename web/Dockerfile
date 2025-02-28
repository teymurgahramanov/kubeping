FROM python:3.9-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
ENV GUNICORN_CMD_ARGS="--workers 2 --bind 0.0.0.0:8000 --log-level info"
ENTRYPOINT ["gunicorn"] 
CMD ["app:app"]