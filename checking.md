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


docker exec -i cdc-stock-consolidation-db-1 psql -U admin -d stockdb -c "INSERT INTO stock (product_id, branch_id, quantity, reserved) VALUES (4, 1, 300, 30);"



docker logs cdc-stock-consolidation-app-1 --tail 10

git status

git add .; git commit -m "Update checking.md with testing steps and fix database connection"

git add .; git commit -m "Refactoring for linter : change go.mod"

git add .; git commit -m "Revamp for workflow Test for Set up job"

git add .; git commit -m "change internal adapter db for lister test"

git add .; git commit -m "change pkg,test,adapter"



git push origin master

--
go test ./... -cover

https://github.com/BINAR-Learning/demo-repository/tree/Module-8

go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; golangci-lint run


---CI/CD

# Pastikan semua perubahan sudah di-commit
git add .github/workflows/ci.yml
git add .github/workflows/go.yml

git commit -m "feat: revamp GitHub Actions workflow for lint and test"

# Push ke GitHub
git push origin master

Di Browser GitHub:

Di Browser GitHub:

Buka repository Anda di https://github.com/talahoo/stock-consolidation
Klik tab "Actions" di bagian atas repository
Anda akan melihat workflow "Go CI/CD" yang baru saja di-push
Workflow akan otomatis berjalan untuk setiap push ke master
Memonitor Workflow:

Di tab Actions, klik pada workflow run yang sedang berjalan
Anda bisa melihat 2 jobs: "Lint" dan "Test"
Klik pada masing-masing job untuk melihat detail prosesnya
Jika ada error, log lengkap akan tersedia di sini
Mengecek Coverage Report:

Setelah workflow selesai, klik pada workflow run yang sukses
Di bagian "Artifacts", Anda akan melihat "coverage-report"
Download dan buka file coverage.out untuk melihat detail coverage
Troubleshooting (jika diperlukan):

Jika lint fail:

Lihat error message di log
Perbaiki masalah formatting atau linting
Commit dan push lagi
Jika test fail:

Cek log untuk error detail
Pastikan environment variables sudah benar
Cek koneksi ke PostgreSQL container
Best Practices:

Selalu test perubahan locally sebelum push
Gunakan golangci-lint run di local untuk cek masalah
Jalankan test dengan go test ./... -cover di local
Mengecek Badge Status (opsional):

Di GitHub, pergi ke tab Actions
Klik "..." di samping workflow
Pilih "Create status badge"
Copy markdown dan tambahkan ke README.md

golangci-lint run

go test ./... -v -short

go test ./pkg/logger -v

go test ./pkg/logger -v -count=1

go test ./... -coverprofile=coverage.out

go test ./... -v -short

go test ./... -cover -short