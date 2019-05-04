package internal

import (
	"reflect"
	"strings"
	"sync"

	"github.com/azer/snakecase"
)

type primaryKey interface {
	PrimaryKey() (string, interface{})
}

type primaryKeyData struct {
	field string
	index int
}

var primaryKeysCache sync.Map

func InferPrimaryKey(record interface{}) (string, interface{}) {
	if pk, ok := record.(primaryKey); ok {
		return pk.PrimaryKey()
	}

	rv := reflectInternalStruct(record)
	rt := rv.Type()

	result, cached := primaryKeysCache.Load(rt)
	if !cached {
		field, index := searchPrimaryKey(rt)
		result = primaryKeyData{
			field: field,
			index: index,
		}

		primaryKeysCache.Store(rt, result)
	}

	pkey := result.(primaryKeyData)
	return pkey.field, rv.Field(pkey.index).Interface()
}

func searchPrimaryKey(rt reflect.Type) (string, int) {
	field := ""
	index := 0

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)

		if tag := sf.Tag.Get("db"); strings.HasSuffix(tag, ",primary") {
			index = i
			if len(tag) > 8 { // has custom field name
				field = tag[:len(tag)-8]
			} else {
				field = snakecase.SnakeCase(sf.Name)
			}

			continue
		}

		// check fallback for id field
		if strings.EqualFold("id", sf.Name) {
			index = i
			field = "id"
		}
	}

	if field == "" {
		panic("grimoire: failed to infer primary key for type " + rt.String())
	}

	return field, index
}
