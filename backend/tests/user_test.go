//nolint:noctx,gosec
package tests_test

import (
	"context"
	"dish_as_a_service/assembly"
	"dish_as_a_service/domain"
	"dish_as_a_service/repository"
	"dish_as_a_service/service"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/json"
	"github.com/Falokut/go-kit/test"
	"github.com/Falokut/go-kit/test/dbt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/bgjob"
)

type UserSuite struct {
	suite.Suite
	test *test.Test

	db          *dbt.TestDb
	userRepo    repository.User
	cli         *http.Client
	serverAddr  string
	userService service.User
}

func TestUser(t *testing.T) {
	t.Parallel()
	suite.Run(t, &UserSuite{})
}

func (t *UserSuite) SetupTest() {
	test, _ := test.New(t.T())
	t.test = test
	t.db = dbt.New(test, db.WithMigrationRunner("../migrations", test.Logger()))
	t.userRepo = repository.NewUser(t.db.Client)
	t.userService = service.NewUser(t.userRepo, nil)

	bgjobDb := bgjob.NewPgStore(t.db.Client.DB.DB)
	bgjobCli := bgjob.NewClient(bgjobDb)

	locatorCfg, err := assembly.Locator(context.Background(), test.Logger(), t.db.Client, nil, bgjobCli, getConfig())
	t.Require().NoError(err)
	server := httptest.NewServer(locatorCfg.HttpRouter)
	t.serverAddr = server.Listener.Addr().String()
	t.cli = server.Client()
	t.T().Cleanup(func() {
		server.Close()
	})
}

func (t *UserSuite) getServerUrl(endpoint string) string {
	return fmt.Sprintf("http://%s/%s", t.serverAddr, endpoint)
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

	resp, err := t.cli.Get(t.getServerUrl("users/is_admin?userId=" + userId))
	t.Require().NoError(err)
	defer resp.Body.Close()

	var isAdmin domain.IsUserAdminResponse
	err = json.NewDecoder(resp.Body).Decode(&isAdmin)
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

	resp, err = t.cli.Get(t.getServerUrl("users/is_admin?userId=" + userId))
	t.Require().NoError(err)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&isAdmin)
	t.Require().NoError(err)
	t.Require().False(isAdmin.IsAdmin)
}

func (t *UserSuite) Test_IsAdmin_NotFound() {
	userId := uuid.NewString()
	resp, err := t.cli.Get(t.getServerUrl("users/is_admin?userId=" + userId))
	t.Require().NoError(err)
	defer resp.Body.Close()
	t.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (t *UserSuite) Test_GetUserIdByTelegramId_HappyPath() {
	var telegramId = rand.Int64()
	err := t.userService.Register(context.Background(), domain.RegisterUser{
		Username: "user_with_tg",
		Name:     "some_name",
		Telegram: &domain.Telegram{
			ChatId: rand.Int64(),
			UserId: telegramId,
		},
	})
	t.Require().NoError(err)

	var userId string
	err = t.db.Get(&userId, "SELECT id FROM users_telegrams WHERE telegram_id=$1", telegramId)
	t.Require().NoError(err)

	resp, err := t.cli.Get(t.getServerUrl(
		fmt.Sprintf("users/get_by_telegram_id/%d", telegramId)),
	)
	t.Require().NoError(err)
	defer resp.Body.Close()

	var userIdRep domain.GetUserIdByTelegramIdResponse
	err = json.NewDecoder(resp.Body).Decode(&userIdRep)
	t.Require().NoError(err)
	t.Require().Equal(userId, userIdRep.UserId)
}

func (t *UserSuite) Test_GetUserIdByTelegramId_NotFound() {
	var telegramId = rand.Int64()
	resp, err := t.cli.Get(t.getServerUrl(
		fmt.Sprintf("users/get_by_telegram_id/%d", telegramId)),
	)
	t.Require().NoError(err)
	defer resp.Body.Close()
	t.Require().Equal(http.StatusNotFound, resp.StatusCode)
}
