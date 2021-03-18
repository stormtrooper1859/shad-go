## dockertest

Это не настоящая задача, а заготовка на будущее.

### Что нужно сделать?

Установить docker и добиться успешного **локального** запуска тестов
```
go test -v ./dockertest/... -count=1
```

Только **после того, как тесты пройдут локально** можете запушить решение в систему.

### С чего начать?

#### Установить docker

https://docs.docker.com/engine/install/

После стандартной процедуры установки на Linux будет создана группа `docker`.
Чтобы использовать docker cli без sudo, нужно добавить себя в эту группу:
```
sudo groupadd docker
```
Для проверки можно запустить
```
docker run hello-world
```

#### Установить docker-compose

https://docs.docker.com/compose/install/

#### Запустить контейнеры не через тесты

В директории `dockertest` выполнить
```
docker-compose up
```

### Что делать, если сразу не заработало?

Поискать решение проблемы в интернете.

Если решение найдено, и проблема выглядит общей, сделать merge request с улучшением README.

Если интернет не помог, спросить в чате.

### docker-compose cheat sheet

Запустить все контейнеры в daemon режиме пересобрав образы:
```
docker-compose up -d --build
```

Остановить все контейнеры:
```
docker-compose down
```

### Docker cheat sheet

Получить список образов
```
docker images
```

Список всех контейнеров:
```
docker ps -a
```

Остановить контейнер:
```
docker stop <NAME>
```

Удалить контейнер:
```
docker rm <NAME>
```

Удалить образ:
```
docker rmi <NAME>
```