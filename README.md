## Overview

This repository is of authenticator service

## Endpoints

Method | Path        | Description                                   |                                                                         
---    |-------------|------------------------------------------------
POST   | `/register` | Регистрация нового пользователя               |
POST   | `/login`    | Логин пользователя                            |
POST   | `/auth`     | Авторизация пользователя                      |
GET    | `/signin`   | Информация о необходимости авторизации        |
POST   | `/logout`   | Разлогин пользователя                         |
GET    | `/sessions` | Получение всех сессий                         |
