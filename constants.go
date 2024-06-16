// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package entrest

import (
	"fmt"

	"entgo.io/ent/entc/gen"
	"github.com/go-openapi/inflect"
	"github.com/stoewer/go-strcase"
)

const OpenAPIVersion = "3.0.3"

var (
	// Add all casing and word-massaging functions here so others can use them if they
	// want to customize the naming of their spec/endpoints/etc.
	Pluralize   = inflect.Pluralize // nolint: errcheck,unused
	KebabCase   = strcase.KebabCase
	Singularize = gen.Funcs["singular"].(func(string) string) // nolint: errcheck,unused
	PascalCase  = gen.Funcs["pascal"].(func(string) string)   // nolint: errcheck,unused
	CamelCase   = gen.Funcs["camel"].(func(string) string)    // nolint: errcheck,unused
	SnakeCase   = gen.Funcs["snake"].(func(string) string)    // nolint: errcheck,unused
)

// Operation represents the CRUD operation(s).
type Operation string

const (
	// OperationCreate represents the create operation (method: POST).
	OperationCreate Operation = "create"
	// OperationRead represents the read operation (method: GET).
	OperationRead Operation = "read"
	// OperationUpdate represents the update operation (method: PATCH).
	OperationUpdate Operation = "update"
	// OperationDelete represents the delete operation (method: DELETE).
	OperationDelete Operation = "delete"
	// OperationList represents the list operation (method: GET).
	OperationList Operation = "list"
)

// AllOperations holds a list of all supported operations.
var AllOperations = []Operation{OperationCreate, OperationRead, OperationUpdate, OperationDelete, OperationList}

// Predicate represents a filtering predicate provided by ent.
type Predicate int

// Mirrored from entgo.io/ent/entc/gen with special groupings and support for bitwise operations.
const (
	// FilterEdge is a special filter which is applied to the edge itself, indicating
	// that all of the edges fields should also be included in filtering options.
	FilterEdge Predicate = 1 << iota

	FilterEQ           // =
	FilterNEQ          // <>
	FilterGT           // >
	FilterGTE          // >=
	FilterLT           // <
	FilterLTE          // <=
	FilterIsNil        // IS NULL / has
	FilterNotNil       // IS NOT NULL / hasNot
	FilterIn           // within
	FilterNotIn        // without
	FilterEqualFold    // equals case-insensitive
	FilterContains     // containing
	FilterContainsFold // containing case-insensitive
	FilterHasPrefix    // startingWith
	FilterHasSuffix    // endingWith

	// FilterGroupEqualExact includes: eq, neq, equal fold, is nil.
	FilterGroupEqualExact = FilterEQ | FilterNEQ | FilterEqualFold | FilterGroupNil
	// FilterGroupEqual includes: eq, neq, equal fold, contains, contains case, prefix, suffix, nil.
	FilterGroupEqual = FilterGroupEqualExact | FilterGroupContains | FilterHasPrefix | FilterHasSuffix
	// FilterGroupContains includes: contains, contains case, is nil.
	FilterGroupContains = FilterContains | FilterContainsFold | FilterGroupNil
	// FilterGroupNil includes: is nil.
	FilterGroupNil = FilterIsNil
	// FilterGroupLength includes: gt, lt (often gte/lte isn't really needed).
	FilterGroupLength = FilterGT | FilterLT
	// FilterGroupArray includes: in, not in.
	FilterGroupArray = FilterIn | FilterNotIn
)

// filterMap maps a predicate to the entgo.io/ent/entc/gen.Op (to get string representation).
var filterMap = map[Predicate]gen.Op{
	FilterEQ:           gen.EQ,
	FilterNEQ:          gen.NEQ,
	FilterGT:           gen.GT,
	FilterGTE:          gen.GTE,
	FilterLT:           gen.LT,
	FilterLTE:          gen.LTE,
	FilterIsNil:        gen.IsNil,
	FilterNotNil:       gen.NotNil,
	FilterIn:           gen.In,
	FilterNotIn:        gen.NotIn,
	FilterEqualFold:    gen.EqualFold,
	FilterContains:     gen.Contains,
	FilterContainsFold: gen.ContainsFold,
	FilterHasPrefix:    gen.HasPrefix,
	FilterHasSuffix:    gen.HasSuffix,
}

