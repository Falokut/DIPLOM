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

func stringToIntSlice(s string) ([]int32, error) {
	if s == "" {
		return []int32{}, nil
	}
	strList := strings.Split(s, ",")
	res := make([]int32, len(strList))

	for i, str := range strList {
		num, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, errors.WithMessage(err, "parse num")
		}
		res[i] = int32(num)
	}
	return res, nil
}
