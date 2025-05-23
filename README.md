# Тестовое задание на стажировку avito-tech 2023
Cервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент). Имеется возможность создания csv отчетов, добавления сегментов для пользователя на ограниченное время и автоматическое добавление сегментов проценту пользователей.

Выполнено:
- [x] Основное задание
- [x] Создание csv отчетов
- [x] TTL сегментов
- [x] Автоматическое добавление пользователей в сегмент
- [x] Swagger документация
- [x] Валидация запросов 
- [x] Покрытие кода тестами (+/-)
- [x] Мониторинг и визуализация метрик с помощью Prometheus и Grafana

---
Библиотеки и технологии:
- PostgreSQL (хранилище)
- pgx (драйвер базы данных)
- golang-migrate/migrate (миграции базы данных)
- viper (конфигурация)
- chi-router (роутер)
- swaggo (создание документации)
- go-playground/validator (валидация запросов)
- zap (логирование)
- docker (деплой)
- asynq (очередь отложенных задач)
- test-containers (тестирование БД)


# Локальный запуск

1. Клонирование проекта
```
git clone https://github.com/dezzerlol/avitotech-test-2023.git
```
2. Запуск сервиса, prometheus, grafana, redis, postgres, и автоматическая миграция. (сервис запускается с задержкой 5 сек. после запуска БД, также должен быть запущен docker)
```
make compose
```

*.env файл с необходимыми переменными оставлен в корне проекта*

Swagger документация доступна по ссылке `http://localhost:8080/swagger/index.html#/`

Для запуска тестов используется команда `make test` (должен быть запущен docker).

# Примеры запросов
[Создание пользователя](#1-создание-пользователя)  
[Создание сегмента](#2-создание-сегмента)  
[Удаление сегмента](#3-удаление-сегмента)  
[Добавление/удаление сегментов пользователя](#4-добавлениеудаление-сегментов-пользователя)  
[Получение всех сегментов пользователя](#5-получение-всех-сегментов-пользователя)  
[Создание отчета добавления/удаления сегментов пользователя](#6-создание-отчета-добавленияудаления-сегментов-пользователя)  
[Скачивание отчета по сегментам](#7-скачивание-отчета-по-сегментам)


### 1. **Создание пользователя**

Используется в случае необходимости вручную добавить пользователя, так как при добавлении сегмента пользователю, id этого пользователя сохраняется автоматически.

Запрос:
```
curl --request POST 'http://localhost:8080/user'
```

Ответ:
```
{"user_id":1}
```

### 2. **Создание сегмента**
Принимает `slug` - название сегмента.

Если указан `user_percent`, то сегмент будет добавлен случайным пользователям в заданном проценте от общего числа (прим. Задали 50%, в таблице сохранено 200 пользователей, таким образом 100 случайных пользователей получат сегмент).

Запрос:
```
curl --request POST -d '{"slug": "AVITO_DISCOUNT_30"}' 'http://localhost:8080/segment'
```

Ответ:
```
{"created_at":"2023-08-28T08:16:24.21653Z"}
```

### 3. **Удаление сегмента**
Принимает `slug` - название сегмента. В случае если удален сегмент, который активен у пользователей, то он будет удален у всех пользователей.

Запрос:
```
curl --request DELETE -d '{"slug": "AVITO_DISCOUNT_30"}' 'http://localhost:8080/segment'
```

Ответ:
```
{"message":"ok"}
```

### 4. **Добавление/удаление сегментов пользователя**
Принимает `user_id` - id пользователя, `add_segments` - список сегментов которые нужно добавить пользователю и `delete_segments` - список сегментов которые нужно удалить.
В случае если указан `ttl` в секундах, добавялет сегменты пользователю на определенный промежуток времени.


Пример запроса, добавляющего 2 сегмента на 86400 секунды (1 день):
```
curl --request POST \
-d '{"user_id": 1, "add_segments": ["AVITO_DISCOUNT_50", "AVITO_DISCOUNT_30"], "ttl":86400, "delete_segments": []}' \
'http://localhost:8080/segment/user'
```

Ответ:
```
{"segments_added":2,"segments_deleted":0}
```

### 5. **Получение всех сегментов пользователя**
Принимает `id пользователя` в качестве url param.

Запрос:
```
curl --request GET 'http://localhost:8080/segment/user/1'
```

Ответ:
```
{"segments":[{"slug":"AVITO_DISCOUNT_30"},{"slug":"AVITO_DISCOUNT_50"}]}
```

### 6. **Создание отчета добавления/удаления сегментов пользователя**
Принимает `id пользователя` в качестве url param, `year` и `month` в виде query param. Возвращает ссылку на скачивание csv отчета.

Запрос:
```
curl --request GET 'http://localhost:8080/segment/history/1?year=2023&month=9'
```

Ответ:
```
{"report_link":"localhost:8080/segment/reports/1-1693218476.csv"}
```

### 7. **Скачивание отчета по сегментам**
Принимает `название файла` (полученное при создании отчета) в качестве url param. Возвращает csv файл с отчетом в формате: id пользователя, slug сегмента, операция (I = создание, D = удаление), дата и время.

Запрос:
```
curl --request GET 'http://localhost:8080/segment/reports/1-1693218476.csv'
```

Ответ (скачивание файла):
```
user_id,segment_slug,operation,executed_at
1,AVITO_DISCOUNT_50,I,2023-08-28 10:25:25
1,AVITO_DISCOUNT_30,I,2023-08-28 10:25:25
1,AVITO_DISCOUNT_50,D,2023-08-28 10:25:55
1,AVITO_DISCOUNT_30,D,2023-08-28 10:25:55
1,AVITO_DISCOUNT_50,I,2023-08-28 10:27:46
1,AVITO_DISCOUNT_30,I,2023-08-28 10:27:49
```
Пример отчета: [файл](/reports/1-1693224806.csv)


# FAQ
1. Почему для истории сегментов используется отдельная таблица, а не поле в таблице сегментов?
    > Для того, чтобы не хранить в таблице сегментов лишние данные, т.к. таблица сегментов будет использоваться чаще чем таблица истории сегментов.

2. Как реализована таблица истории сегментов?
    > С помощью PostgreSQL триггера, который при добавлени/удалении сегмента, добавляет запись в таблицу истории сегментов.

3. Как реализовано автоматическое удаление пользователя из сегмента?
    > Пользователь указывает ttl в секундах, через которое нужно удалить сегмент, добавляется отложенная задача, которая будет выполняться через указанное время и удалять сегмент у пользователя. Для этого используется asynq, который позволяет добавлять отложенные задачи в очередь. В качестве брокера сообщений используется Redis.

4. Реализация отчетов.
    > При каждом добавлении/удалении сегментов у пользователя, срабатывает триггер PostgreSQL, который сохраняет запись в таблице истории. При запросе отчета от пользователя генерируется файл и ссылка на скачивание этого файла. Пользователь переходит по ссылке и скачивает отчет. (файл сохраняется внутри проекта, для production лучше переписать код и использовать облачное хранилище).