FROM python:3.13-alpine
WORKDIR /app
RUN chmod 777 /app

COPY requirements.txt ./
RUN pip install -r requirements.txt

USER 2000:2000
RUN id | tee 2000.log

USER 1000:1000
RUN id | tee 1000.log

COPY main.py ./
CMD [ "python", "./main.py"]
