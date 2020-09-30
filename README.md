# DzProxy

#### Требования
- Mongodb :27017
- Свободные :3000, :8080, :5001 порты


### Запуск
- установить переменную в DzProxy/main.go folderForCerts = путь до DzProxy/certs
- go run main.go
- для запуска доп сервера для проверки атаки: go run ./testParam/main.go
- добавить root сертификат в браузер

### Работа с монго
- Посмотреть сохраненные запросы: http://localhost:3000/requests
- Сделать запрос и получить вывод: http://localhost:3000/request/[id] (id выбрать на странице с запросами)
- Атака: http://localhost:3000/attack/[id] 
- Ограничение по кол-ву запросов 500
