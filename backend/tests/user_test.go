//nolint:noctx,funlen
package tests_test

import (
	"context"
	"dish_as_a_service/assembly"
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"

	"dish_as_a_service/repository"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Falokut/go-kit/test/fake"
	"github.com/Falokut/go-kit/test/telegramt"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/http/client"
	"github.com/Falokut/go-kit/test"
	"github.com/Falokut/go-kit/test/dbt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/bgjob"
)

type UserSuite struct {
	suite.Suite
	test *test.Test

	db       *dbt.TestDb
	userRepo repository.User
	cli      *client.Client
}

func TestUser(t *testing.T) {
	t.Parallel()
	suite.Run(t, &UserSuite{})
}

// nolint:dupl
func (t *UserSuite) SetupTest() {
	test, _ := test.New(t.T())
	t.test = test
	t.db = dbt.New(test, db.WithMigrationRunner("../migrations", test.Logger()))
	t.userRepo = repository.NewUser(t.db.Client)

	bgjobDb := bgjob.NewPgStore(t.db.Client.DB.DB)
	bgjobCli := bgjob.NewClient(bgjobDb)
	tgBot, _ := telegramt.TestBot(test)
	locatorCfg, err := assembly.Locator(
		context.Background(),
		test.Logger(),
		t.db.Client,
		nil,
		tgBot,
		bgjobCli,
		getConfig(),
	)
	t.Require().NoError(err)
	server := httptest.NewServer(locatorCfg.HttpRouter)
	t.cli = client.NewWithClient(server.Client())
	t.cli.GlobalRequestConfig().BaseUrl = fmt.Sprintf("http://%s", server.Listener.Addr().String())
	t.T().Cleanup(func() {
		server.Close()
	})
}

func (t *UserSuite) Test_IsAdmin_HappyPath() {
	var userId string
	err := t.db.Get(&userId,
		`INSERT INTO users(username,name,admin)
		VALUES($1,$2,$3)
		RETURNING id;`,
		"@test_admin",
		"test",
		true,
	)
	t.Require().NoError(err)

	var isAdmin domain.IsUserAdminResponse
	_, err = t.cli.Get("/users/is_admin").
		QueryParams(map[string]any{
			"userId": userId,
		}).
		StatusCodeToError().
		JsonResponseBody(&isAdmin).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().True(isAdmin.IsAdmin)

	err = t.db.Get(&userId,
		`INSERT INTO users(username,name,admin)
		VALUES($1,$2,$3)
		RETURNING id;`,
		"@test_no_admin",
		"test",
		false,
	)
	t.Require().NoError(err)
	_, err = t.cli.Get("/users/is_admin").
		QueryParams(map[string]any{
			"userId": userId,
		}).
		StatusCodeToError().
		JsonResponseBody(&isAdmin).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().False(isAdmin.IsAdmin)
}

func (t *UserSuite) Test_IsAdmin_NotFound() {
	userId := uuid.NewString()
	resp, err := t.cli.Get("/users/is_admin").
		QueryParams(map[string]any{
			"userId": userId,
		}).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().Equal(http.StatusNotFound, resp.StatusCode())
}

func (t *UserSuite) Test_GetUserIdByTelegramId_HappyPath() {
	var telegramId = fake.It[int64]()
	err := t.userRepo.Register(context.Background(), entity.RegisterUser{
		Id:       uuid.NewString(),
		Username: fake.It[string](),
		Name:     fake.It[string](),
		Telegram: &entity.Telegram{
			ChatId: fake.It[int64](),
			UserId: telegramId,
		},
	})
	t.Require().NoError(err)

	var userId string
	err = t.db.Get(&userId, "SELECT id FROM users_telegrams WHERE telegram_id=$1", telegramId)
	t.Require().NoError(err)

	var userIdRep domain.GetUserIdByTelegramIdResponse
	_, err = t.cli.Get(fmt.Sprintf("/users/get_by_telegram_id/%d", telegramId)).
		JsonResponseBody(&userIdRep).
		StatusCodeToError().
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().Equal(userId, userIdRep.UserId)
}

func (t *UserSuite) Test_GetUserIdByTelegramId_NotFound() {
	var telegramId = fake.It[int64]()
	resp, err := t.cli.Get(
		fmt.Sprintf("/users/get_by_telegram_id/%d", telegramId),
	).Do(context.Background())
	t.Require().NoError(err)
	t.Require().Equal(http.StatusNotFound, resp.StatusCode())
}
