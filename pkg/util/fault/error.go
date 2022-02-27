package fault

import (
	"errors"
	"fmt"
	"net/http"
	"recipe-app/pkg/domain/constant"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type RecipeError struct {
	HTTPStatus int
	Status     int
	Message    string
	Debug      string
	Err        error
	Validation map[string]string
}

func (e *RecipeError) Error() string {
	return e.Message
}

func Whs400Error(s, msg string) error {
	return &RecipeError{ //nolint:exhaustivestruct //no need
		HTTPStatus: http.StatusBadRequest,
		Debug:      s,
		Message:    msg,
	}
}

func Whs404Error(s, msg string) error {
	return &RecipeError{ //nolint:exhaustivestruct //no need
		HTTPStatus: http.StatusNotFound,
		Debug:      s,
		Message:    msg,
	}
}

func Whs409Error(s, msg string) error {
	return &RecipeError{ //nolint:exhaustivestruct //no need
		HTTPStatus: http.StatusConflict,
		Debug:      s,
		Message:    msg,
	}
}

func Whs500Error(s, msg string) error {
	return &RecipeError{ //nolint:exhaustivestruct //no need
		HTTPStatus: http.StatusInternalServerError,
		Debug:      s,
		Message:    msg,
	}
}

func WhsCtxKVError(key string) error {
	return Whs500Error(
		"",
		fmt.Sprintf("couldn't type-assert/parse value from ctx for a key {%s}", key),
	)
}

func WhsValidateError(msg string, vMap map[string]string) error {
	return &RecipeError{ //nolint:exhaustivestruct //no need
		HTTPStatus: http.StatusBadRequest, Message: msg, Validation: vMap,
	}
}

var (
	EmptyArgs              []interface{}
	ErrJuncNotFound        = errors.New("couldn't find junction table")
	ErrNoImplFound         = errors.New("couldn't find implementation")
	ErrNotFoundProductCard = errors.New("не получилось найти продукты в сервисе product-card")
	ErrNotFoundTradeItem   = errors.New("не получилось найти товарные единицы в сервисе trade-item")
)

type QueryStringComposeError struct {
	action constant.SQLAction
	table  constant.Table
	err    error
}

func NewQueryStringComposeError(action constant.SQLAction, table constant.Table, err error) *QueryStringComposeError {
	return &QueryStringComposeError{
		action: action,
		table:  table,
		err:    err,
	}
}

func (e QueryStringComposeError) Error() string {
	return fmt.Sprintf("querystring compose error, action: {%s}, table {%s}, err: {%v}", e.action, e.table, e.err)
}

type DBErrorReason string

const (
	NotFound             DBErrorReason = "entry not found in db"
	ParentNotFound       DBErrorReason = "entry parent not found in db (foreign key violation)"
	AlreadyExists        DBErrorReason = "entry already exists in db"
	NoRowsChanged        DBErrorReason = "no rows changed in db"
	Unhandled            DBErrorReason = "unhandled error from db"
	FailsCheckConstraint DBErrorReason = "fails check constraint in db"
)

// DBRaisedError is a wrapper error for an error raised in DB.
// Reason is a reason of error.
// Stmt is a compiled query.
// Args are objects passed towards pgx call .
// Err is a nested pgconn.PgError object.
// Debug is a debug flag for error.
type DBRaisedError struct {
	Reason DBErrorReason
	Stmt   string
	Args   []interface{}
	Err    error
	Debug  bool
}

func (e DBRaisedError) Error() string {
	if !e.Debug {
		return fmt.Sprintf("%s, stmt: {%s}, args: {%v}", e.Reason, e.Stmt, e.Args)
	}

	return fmt.Sprintf("%s, stmt: {%s}, args: {%v}, err: {%v}", e.Reason, e.Stmt, e.Args, e.Err)
}

func FailCheckInDBError(stmt string, args []interface{}, err error) *DBRaisedError {
	return &DBRaisedError{Reason: FailsCheckConstraint, Stmt: stmt, Args: args, Err: err} //nolint:exhaustivestruct // error response
}

func NotFoundInDBError(stmt string, args []interface{}) *DBRaisedError {
	return &DBRaisedError{Reason: NotFound, Stmt: stmt, Args: args} //nolint:exhaustivestruct // error response
}

func EntryAlreadyExistsInDBError(stmt string, args []interface{}, err error) *DBRaisedError {
	return &DBRaisedError{Reason: AlreadyExists, Stmt: stmt, Args: args, Err: err} //nolint:exhaustivestruct // error response
}

func ParentObjectNotFoundInDBError(stmt string, args []interface{}, err error) *DBRaisedError {
	return &DBRaisedError{Reason: ParentNotFound, Stmt: stmt, Args: args, Err: err} //nolint:exhaustivestruct // error response
}

func NoRowsChangedInDBError(stmt string, args []interface{}) *DBRaisedError {
	return &DBRaisedError{Reason: NoRowsChanged, Stmt: stmt, Args: args} //nolint:exhaustivestruct // error response
}

func UnhandledDBError(stmt string, args []interface{}, err error) *DBRaisedError {
	return &DBRaisedError{Reason: Unhandled, Stmt: stmt, Args: args, Err: err, Debug: true}
}

func SanitizeDBError(err error, stmt string, args []interface{}) error {
	var matchDBErr *DBRaisedError
	if ok := errors.As(err, &matchDBErr); ok {
		return err
	}

	var match *pgconn.PgError
	if ok := errors.As(err, &match); ok {
		switch match.Code {
		case constant.PgUniqueConstraintViolation:
			return EntryAlreadyExistsInDBError(stmt, args, err)
		case constant.PgForeignKeyViolation:
			return ParentObjectNotFoundInDBError(stmt, args, err)
		case constant.PgCheckConstraint:
			return FailCheckInDBError(stmt, args, err)
		case constant.PgRelationDoesNotExist:
			return NotFoundInDBError(stmt, args)
		// Add the other cases as soon as they're raised.
		default:
			return UnhandledDBError(stmt, args, err)
		}
	}

	if match := errors.Is(err, pgx.ErrNoRows); match {
		return NotFoundInDBError(stmt, args)
	}

	return err
}

func SanitizeServiceError(err error) error {
	var expectedErr *DBRaisedError
	if match := errors.As(err, &expectedErr); match {
		switch expectedErr.Reason {
		case NotFound, NoRowsChanged:
			return Whs404Error(err.Error(), constant.MsgNotFoundErr)
		case ParentNotFound, FailsCheckConstraint:
			return Whs400Error(err.Error(), constant.MsgRequestBodyErr)
		case AlreadyExists:
			return Whs409Error(err.Error(), constant.MsgAlreadyExists)
		case Unhandled:
			return Whs500Error(err.Error(), constant.MsgUnhandledErr)
		default:
			return Whs500Error(err.Error(), constant.MsgUnhandledErr)
		}
	}

	var msmtErr *DBRaisedError
	if match := errors.As(err, &msmtErr); match {
		return Whs400Error(err.Error(), constant.MsgRequestBodyErr)
	}

	return Whs500Error(err.Error(), constant.MsgUnhandledErr)
}

type Cmp string

const (
	Less = "less"
	More = "more"
)

type MeasurementKind string

var (
	Length MeasurementKind = "length"
	Width  MeasurementKind = "width"
	Height MeasurementKind = "height"
	Volume MeasurementKind = "volume"
	Weight MeasurementKind = "weight"
)

type MeasurementError struct {
	Kind      MeasurementKind
	Direction Cmp
}

func (e *MeasurementError) Error() string {
	return fmt.Sprintf("%s is %s than needed", e.Kind, e.Direction)
}
