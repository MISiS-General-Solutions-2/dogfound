# Запуск сервиса

### Требования:
git
docker версии 18.09 или выше
docker-compose

### Сборка и запуск контейнера

Скачать содержимое репозитория
```git clone https://github.com/MISiS-General-Solutions-2/dogfound```
Для использования BuildKit задать переменные окружения DOCKER_BUILDKIT=1 и COMPOSE_DOCKER_CLI_BUILD=1. Например, на системе linux, это можно сделать коммандами
```export DOCKER_BUILDKIT=1```
```export COMPOSE_DOCKER_CLI_BUILD=1```
**Изображения из папки data/new_images будут классифицированы и добавлены в базу данных, если поместить их туда перед сборкой контейнера.**

Перейти в директорию с проекта
```cd dogfound```
Установить адрес сервера, на котором будет запущен сервер в front/Dockerfile
```REACT_APP_API_URL="ddd.ddd.ddd.ddd:pppp"```
Собрать и запустить контейнер
```docker-compose up```
Чтобы пересобрать контейнер с учетом изменений в исходном коде
```docker-compose build```

Веб-страница сервиса доступна на порту 1022.

# API
### Публичное API
**API доступно на порту 1022**
### POST /api/image/by-classes
```
{
    ["color"]: <цвет собаки>
    [ "tail"]: <длина хвоста собаки>
    [ "cam_id"]: <айди камеры>
    ["t0"]: <время, начиная с которого искать собаку>
    ["t1"]: <время, по которое искать собаку>
}
```

- color:	int	- Метка класса цвета животного
- tail:	int - Метка класса хвоста животного
- cam_id: string - Айди камеры, на которой осуществлять поиск
- t0: int - Unix время, начиная с которого искать собаку
- t1: int - Unix время, заканчивая которым искать собаку.
 
Ответ
200 OK
```
{
    "filename": <имя файла>,
    "address": <адрес>,
    "cam_id": <айди камеры>,
    "timestamp": <временная метка снимка>,
    "lonlat": []: <широта и долгота>,
    "breed": <порода собаки>
    "additional": {
      "crop": [], <координаты кропа собаки>,
    }
```
- filename:	string	- Имя файла - снимка с камеры
- address:	string - Адрес камеры. Может быть пустым, если адрес камеры не был найден.
- cam_id: string - Айди камеры. Может быть пустым, если id камеры не был распознан.
- timestamp: int - Unix время снимка с точностью до дня. Может быть 0, если время не было распознано.
- lonlat: [2]float - LonLat координаты камеры. Могут быть 0, если адрес камеры не был найден.
- breed: string - Порода собаки. Может быть пустым, если порода не определена.
- crop: [4]int - Координаты прямоугольника с собакой на изображении. Координаты левого верхнего и правого нижнего угла изображения в системе координат OpenCV.
- 
### GET /api/image/:name
name - имя файла, полученное из запроса /api/image/by-classes.

Поля запроса:
- omit_crop - Не выделять собаку на изображении. Не выделяет собаку, если значение равно 1.

Ответ
200 OK
В ответе возвращается запрошенное изображение.

### Приватное API
**Приватное API доступно на порту 6000**
### PUT /image/upload
Запрос в формате multipart/form-data. Добавляет изображения из поля file в базу данных, они в дальнейшем размечаются и используются при выдаче результата.

### Добавление изображений в базу
- Контейнер использует volume data, в котором папка /opt/dogfound/data/new_images проверяется на наличие новых изображений по умолчанию каждые 5 секунд. Можно добавить изображения в эту папку, например, командой
```docker cp ./core/my_data/new_images dogfound-core-1:/opt/dogfound/data/```
- Можно воспользоваться приватным API /image/upload

# Ноутбук с построением csv файла на датасете и AI
Находится в neural_network/notebooks/build_test_csv.ipynb
"Боевая" логика для сервиса - neural_network/detect.py
Обученные модели - YOLOv5l и resnet38 в neural_network/models
