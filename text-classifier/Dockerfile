# ============================
# STAGE 1 — BUILDER
# ============================
FROM python:3.11 AS builder

# Versión de Poetry
ENV POETRY_VERSION=2.2.1

# Instalación recomendada por Poetry (dentro de ~/.local)
ENV POETRY_HOME="/root/.local"
ENV PATH="/root/.local/bin:$PATH"

# Poetry crea .venv dentro del proyecto
ENV POETRY_VIRTUALENVS_IN_PROJECT=true
ENV POETRY_NO_INTERACTION=1

# Dependencias necesarias
RUN apt-get update && apt-get install -y --no-install-recommends \
        curl \
        python3-venv \
        build-essential \
    && rm -rf /var/lib/apt/lists/*

# Instalar Poetry
RUN curl -sSL https://install.python-poetry.org | python3 -

WORKDIR /app

# Copiar archivos de Poetry
COPY pyproject.toml poetry.lock ./

# Instalar dependencias (sin dev)
RUN poetry install --no-root

# Copiar el resto de la app
COPY . .

# ============================
# STAGE 2 — FINAL
# ============================
FROM python:3.11-slim

WORKDIR /app

# Copiar app completa + .venv
COPY --from=builder /app /app

# Activar virtualenv
ENV PATH="/app/.venv/bin:$PATH"

RUN apt-get update && apt-get install -y --no-install-recommends \
        curl \
    && rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD ["curl", "-f", "http://localhost:8000/docs"] || exit 1

CMD ["python", "api.py"]
