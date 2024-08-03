package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

//nolint:ireturn
func bindRequest[T any](ctx echo.Context) (T, error) {
	var req T
	err := ctx.Bind(&req)
	if err != nil {
		_ = ctx.NoContent(http.StatusBadRequest)
		return req, err
	}

	err = ctx.Validate(req)
	if err != nil {
		_ = ctx.String(http.StatusBadRequest, err.Error())
		return req, err
	}

	return req, nil
}

func stringToIntSlice(str string) ([]int32, error) {
	if str == "" {
		return []int32{}, nil
	}
	strs := strings.Split(str, ",")
	res := make([]int32, len(strs))

	for i, str := range strs {
		num, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, errors.WithMessage(err, "parse num")
		}
		res[i] = int32(num)
	}
	return res, nil
}
