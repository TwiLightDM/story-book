# Story Book

Проект **Story Book** — это backend-приложение на Go с использованием PostgreSQL, Redis и MongoDB.

## Технологии

* **Go:** 1.25.0
* **PostgreSQL:** 17
* **Redis:** 7
* **Docker & Docker Compose** для локального развертывания сервисов

## Установка и запуск

### 1. Установите язык Go версии 1.25+
https://go.dev/doc/install

### 2. Клонируйте репозиторий

```bash
git clone https://github.com/TwiLightDM/story-book.git
cd story-book
```
### 3. Запуск сервисов через Docker Compose

```bash
docker-compose up
```
### 4. Запуск Go-приложения

```bash
go mod tidy
go run main.go
```
