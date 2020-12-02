# docker build -t blog-docker .

docker run  --name blog1 --link mysql:mysql -p 8080:80 -d blog-docker
