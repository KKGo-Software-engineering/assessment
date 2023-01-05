unit test
=========

go test -v --tags=unit github.com/sutthiphong2005/assessment/rest/handler


integration test
==========
1) run
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests
2) tear down
docker-compose -f docker-compose.test.yml down


server.go
==========
export PORT=:2565
export DATABASE_URL=postgres://pupffhjj:cTLk0BZ4OkVPGze0vhiED7wOZjO5ZMyN@tiny.db.elephantsql.com/pupffhjj

go run server.go
