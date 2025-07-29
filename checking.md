## Check2 Command
docker compose down -v
docker exec -it cdc-stock-consolidation-db-1 sh
psql -U admin -d stockdb


stockdb=# \dt
       List of relations
 Schema | Name  | Type  | Owner
--------+-------+-------+-------
 public | stock | table | admin
(1 row)

 INSERT INTO stock (product_id, branch_id, quantity, reserved) VALUES (1, 4, 100, 10);

 docker exec -it cdc-stock-consolidation-db-1 psql -U admin -d stockdb
LISTEN stock_changes;

NSERT 0 1
Asynchronous notification "stock_changes" with payload "{"id" : "1988f735-0bf7-4368-9dfc-13db193247a8", "product_id" : 1, "branch_id" : 2, "quantity" : 100, "reserved" : 10, "created_at" : "2025-07-29T05:17:55.443242", "updated_at" : "2025-07-29T05:17:55.443242"}" received from server process with PID 46.
stockdb=#  LISTEN stock_changes;

INSERT INTO stock (product_id, branch_id, quantity, reserved) VALUES (4, 2, 300, 30);

update stock
set quantity=2
where product_id = 3 and branch_id = 1;

LISTEN
Asynchronous notification "stock_changes" with payload "{"id" : "82a2279d-f97d-4c5f-988e-fb2c0a191555", "product_id" : 1, "branch_id" : 4, "quantity" : 100, "reserved" : 10, "created_at" : "2025-07-29T05:18:09.9969", "updated_at" : "2025-07-29T05:18:09.9969"}" received from server process with PID 56.
stockdb=#

docker-compose down; docker-compose up --build -d

docker exec -it   cdc-stock-consolidation-app-1 sh

docker exec cdc-stock-consolidation-app-1 env

docker logs -f cdc-stock-consolidation-app-1

docker exec -it cdc-stock-consolidation-app-1 /bin/sh


docker exec -i cdc-stock-consolidation-db-1 psql -U admin -d stockdb -c "INSERT INTO stock (product_id, branch_id, quantity, reserved) VALUES (4, 1, 300, 30);"

docker logs cdc-stock-consolidation-app-1

curl -v http://localhost:3000/health


cd D:/Training/AIEnhancementCourse/Assigment/FinalProject/cdc-stock-consolidation/; $env:DB_HOST="localhost"; $env:DB_PORT="5432"; $env:DB_USER="admin"; $env:DB_PASSWORD="admin"; $env:DB_NAME="stockdb"; $env:SERVICE_PORT="3000"; $env:HQ_END_POINT="http://localhost:8080"; $env:HQ_BASIC_AUTHORIZATION="Basic dXNlcjpwYXNz"; go test ./... -cover

cd D:/Training/AIEnhancementCourse/Assigment/FinalProject/cdc-stock-consolidation; $env:DB_HOST="localhost"; $env:DB_PORT="5432"; $env:DB_USER="admin"; $env:DB_PASSWORD="admin"; $env:DB_NAME="stockdb"; $env:SERVICE_PORT="3000"; $env:HQ_END_POINT="http://localhost:8080"; $env:HQ_BASIC_AUTHORIZATION="Basic dXNlcjpwYXNz"; go test ./... -cover

docker exec cdc-stock-consolidation-app-1 env | Select-String "DB_"


git init; git add .; git commit -m "Initial commit: CDC Stock Consolidation System"; git remote add origin https://github.com/talahoo/stock-consolidation.git; git push -u origin main

git push -u origin master

git add README.md; git commit -m "Add comprehensive README.md"; git push origin master

go test ./... -coverprofile=coverage.out; go tool cover -func=coverage.out


---- step by step checking

docker compose down -v

docker compose up --build -d

docker compose ps

docker logs cdc-stock-consolidation-app-1

docker exec -i cdc-stock-consolidation-db-1 psql -U admin -d stockdb -c "INSERT INTO stock (product_id, branch_id, quantity, reserved) VALUES (6, 3, 200, 25);"

docker logs cdc-stock-consolidation-app-1 --tail 10