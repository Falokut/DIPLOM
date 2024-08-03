package routes

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Falokut/go-kit/encoding/json"
	"github.com/labstack/echo/v4"
)

type JsonSerializer struct{}

func (s JsonSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

func (s JsonSerializer) Deserialize(c echo.Context, i interface{}) error {
	err := json.NewDecoder(c.Request().Body).Decode(i)
	if err == nil {
		return nil
	}

	ute := &json2.UnmarshalTypeError{}
	se := &json2.SyntaxError{}
	switch {
	case errors.As(err, &ute):
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v",
				ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	case errors.As(err, &se):
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Syntax error: offset=%v, error=%v",
				se.Offset, se.Error())).SetInternal(err)
	}

	return err
}
