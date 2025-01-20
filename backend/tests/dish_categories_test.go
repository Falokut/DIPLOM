// nolint:noctx,funlen
package tests_test

import (
	"context"
	"dish_as_a_service/assembly"
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"
	"dish_as_a_service/jwt"
	"dish_as_a_service/repository"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/Falokut/go-kit/http/client"
	"github.com/Falokut/go-kit/json"
	"github.com/Falokut/go-kit/test"
	"github.com/Falokut/go-kit/test/dbt"
	"github.com/Falokut/go-kit/test/fake"
	"github.com/Falokut/go-kit/test/telegramt"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/bgjob"
)

type DishCategoriesSuite struct {
	suite.Suite
	test             *test.Test
	adminAccessToken string

	db                 *dbt.TestDb
	dishCategoriesRepo repository.DishesCategories
	dishRepo           repository.Dish
	cli                *client.Client
}

func TestDishCategories(t *testing.T) {
	t.Parallel()
	suite.Run(t, &DishCategoriesSuite{})
}

// nolint:dupl
func (t *DishCategoriesSuite) SetupTest() {
	test, _ := test.New(t.T())
	t.test = test
	t.db = dbt.New(test, db.WithMigrationRunner("../migrations", test.Logger()))
	t.dishCategoriesRepo = repository.NewDishesCategories(t.db.Client)
	t.dishRepo = repository.NewDish(t.db.Client)

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
	t.cli.GlobalRequestConfig().BaseUrl = fmt.Sprintf("http://%s", server.Listener.Addr())

	var userId string
	err = t.db.Get(&userId,
		`INSERT INTO users(username,name,admin)
		VALUES($1,$2,$3)
		RETURNING id;`,
		"@test",
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

	t.T().Cleanup(func() {
		server.Close()
	})
}

func (t *DishCategoriesSuite) Test_GetAllCategories_HappyPath() {
	var categories []domain.DishCategory
	_, err := t.cli.Get("/dishes/all_categories").
		JsonResponseBody(&categories).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().ElementsMatch([]domain.DishCategory{
		{
			Id:   1,
			Name: "Горячее",
		},
		{
			Id:   2,
			Name: "Холодное",
		},
		{
			Id:   3,
			Name: "Напиток",
		},
		{
			Id:   4,
			Name: "Острое",
		},
		{
			Id:   5,
			Name: "Рыба",
		},
		{
			Id:   6,
			Name: "Вегетарианское",
		},
		{
			Id:   7,
			Name: "Мясное",
		},
	}, categories)
}

func (t *DishCategoriesSuite) Test_GetDishesCategories_HappyPath() {
	var categories []domain.DishCategory
	_, err := t.cli.Get("/dishes/categories").
		JsonResponseBody(&categories).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().Empty(categories)

	err = t.dishRepo.AddDish(context.Background(), &entity.AddDishRequest{
		Name:        "dish",
		Description: "desc",
		ImageId:     "image_id",
		Price:       1000,
		Categories:  []int32{1, 2},
	})
	t.Require().NoError(err)

	_, err = t.cli.Get("/dishes/categories").
		JsonResponseBody(&categories).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().ElementsMatch([]domain.DishCategory{
		{
			Id:   1,
			Name: "Горячее",
		},
		{
			Id:   2,
			Name: "Холодное",
		},
	}, categories)
}

func (t *DishCategoriesSuite) Test_GetCategory_HappyPath() {
	const categoryId = 3
	var category domain.DishCategory
	_, err := t.cli.Get(fmt.Sprintf("/dishes/categories/%d", categoryId)).
		JsonResponseBody(&category).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().Equal(domain.DishCategory{Id: categoryId, Name: "Напиток"}, category)
}

func (t *DishCategoriesSuite) Test_AddCategory_HappyPath() {
	categoryName := fake.It[string]()
	req := domain.AddCategoryRequest{Name: categoryName}
	var resp domain.AddCategoryResponse
	_, err := t.cli.Post("/dishes/categories").
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(req).
		JsonResponseBody(&resp).
		Do(context.Background())
	t.Require().NoError(err)

	category, err := t.dishCategoriesRepo.GetCategory(context.Background(), resp.Id)
	t.Require().NoError(err)
	t.Require().Equal(categoryName, category.Name)
}

func (t *DishCategoriesSuite) Test_DeleteCategory_HappyPath() {
	categoryName := fake.It[string]()
	req := domain.AddCategoryRequest{Name: categoryName}
	var resp domain.AddCategoryResponse
	_, err := t.cli.Post("/dishes/categories").
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(req).
		JsonResponseBody(&resp).
		Do(context.Background())
	t.Require().NoError(err)

	category, err := t.dishCategoriesRepo.GetCategory(context.Background(), resp.Id)
	t.Require().NoError(err)
	t.Require().Equal(categoryName, category.Name)

	err = t.cli.Delete(fmt.Sprintf("dishes/categories/%d", resp.Id)).
		Header(domain.AuthHeaderName, t.adminAccessToken).
		DoWithoutResponse(context.Background())
	t.Require().NoError(err)

	getResp, err := t.cli.Get(fmt.Sprintf("/dishes/categories/%d", resp.Id)).
		Do(context.Background())
	t.Require().NoError(err)
	t.Require().NoError(err)
	t.Require().EqualValues(http.StatusNotFound, getResp.StatusCode())

	respBody, err := getResp.Body()
	t.Require().NoError(err)

	var errorResp apierrors.Error
	err = json.Unmarshal(respBody, &errorResp)
	t.Require().NoError(err)
	t.Require().EqualValues(domain.ErrCodeDishCategoryNotFound, errorResp.ErrorCode)
}

func (t *DishCategoriesSuite) Test_RenameCategory_HappyPath() {
	categoryName := fake.It[string]()
	req := domain.AddCategoryRequest{Name: categoryName}
	var resp domain.AddCategoryResponse
	_, err := t.cli.Post("/dishes/categories").
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(req).
		JsonResponseBody(&resp).
		Do(context.Background())
	t.Require().NoError(err)

	category, err := t.dishCategoriesRepo.GetCategory(context.Background(), resp.Id)
	t.Require().NoError(err)
	t.Require().Equal(categoryName, category.Name)

	newCategoryName := fake.It[string]()
	renameReq := domain.RenameCategoryRequest{Name: newCategoryName}
	_, err = t.cli.Post(fmt.Sprintf("/dishes/categories/%d", resp.Id)).
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(renameReq).
		StatusCodeToError().
		Do(context.Background())
	t.Require().NoError(err)

	category, err = t.dishCategoriesRepo.GetCategory(context.Background(), resp.Id)
	t.Require().NoError(err)
	t.Require().Equal(newCategoryName, category.Name)
}

func (t *DishCategoriesSuite) Test_RenameCategory_Conflict() {
	categoryName := fake.It[string]()
	req := domain.AddCategoryRequest{Name: categoryName}
	var resp domain.AddCategoryResponse
	_, err := t.cli.Post("/dishes/categories").
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(req).
		JsonResponseBody(&resp).
		Do(context.Background())
	t.Require().NoError(err)

	category, err := t.dishCategoriesRepo.GetCategory(context.Background(), resp.Id)
	t.Require().NoError(err)
	t.Require().Equal(categoryName, category.Name)

	categoryName = fake.It[string]()
	_, err = t.cli.Post("/dishes/categories").
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(domain.AddCategoryRequest{Name: categoryName}).
		Do(context.Background())
	t.Require().NoError(err)

	renameReq := domain.RenameCategoryRequest{Name: categoryName}
	renameResp, err := t.cli.Post(fmt.Sprintf("/dishes/categories/%d", resp.Id)).
		Header(domain.AuthHeaderName, t.adminAccessToken).
		JsonRequestBody(renameReq).
		Do(context.Background())
	t.Require().NoError(err)

	_, err = t.dishCategoriesRepo.GetCategory(context.Background(), resp.Id)
	t.Require().NoError(err)
	t.Require().EqualValues(http.StatusConflict, renameResp.StatusCode())
}