// String returns the gen.Op string representation of a predicate.
func (p Predicate) String() string {
	if _, ok := filterMap[p]; ok {
		return filterMap[p].Name()
	}
	panic("predicate.String() called with grouped predicate, use Explode() first")
}

// Has returns if the predicate has the provided predicate.
func (p Predicate) Has(v Predicate) bool {
	return p&v != 0
}

// Add adds the provided predicate to the current predicate.
func (p Predicate) Add(v Predicate) Predicate {
	p |= v
	return p
}

// Remove removes the provided predicate from the current predicate.
func (p Predicate) Remove(v Predicate) Predicate {
	p &^= v
	return p
}

// Explode returns all individual predicates as []gen.Op.
func (p Predicate) Explode() (ops []gen.Op) {
	for pred, op := range filterMap {
		if p.Has(pred) {
			ops = append(ops, op)
		}
	}
	return ops
}

// predicateFormat returns the query string representation of a filter predicate.
func predicateFormat(op gen.Op) string {
	switch op {
	case gen.IsNil:
		return "null"
	case gen.HasPrefix:
		return "prefix"
	case gen.HasSuffix:
		return "suffix"
	case gen.EqualFold:
		return "eqFold"
	default:
		return CamelCase(SnakeCase(op.Name()))
	}
}

// predicateDescription returns the description of a filter predicate.
func predicateDescription(f *gen.Field, op gen.Op) string {
	switch op {
	case gen.EQ:
		return fmt.Sprintf("Filters field %q to be equal to the provided value.", f.Name)
	case gen.NEQ:
		return fmt.Sprintf("Filters field %q to be not equal to the provided value.", f.Name)
	case gen.GT:
		if f.IsString() {
			return fmt.Sprintf("Filters field %q to be longer than the provided value.", f.Name)
		}
		return fmt.Sprintf("Filters field %q to be greater than the provided value.", f.Name)
	case gen.GTE:
		if f.IsString() {
			return fmt.Sprintf("Filters field %q to be longer than or equal in length to the provided value.", f.Name)
		}
		return fmt.Sprintf("Filters field %q to be greater than or equal to the provided value.", f.Name)
	case gen.LT:
		if f.IsString() {
			return fmt.Sprintf("Filters field %q to be shorter than the provided value.", f.Name)
		}
		return fmt.Sprintf("Filters field %q to be less than the provided value.", f.Name)
	case gen.LTE:
		if f.IsString() {
			return fmt.Sprintf("Filters field %q to be shorter than or equal in length to the provided value.", f.Name)
		}
		return fmt.Sprintf("Filters field %q to be less than or equal to the provided value.", f.Name)
	case gen.IsNil:
		return fmt.Sprintf("Filters field %q to be null/nil.", f.Name)
	case gen.NotNil:
		return fmt.Sprintf("Filters field %q to be not null/nil.", f.Name)
	case gen.In:
		return fmt.Sprintf("Filters field %q to be within the provided values.", f.Name)
	case gen.NotIn:
		return fmt.Sprintf("Filters field %q to be not within the provided values.", f.Name)
	case gen.EqualFold:
		return fmt.Sprintf("Filters field %q to be equal to the provided value, case-insensitive.", f.Name)
	case gen.Contains:
		return fmt.Sprintf("Filters field %q to contain the provided value.", f.Name)
	case gen.ContainsFold:
		return fmt.Sprintf("Filters field %q to contain the provided value, case-insensitive.", f.Name)
	case gen.HasPrefix:
		return fmt.Sprintf("Filters field %q to start with the provided value.", f.Name)
	case gen.HasSuffix:
		return fmt.Sprintf("Filters field %q to end with the provided value.", f.Name)
	default:
		panic("unknown predicate")
	}
}

const (
	defaultMinItemsPerPage = 1
	defaultMaxItemsPerPage = 100
	defaultItemsPerPage    = 10
)

type SupportedHTTPHandler string

const (
	HandlerNone    SupportedHTTPHandler = ""
	HandlerGeneric SupportedHTTPHandler = "generic"
	HandlerChi     SupportedHTTPHandler = "chi"
)

var AllSupportedHTTPHandlers = []SupportedHTTPHandler{
	HandlerNone,
	HandlerGeneric,
	HandlerChi,
}
