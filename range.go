package kstrct

import (
	"iter"
	"reflect"
	"sync"
)

type FieldCtx struct {
	NumFields int
	Index     int
	Field     reflect.Value
	Name      string
	Value     any
	Type      string
	Tags      []string
}

func (ctx *FieldCtx) Reset() {
	ctx.NumFields = 0
	ctx.Index = 0
	ctx.Field = reflect.Value{}
	ctx.Name = ""
	ctx.Value = nil
	ctx.Type = ""
	ctx.Tags = ctx.Tags[:0]
}

func (ctx *FieldCtx) SetFieldValue(value any) error {
	return TrySet(ctx.Field, value)
}

var fieldCtxPool = sync.Pool{
	New: func() any {
		return &FieldCtx{
			Tags: []string{},
		}
	},
}

// From returns an iterator for struct fields
func From(strctPtr any, tagsToGet ...string) iter.Seq2[int, FieldCtx] {
	return func(yield func(int, FieldCtx) bool) {
		rs := reflect.ValueOf(strctPtr).Elem()
		rt := rs.Type()
		numFields := rs.NumField()

		// Use the same caching mechanism as Fill
		cacheKey := rt.String()
		cache, ok := cacheFieldsIndex.Get(cacheKey)
		if !ok {
			cache = &fieldCache{
				names: make([]string, numFields),
			}
			for i := 0; i < numFields; i++ {
				cache.names[i] = ToSnakeCase(rt.Field(i).Name)
			}
			cacheFieldsIndex.Set(cacheKey, cache)
		}

		// Get a single context to reuse
		ctx := fieldCtxPool.Get().(*FieldCtx)
		defer fieldCtxPool.Put(ctx)

		for i := 0; i < numFields; i++ {
			ctx.Reset()
			f := rs.Field(i)
			val := f.Interface()
			ctx.Field = f
			ctx.Name = cache.names[i] // Use cached name
			ctx.Value = val
			ctx.Type = f.Type().Name()
			ctx.NumFields = numFields
			ctx.Index = i
			ctx.Tags = ctx.Tags[:0]
			for _, t := range tagsToGet {
				if ftag, ok := rt.Field(i).Tag.Lookup(t); ok {
					ctx.Tags = append(ctx.Tags, ftag)
				}
			}

			if !yield(i, *ctx) {
				return
			}
		}
	}
}

// Range iterates over struct fields and calls fn for each field
func Range(strctPtr any, fn func(FieldCtx) bool, tagsToGet ...string) {
	rs := reflect.ValueOf(strctPtr).Elem()
	rt := rs.Type()
	numFields := rs.NumField()

	// Use the same caching mechanism as Fill
	cacheKey := rt.String()
	cache, ok := cacheFieldsIndex.Get(cacheKey)
	if !ok {
		cache = &fieldCache{
			names: make([]string, numFields),
		}
		for i := 0; i < numFields; i++ {
			cache.names[i] = ToSnakeCase(rt.Field(i).Name)
		}
		cacheFieldsIndex.Set(cacheKey, cache)
	}

	// Get a single context to reuse
	ctx := fieldCtxPool.Get().(*FieldCtx)
	defer fieldCtxPool.Put(ctx)

	for i := 0; i < numFields; i++ {
		ctx.Reset()
		f := rs.Field(i)
		val := f.Interface()
		ctx.Field = f
		ctx.Name = cache.names[i] // Use cached name
		ctx.Value = val
		ctx.Type = f.Type().Name()
		ctx.NumFields = numFields
		ctx.Index = i
		ctx.Tags = ctx.Tags[:0]
		for _, t := range tagsToGet {
			if ftag, ok := rt.Field(i).Tag.Lookup(t); ok {
				ctx.Tags = append(ctx.Tags, ftag)
			}
		}

		if !fn(*ctx) {
			return
		}
	}
}
