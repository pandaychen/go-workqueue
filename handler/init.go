package handler

//add validator

import (
	"github.com/go-playground/validator/v10"
)

var g_validator *validator.Validate = validator.New()
