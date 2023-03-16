package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"reflect"
	"time"

	"github.com/cucumber/godog"
	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tiketdatasofianhadiant0/bdd-demo/database"
	"github.com/tiketdatasofianhadiant0/bdd-demo/models"
	"github.com/tiketdatasofianhadiant0/bdd-demo/routes"
)

type godogsResponseCtxKey struct{}

func init() {
	ctx := context.Background()

	const dbname = "test-db"
	const user = "postgres"
	const password = "password"

	port, _ := nat.NewPort("tcp", "5432")

	container, _ := startContainer(ctx,
		WithPort(port.Port()),
		WithInitialDatabase(user, password, dbname),
		WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5*time.Second)),
	)
	containerPort, _ := container.MappedPort(ctx, port)
	host, _ := container.Host(ctx)

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", containerPort.Port())
	os.Setenv("DB_USER", user)
	os.Setenv("DB_PASS", password)
	os.Setenv("DB_NAME", dbname)
}

type apiFeature struct {
	app *fiber.App
}

type response struct {
	status int
	body   any
}

func (a *apiFeature) resetResponse(*godog.Scenario) {
	a.app = fiber.New()
	routes.SetupRoutes(a.app)
	database.ConnectDb()
}

func (a *apiFeature) iSendRequestToWithPayload(ctx context.Context, method, route string, payloadDoc *godog.DocString) (context.Context, error) {
	var reqBody []byte

	if payloadDoc != nil {
		payloadMap := models.Book{}
		err := json.Unmarshal([]byte(payloadDoc.Content), &payloadMap)
		if err != nil {
			panic(err)
		}

		reqBody, _ = json.Marshal(payloadMap)
	}

	req := httptest.NewRequest(method, route, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := a.app.Test(req)
	var createdBooks []models.Book
	json.NewDecoder(resp.Body).Decode(&createdBooks)

	actual := response{
		status: resp.StatusCode,
		body:   createdBooks,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (a *apiFeature) theResponseCodeShouldBe(ctx context.Context, expectedStatus int) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if expectedStatus != resp.status {
		if resp.status >= 400 {
			return fmt.Errorf("expected response code to be: %d, but actual is: %d, response message: %s", expectedStatus, resp.status, resp.body)
		}
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", expectedStatus, resp.status)
	}

	return nil
}

func (a *apiFeature) theResponsePayloadShouldMatchJson(ctx context.Context, expectedBody *godog.DocString) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	books := make([]models.Book, 0)

	err := json.Unmarshal([]byte(expectedBody.Content), &books)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(actualResp.body, books) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expectedBody, actualResp.body)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		api.resetResponse(sc)
		return ctx, nil
	})

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" with payload:$`, api.iSendRequestToWithPayload)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response payload should match json:$`, api.theResponsePayloadShouldMatchJson)
}
