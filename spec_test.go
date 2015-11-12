package spec_test

import (
	"github.com/lysu/struct-spec"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type AnType struct {
	A string `a:"b"`
}

func TestExtract(t *testing.T) {
	a := AnType{}
	ss := spec.StructSpecForType("a", reflect.TypeOf(a))
	assert.Equal(t, "a", ss.TagName)
	assert.Equal(t, "b", ss.Items[0].Name)
}
