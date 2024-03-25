FOR /f "tokens=*" %%i IN ('docker ps -a -q') DO docker stop %%i
FOR /f "tokens=*" %%i IN ('docker ps -a -q') DO docker rm %%i

docker image prune -f
FOR /f "tokens=*" %%i IN ('docker images -a -q') DO docker rmi -f %%i
