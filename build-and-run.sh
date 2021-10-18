docker build . --tag bradj.ca/car-audio-db &&
docker run -p 8080:8080 -v /home/bradsk88/GolandProjects/CarAudioDatabase/dbcreds.txt:/dbcreds.txt bradj.ca/car-audio-db;

