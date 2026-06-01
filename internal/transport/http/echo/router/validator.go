package transport_http_echo_router

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/egorkto/Chat-go/internal/domain"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type RequestValidator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewRequestValidator() (RequestValidator, error) {
	val := getNewValidate()

	trans, err := getNewTranslator()
	if err != nil {
		return RequestValidator{}, fmt.Errorf("get new translator: %w", err)
	}

	enTranslations.RegisterDefaultTranslations(val, trans)

	return RequestValidator{
		validate: val,
		trans:    trans,
	}, nil
}

func (rv RequestValidator) Validate(i any) error {
	err := rv.validate.Struct(i)
	if err == nil {
		return nil
	}

	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		valErrs := make(map[string]string, len(valErr))
		for _, fErr := range valErr {
			valErrs[fErr.Field()] = fErr.Translate(rv.trans)
		}
		return domain.NewValidationError(valErrs)
	}

	return fmt.Errorf("struct: %w", err)
}

func getNewValidate() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}

func getNewTranslator() (ut.Translator, error) {
	uni := ut.New(en.New())
	trans, notFound := uni.GetTranslator("en")
	if notFound {
		return nil, fmt.Errorf("translator not found")
	}

	return trans, nil
}

func (rv RequestValidator) RegisterTranslations(translations map[string]string) {
	for tag, translation := range translations {
		rv.validate.RegisterTranslation(tag, rv.trans,
			func(ut ut.Translator) error {
				return ut.Add(tag, translation, true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(tag, fe.Field(), fe.Param())
				return t
			},
		)
	}
}
