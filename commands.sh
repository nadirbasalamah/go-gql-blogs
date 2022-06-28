# pull the mongoDB image from the Docker Hub
docker pull mongo
# create a mongoDB container
docker run -itd --name my-mongo -v /your/volume/mongodb-volume:/data/db -p 27017:27017 mongo