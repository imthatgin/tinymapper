package tinymapper

import (
	"reflect"

	"github.com/imthatgin/tinymapper/structs"
)

// This is grabbed from https://github.com/geraldo-labs/merge-struct
// and is meant to be a simple way to map fields from A to B.

// patchStruct updates the target struct in-place with non-zero values from the patch struct.
// Only fields with the same name and type get updated. Fields in the patch struct can be
// pointers to the target's type.
// Skips fields with incorrect type, rather than erroring.
// Returns true if any value has been changed.
func patchStruct(target, patch interface{}) {
	dst := structs.New(target)
	fields := structs.New(patch).Fields() // work stack

	for N := len(fields); N > 0; N = len(fields) {
		var srcField = fields[N-1] // pop the top
		fields = fields[:N-1]

		if !srcField.IsExported() {
			continue // skip unexported fields
		}
		if srcField.IsEmbedded() {
			// add the embedded fields into the work stack
			fields = append(fields, srcField.Fields()...)
			continue
		}
		if srcField.IsZero() {
			continue // skip zero-value fields
		}

		var name = srcField.Name()

		var dstField, ok = dst.FieldOk(name)
		if !ok {
			continue // skip non-existing fields
		}

		var srcValue = reflect.ValueOf(srcField.Value())
		srcValue = reflect.Indirect(srcValue)
		// If these are not the same, we skip, and add them to the "skipped list"
		if skind, dkind := srcValue.Kind(), dstField.Kind(); skind != dkind {
			continue
		}

		srcType := reflect.TypeOf(srcValue.Interface())
		dstType := reflect.TypeOf(dstField.Value())
		if srcType != dstType {
			continue
		}

		_ = dstField.Set(srcValue.Interface())
	}
}
