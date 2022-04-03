package domain

import (
	"encoding/json"
	"errors"
	"recipe-app/pkg/util/writer"
	"strings"
)

type CreatedObjectView struct {
	ID uint64 `json:"id"`
	writer.ServiceResponse
}

type Locale string

const (
	LocaleKk Locale = "KK"
	LocaleEn Locale = "EN"
	LocaleRu Locale = "RU"
)

func (l Locale) String() string {
	return string(l)
}

var LocaleEnum = struct {
	KK Locale
	EN Locale
	RU Locale
}{KK: LocaleKk, EN: LocaleEn, RU: LocaleRu}

func (l *Locale) As() string {
	return string(*l)
}

type RestCtxKey string

func (k RestCtxKey) String() string {
	return string(k)
}

type Stateable interface {
	StateFieldVal() interface{}
}

// Asable is an alias of String().
type Asable interface {
	As() string
}

type StateAsable interface {
	Stateable
	Asable
}

type StateStr string

const (
	ActiveStr   StateStr = "active"
	DisabledStr StateStr = "disabled"
	DeletedStr  StateStr = "deleted"
)

func (s *StateStr) As() string {
	return string(*s)
}

func (s *StateStr) StateFieldVal() interface{} {
	return s
}

type StateTitleSettable interface {
	SetStateTitle(func(s StateAsable) string)
}

func SetStateTitleEnMasse(f func(s StateAsable) string, settables ...StateTitleSettable) {
	for _, settable := range settables {
		settable.SetStateTitle(f)
	}
}

var ErrUnmatchedState = errors.New("unmatched state")

type State uint8

const (
	active   = "active"
	disabled = "disabled"
	deleted  = "deleted"
)

// It would be nice to consider it as a proper pgtype.
const (
	Active State = iota + 1
	Disabled
	Deleted
)

func (s State) String() string {
	switch s {
	case Active:
		return active
	case Disabled:
		return disabled
	case Deleted:
		return deleted
	default:
		return ""
	}
}

func (s *State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String()) // nolint:wrapcheck // Using MarshalJSON only for state
}

func MatchState(s string) (state State, ok bool) {
	switch strings.ToLower(s) {
	case "", active: // the empty string comes from client, thus I have no idea besides just casting act like that.
		return Active, true
	case disabled:
		return Disabled, true
	case deleted:
		return Deleted, true
	default:
		return State(4), false // nolint:gomnd // Unknown default state
	}
}

func (s *State) UnmarshalJSON(data []byte) error {
	var maybeStr string

	if err := json.Unmarshal(data, &maybeStr); err != nil {
		return err // nolint:wrapcheck // Using UnmarshalJSON only for state
	}

	state, ok := MatchState(maybeStr)
	if !ok {
		return ErrUnmatchedState
	}

	*s = state

	return nil
}

var StateEnum = struct {
	Active   State
	Disabled State
	Deleted  State
}{Active: Active, Disabled: Disabled, Deleted: Deleted}

type DBVMap map[string]interface{}

type Pageable struct {
	Content       interface{} `json:"content"`
	PageNumber    uint64      `json:"number"`
	PageSize      uint64      `json:"number_of_elements"`
	ElementsCount uint64      `json:"total_elements"`
	TotalPages    uint64      `json:"total_pages"`
}

type PageableQueryParams struct {
	Page *uint64 `schema:"page"`
	Size *uint64 `schema:"size"`
}

type CreatedUpdatedDate struct {
	CreatedDate DBTime `json:"created_date"`
	UpdatedDate DBTime `json:"updated_date"`
}

type PickingTaskState string

const (
	New       PickingTaskState = "new"
	Suspended PickingTaskState = "suspended"
)
