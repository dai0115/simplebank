package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/simplebank/util"
)

var validCurrency validator.Func = func(fieldlevel validator.FieldLevel) bool {
	// check whether currency field is string
	if currency, ok := fieldlevel.Field().Interface().(string); ok {
		// check currency is supported
		return util.IsSupportedCurrency(currency)
	}
	return false
}
