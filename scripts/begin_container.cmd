docker run --name oracle12c -d -p 1580:8080 -p 1521:1521 truevoly/oracle-12c

docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=SQLServer1433" -p 1433:1433 -d mcr.microsoft.com/mssql/server:2017-latest

docker run --detach --name some-mariadb -p 13306:3306 --env MARIADB_USER=example-user --env MARIADB_PASSWORD= --env MARIADB_ALLOW_EMPTY_ROOT_PASSWORD=true  mariadb:latest

docker run -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pgsql -it --rm postgres
