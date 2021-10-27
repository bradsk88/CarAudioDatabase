docker build . --tag bradj.ca/car-audio-db &&
docker run -p 8080:8080 \
-v /home/bradsk88/GolandProjects/CarAudioDatabase/dbcreds.txt:/dbcreds.txt \
-v /home/bradsk88/client_secret_234930237005-t3q14pe976i40c5v6ghvqmi3j6lhd2ut.apps.googleusercontent.com.json:/credentials.json \
-v /home/bradsk88/GolandProjects/CarAudioDatabase/session.key:/session.key \
bradj.ca/car-audio-db;

