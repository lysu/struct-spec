package spec

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

var (
	structSpecMutex sync.RWMutex
	structSpecCache = make(map[reflect.Type]*StructSpec)
)

// FieldSpec present file name and reflect index
type FieldSpec struct {
	Name  string
	Index []int
}

// StructSpec present files that be tagged
type StructSpec struct {
	TagName string
	Index   map[string]*FieldSpec
	Items   []*FieldSpec
}

// FieldSpec use to take fieldSpec by tagged-name
func (ss *StructSpec) FieldSpec(name string) *FieldSpec {
	return ss.Index[name]
}

// StructSpecForType use to extract struct spec from type.
func StructSpecForType(tagName string, t reflect.Type) *StructSpec {

	structSpecMutex.RLock()
	ss, found := structSpecCache[t]
	structSpecMutex.RUnlock()
	if found {
		return ss
	}

	structSpecMutex.Lock()
	defer structSpecMutex.Unlock()
	ss, found = structSpecCache[t]
	if found {
		return ss
	}

	ss = &StructSpec{Index: make(map[string]*FieldSpec), TagName: tagName}
	compileStructSpec(tagName, t, make(map[string]int), nil, ss)
	structSpecCache[t] = ss
	return ss
}

func compileStructSpec(tagName string, t reflect.Type, depth map[string]int, index []int, ss *StructSpec) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		switch {
		case f.PkgPath != "":
		// Ignore unexported fields.
		case f.Anonymous:
			// protection against infinite recursion.
			if f.Type.Kind() == reflect.Struct {
				compileStructSpec(tagName, f.Type, depth, append(index, i), ss)
			}
		default:
			fs := &FieldSpec{Name: f.Name}
			tag := f.Tag.Get(tagName)
			p := strings.Split(tag, ",")
			if len(p) > 0 {
				if p[0] == "-" {
					continue
				}
				if len(p[0]) > 0 {
					fs.Name = p[0]
				}
				for _, s := range p[1:] {
					switch s {
					//case "omitempty":
					//  fs.omitempty = true
					default:
						panic(errors.New("Unknown field flag " + s + " for type " + t.Name()))
					}
				}
			}
			d, found := depth[fs.Name]
			if !found {
				d = 1 << 30
			}
			switch {
			case len(index) == d:
				// At same depth, remove from result.
				delete(ss.Index, fs.Name)
				j := 0
				for i := 0; i < len(ss.Items); i++ {
					if fs.Name != ss.Items[i].Name {
						ss.Items[j] = ss.Items[i]
						j++
					}
				}
				ss.Items = ss.Items[:j]
			case len(index) < d:
				fs.Index = make([]int, len(index)+1)
				copy(fs.Index, index)
				fs.Index[len(index)] = i
				depth[fs.Name] = len(index)
				ss.Index[fs.Name] = fs
				ss.Items = append(ss.Items, fs)
			}
		}
	}
}
