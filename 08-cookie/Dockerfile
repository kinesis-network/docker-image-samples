FROM python
WORKDIR /app
COPY main.py requirements.txt .
RUN pip install --no-cache-dir --upgrade -r requirements.txt
CMD ["python", "main.py"]
