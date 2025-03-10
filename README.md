## Overview

This repository is of orders service

## Endpoints

Method | Path                              | Description                                   |                                                                         
---    |-----------------------------------|------------------------------------------------
GET    | `/health`                         | Health page                                   |
GET    | `/metrics`                        | Страница с метриками                          |
GET    | `/v1/orders/list`                 | Получение заказов системы                     |
GET    | `/v1/orders/get/{id}`             | Получение заказа системы по id                |
GET    | `/v1/orders/get-by-user/{userId}` | Получение заказов пользователя                |
POST   | `/v1/orders/create`               | Создание нового заказа                        |