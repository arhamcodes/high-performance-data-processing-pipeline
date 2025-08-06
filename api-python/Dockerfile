FROM python:3.12-slim
WORKDIR /high-performance-data-processing-pipeline/api-python
RUN python -m pip install --upgrade pip && \
    python -m pip install --no-cache-dir uv build
COPY pyproject.toml uv.lock ./
COPY . .
EXPOSE 8000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]