# Запуск сервиса

git clone https://github.com/MISiS-General-Solutions-2/dogfound
docker-compose up

Для добавления изображений в базу данных добавить их в /opt/dogfound/data/img контейнера, например командой
docker cp ./core/my_data/img dogfound-core-1:/opt/dogfound/data/

Или же можно добавить изображения перед сборкий контейнера в папку data/img.

Веб-страница сервиса доступна на порте 3001. На локальной машине доступ можно получить по адресу http://localhost:3001/