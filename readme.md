# Запуск сервиса

git clone https://github.com/MISiS-General-Solutions-2/dogfound
docker-compose up

Для добавления изображений в базу данных добавить их в /opt/dogfound/data/img контейнера, например командой
docker cp ./core/my_data/img dogfound-core-1:/opt/dogfound/data/

Или же можно добавить изображения перед сборкий контейнера в папку data/img.

# Запуск веб-страницы

Необходим установленный Node.js

Из папки front
npm install yarn -
yarn start

Веб-страница сервиса доступна на порте 1022. На локальной машине доступ можно получить по адресу http://localhost:1022/

# Ноутбук с построением csv файла на датасете и AI
Находится в neural_network/notebooks/build_test_csv.ipynb
"Боевая" логика для сервиса - neural_network/detect.py
Обученные модели - YOLOv5l и resnet38 в neural_network/models
