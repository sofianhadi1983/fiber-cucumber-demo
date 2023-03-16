# gofiber & Godog
## Ingredients

- [gofiber](https://gofiber.io/)
- [Gorm](https://gorm.io/)
- [Docker](https://www.docker.com/)
- [Testcontainer](https://golang.testcontainers.org/)
- [Godog](https://github.com/cucumber/godog)
- [Gherkin Reference](https://cucumber.io/docs/gherkin/reference/)

## How to run
In project root folder, please run following command:
```
go test -test.v -test.run ^TestFeatures$ ./tests
```