package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"recipe-app/pkg/domain/constant"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	ruTranslations "github.com/go-playground/validator/v10/translations/ru"
)

const (
	splitCount    = 2
	fstArgOfSplit = 0
	sndArgOfSplit = 1
)

func Validate(a interface{}) (msg *string, vMap map[string]string) {
	// I have no idea how to separate the initialization, since it keeps doing fatal panics.
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", splitCount)[fstArgOfSplit]
		if name == "-" {
			return ""
		}

		return name
	})

	uni := ut.New(en.New(), ru.New())

	trans, ok := uni.GetTranslator("ru")
	if !ok {
		log.Printf("couldn't get locale[ru] for validator")
	}

	if err := ruTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Printf("couldn't register locale[ru] for validator")
	}

	vMap = make(map[string]string)

	if err := v.Struct(a); err != nil {
		var valErrs validator.ValidationErrors
		if match := errors.As(err, &valErrs); match {
			for _, e := range valErrs {
				jsonTag := strings.SplitN(e.Namespace(), ".", splitCount)[sndArgOfSplit]
				vMap[jsonTag] = e.Translate(trans)
			}
		}

		msg := constant.MsgRequiredErr

		return &msg, vMap
	}

	return nil, vMap
}
