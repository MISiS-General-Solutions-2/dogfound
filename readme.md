# Запуск сервиса

### Требования:
git
docker версии 18.09 или выше
docker-compose

### Сборка и запуск контейнера

Скачать содержимое репозитория
```git clone https://github.com/MISiS-General-Solutions-2/dogfound```
Задать переменную окружения DOCKER_BUILDKIT=1. Например, это можно сделать коммандой
```export DOCKER_BUILDKIT=1```
**Изображения из папки data/new_images будут классифицированы и добавлены в базу данных, если поместить их туда перед сборкой контейнера.**

Перейти в директорию с проекта
```cd dogfound```
Установить адрес сервера, на котором будет запущен сервер в front/Dockerfile
```REACT_APP_API_URL="ddd.ddd.ddd.ddd:pppp"```
Собрать и запустить контейнер
```docker-compose up```
Чтобы пересобрать контейнер с учетом изменений в исходном коде
```docker-compose build```

Для добавления изображений в базу данных добавить их в /opt/dogfound/data/new_images контейнера, например командой
```docker cp ./core/my_data/new_images dogfound-core-1:/opt/dogfound/data/```

Веб-страница сервиса доступна на порту 1022.

# Ноутбук с построением csv файла на датасете и AI
Находится в neural_network/notebooks/build_test_csv.ipynb
"Боевая" логика для сервиса - neural_network/detect.py
Обученные модели - YOLOv5l и resnet38 в neural_network/models
