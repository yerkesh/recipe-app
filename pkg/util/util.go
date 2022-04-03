package util

import (
	"fmt"
	"net/http"
	"recipe-app/pkg/domain"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

const (
	numBase       = 10
	bitSize       = 64
	splitCount    = 2
	fstArgOfSplit = 0
)

func CurUtcTime() time.Time {
	return time.Now().UTC()
}

func CurTime() time.Time {
	return time.Now()
}

func ParseUint64(id string) (parsedU uint64, err error) {
	if parsedU, err = strconv.ParseUint(id, numBase, bitSize); err != nil {
		return 0, fmt.Errorf("coudln't parse string {%s}, since err: {%w}", id, err)
	}

	return parsedU, nil
}

func ParseBool(str string) (parsedU bool, err error) {
	if parsedU, err = strconv.ParseBool(str); err != nil {
		return false, fmt.Errorf("coudln't parse string {%s}, since err: {%w}", str, err)
	}

	return parsedU, nil
}

type TagPolicy func(f *reflect.StructField) (tag string, ok bool)

func QueryParamTagPolicy(f *reflect.StructField) (tag string, ok bool) {
	schema := func(pre string) (string, bool) {
		tag := splitter(f.Tag.Get("schema"))
		if tag == "" {
			return pre + strcase.ToSnake(f.Name), true
		}

		return pre + tag, true
	}

	if pfx := splitter(f.Tag.Get("pfx")); pfx != "" {
		return schema(pfx + ".")
	}

	return schema("")
}

func splitter(tag string) string {
	return strings.SplitN(tag, ",", splitCount)[fstArgOfSplit]
}

func NonNilFieldsAsMap(ptrStruct interface{}, policy TagPolicy) (kvps domain.DBVMap) {
	if maybePtr := reflect.ValueOf(ptrStruct); maybePtr.Kind() == reflect.Ptr && maybePtr.IsNil() {
		return kvps
	}

	kvps = make(domain.DBVMap)
	nonNilStructVal := reflect.Indirect(reflect.ValueOf(ptrStruct)).Interface()

	structVal, structFields := reflect.ValueOf(nonNilStructVal), reflect.TypeOf(nonNilStructVal)
	for i := 0; i < structVal.NumField(); i++ {
		field := structFields.Field(i)

		gotTag, ok := policy(&field)
		if !ok {
			continue
		}

		fVal := structVal.Field(i)
		if fVal.Kind() == reflect.Ptr && fVal.IsNil() {
			continue
		}

		setVal := reflect.Indirect(fVal).Interface()
		switch reflect.ValueOf(setVal).Kind() { //nolint:exhaustive // no need to case switch on the rest of kinds
		case reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Ptr,
			reflect.Slice,
			reflect.Struct,
			reflect.UnsafePointer:
			continue
		}

		kvps[gotTag] = reflect.Indirect(fVal).Interface()
	}

	return kvps
}

func SetDefaultSizePQPIfNil(pqp *domain.PageableQueryParams) {
	if pqp.Size == nil {
		var Ten uint64 = 10

		pqp.Size = &Ten
	}

	if pqp.Page == nil || *(pqp.Page) == 0 {
		var One uint64 = 1

		pqp.Page = &One
	}
}

func TotalPageCounter(res *domain.Pageable) {
	res.TotalPages = res.ElementsCount / res.PageSize

	if res.ElementsCount%res.PageSize != 0 {
		res.TotalPages++
	}
}

func LocaleFromCookie(cook *http.Cookie) *domain.Locale {
	if cook == nil {
		return &domain.LocaleEnum.RU
	}

	switch cook.Value {
	case domain.LocaleEnum.EN.String():
		return &domain.LocaleEnum.EN
	case domain.LocaleEnum.KK.String():
		return &domain.LocaleEnum.KK
	default:
		return &domain.LocaleEnum.RU
	}
}

func LocaleNilFromCookie(cook *http.Cookie) *domain.Locale {
	switch {
	case cook == nil:
		return nil
	case cook.Value == domain.LocaleEnum.EN.As():
		return &domain.LocaleEnum.EN
	case cook.Value == domain.LocaleEnum.KK.As():
		return &domain.LocaleEnum.KK
	case cook.Value == domain.LocaleEnum.RU.As():
		return &domain.LocaleEnum.RU
	default:
		return nil
	}
}

func ConvertUintToString(un uint64) string {
	return strconv.FormatUint(un, numBase)
}
