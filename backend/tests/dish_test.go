package tests_test

import (
	"net/http"
	"testing"

	"github.com/Falokut/go-kit/test"
	"github.com/Falokut/go-kit/test/dbt"
	"github.com/stretchr/testify/suite"
)

type DishSuite struct {
	suite.Suite
	test *test.Test

	db  *dbt.TestDb
	cli http.Client
}

func TestDish(t *testing.T) {
	t.Parallel()
	suite.Run(t, &DishSuite{})
}

func (t *DishSuite) SetupTest() {
	test, _ := test.New(t.T())
	t.test = test
}
