FROM python:3.12-slim

WORKDIR /app
COPY proxy.py requirements.txt ./
RUN pip install -r requirements.txt

EXPOSE 8000
CMD ["python", "proxy.py"]
