//nolint:noctx,funlen
package tests_test

import (
	"context"
	"dish_as_a_service/assembly"
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"
	"dish_as_a_service/jwt"

	"dish_as_a_service/repository"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Falokut/go-kit/test/telegramt"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/http/client"
	"github.com/Falokut/go-kit/test"
	"github.com/Falokut/go-kit/test/dbt"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/bgjob"
)

type AuthSuite struct {
	suite.Suite
	test             *test.Test
	adminAccessToken string
	userAccessToken  string

	db       *dbt.TestDb
	userRepo repository.User
	cli      *client.Client
}

func TestUser(t *testing.T) {
	t.Parallel()
	suite.Run(t, &AuthSuite{})
}

// nolint:dupl
func (t *AuthSuite) SetupTest() {
	test, _ := test.New(t.T())
	t.test = test
	t.db = dbt.New(test, db.WithMigrationRunner("../migrations", test.Logger()))
	t.userRepo = repository.NewUser(t.db.Client)

	bgjobDb := bgjob.NewPgStore(t.db.Client.DB.DB)
	bgjobCli := bgjob.NewClient(bgjobDb)
	tgBot, _ := telegramt.TestBot(test)
	cfg := getConfig()
	locatorCfg, err := assembly.Locator(
		context.Background(),
		test.Logger(),
		t.db.Client,
		nil,
		tgBot,
		bgjobCli,
		cfg,
	)
	t.Require().NoError(err)
	server := httptest.NewServer(locatorCfg.HttpRouter)
	t.cli = client.NewWithClient(server.Client())
	t.cli.GlobalRequestConfig().BaseUrl = fmt.Sprintf("http://%s", server.Listener.Addr().String())

	var userId string
	err = t.db.Get(&userId,
		`INSERT INTO users(username,name,admin)
		VALUES($1,$2,$3)
		RETURNING id;`,
		"@admin",
		"test",
		true,
	)
	t.Require().NoError(err)
	jwtGen, err := jwt.GenerateToken(cfg.Auth.Access.Secret, cfg.Auth.Access.TTL, entity.TokenUserInfo{
		UserId:   userId,
		RoleName: domain.AdminType,
	})
	t.Require().NoError(err)

	t.adminAccessToken = domain.BearerToken + " " + jwtGen.Token

	err = t.db.Get(&userId,
		`INSERT INTO users(username,name,admin)
		VALUES($1,$2,$3)
		RETURNING id;`,
		"@user",
		"test",
		false,
	)
	t.Require().NoError(err)
	jwtGen, err = jwt.GenerateToken(cfg.Auth.Access.Secret, cfg.Auth.Access.TTL, entity.TokenUserInfo{
		UserId:   userId,
		RoleName: domain.UserType,
	})
	t.Require().NoError(err)

	t.userAccessToken = domain.BearerToken + " " + jwtGen.Token

	t.T().Cleanup(func() {
		server.Close()
	})
}

func (t *AuthSuite) Test_HasAdminPrivileges_HappyPath() {
	var resp domain.HasAdminPrivilegesResponse
	_, err := t.cli.Get("/has_admin_privileges").
		Header(domain.AuthHeaderName, t.adminAccessToken).
		StatusCodeToError().
		JsonResponseBody(&resp).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().True(resp.HasAdminPrivileges)

	t.Require().NoError(err)
	_, err = t.cli.Get("/has_admin_privileges").
		Header(domain.AuthHeaderName, t.userAccessToken).
		StatusCodeToError().
		JsonResponseBody(&resp).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().False(resp.HasAdminPrivileges)
}
