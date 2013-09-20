/*
 * Copyright (c) 2013 Conformal Systems <info@conformal.com>
 *
 * This file originated from: http://opensource.conformal.com/
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

/*
Go bindings for GLib 2.  Supports version 2.36 and later.
*/
package glib

// #cgo pkg-config: glib-2.0 gobject-2.0
// #include <glib.h>
// #include <glib-object.h>
// #include "glib.go.h"
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

func init() {
	// call g_type_init() if the version is lower than 2.36
	c := C.glib_check_version(C.guint(2), C.guint(36), C.guint(0))
	if c != nil {
		C.g_type_init()
	}
	closures.m = make(map[*C.GClosure]reflect.Value)
}

/*
 * Type conversions
 */

func gbool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}
func gobool(b C.gboolean) bool {
	if b != 0 {
		return true
	}
	return false
}

/*
 * Unexported vars
 */

var (
	nilPtrErr = errors.New("cgo returned unexpected nil pointer")
	closures = struct {
		sync.RWMutex
		m map[*C.GClosure]reflect.Value
	}{}
	signals = make(map[SignalHandle]*C.GClosure)
)

/*
 * closureNew() creates a new GClosure and adds its callback function
 * to the internally-maintained map.
 */
func closureNew(f interface{}) *C.GClosure {
	closure := C._g_closure_new()
	closures.Lock()
	closures.m[closure] = reflect.ValueOf(f)
	closures.Unlock()
	return closure
}

/*
 * Constants
 */

// Type is a representation of GLib's GType.
type Type uint

const (
	TYPE_INVALID   Type = C.G_TYPE_INVALID
	TYPE_NONE           = C.G_TYPE_NONE
	TYPE_INTERFACE      = C.G_TYPE_INTERFACE
	TYPE_CHAR           = C.G_TYPE_CHAR
	TYPE_UCHAR          = C.G_TYPE_UCHAR
	TYPE_BOOLEAN        = C.G_TYPE_BOOLEAN
	TYPE_INT            = C.G_TYPE_INT
	TYPE_UINT           = C.G_TYPE_UINT
	TYPE_LONG           = C.G_TYPE_LONG
	TYPE_ULONG          = C.G_TYPE_ULONG
	TYPE_INT64          = C.G_TYPE_INT64
	TYPE_UINT64         = C.G_TYPE_UINT64
	TYPE_ENUM           = C.G_TYPE_ENUM
	TYPE_FLAGS          = C.G_TYPE_FLAGS
	TYPE_FLOAT          = C.G_TYPE_FLOAT
	TYPE_DOUBLE         = C.G_TYPE_DOUBLE
	TYPE_STRING         = C.G_TYPE_STRING
	TYPE_POINTER        = C.G_TYPE_POINTER
	TYPE_BOXED          = C.G_TYPE_BOXED
	TYPE_PARAM          = C.G_TYPE_PARAM
	TYPE_OBJECT         = C.G_TYPE_OBJECT
	TYPE_VARIANT        = C.G_TYPE_VARIANT
)

func (t Type) Name() string {
	return C.GoString((*C.char)(C.g_type_name(C.GType(t))))
}

// UserDirectory is a representation of GLib's GUserDirectory.
type UserDirectory int

const (
	USER_DIRECTORY_DESKTOP      UserDirectory = C.G_USER_DIRECTORY_DESKTOP
	USER_DIRECTORY_DOCUMENTS                  = C.G_USER_DIRECTORY_DOCUMENTS
	USER_DIRECTORY_DOWNLOAD                   = C.G_USER_DIRECTORY_DOWNLOAD
	USER_DIRECTORY_MUSIC                      = C.G_USER_DIRECTORY_MUSIC
	USER_DIRECTORY_PICTURES                   = C.G_USER_DIRECTORY_PICTURES
	USER_DIRECTORY_PUBLIC_SHARE               = C.G_USER_DIRECTORY_PUBLIC_SHARE
	USER_DIRECTORY_TEMPLATES                  = C.G_USER_DIRECTORY_TEMPLATES
	USER_DIRECTORY_VIDEOS                     = C.G_USER_DIRECTORY_VIDEOS
)

const USER_N_DIRECTORIES int = C.G_USER_N_DIRECTORIES

/*
 * Events
 */

type SignalHandle uint64

// Connect() is a wrapper around g_signal_connect_closure().
func (v *Object) Connect(detailed_signal string, f interface{}) SignalHandle {
	cstr := C.CString(detailed_signal)
	defer C.free(unsafe.Pointer(cstr))
	closure := closureNew(f)
	c := C.g_signal_connect_closure(C.gpointer(v.Native()), (*C.gchar)(cstr), closure, gbool(false))
	h := SignalHandle(c)
	return h
}

// goMarshal() is called by the GLib runtime when a closure needs to be invoked.
// The closure will be invoked with as many arguments as it can take, from 0 to
// the full amount provided by the call. If the closure asks for more parameters
// than there are to give, this method panics.
//
//export goMarshal
func goMarshal(closure *C.GClosure, return_value *C.GValue, n_param_values C.guint, param_values *C.GValue, invocation_hint C.gpointer, marshal_data C.gpointer) {
	var (
		go_params []reflect.Value
		callback = closures.m[closure]
		numIn = callback.Type().NumIn()
		numParams = int(n_param_values)
		ret []reflect.Value
	)
	if numIn == 0 {
		go_params = make([]reflect.Value, 0)
		ret = callback.Call(go_params)
	} else if numIn <= numParams {
		params := valueSlice(numParams, param_values)
		go_params = make([]reflect.Value, numIn)
		for i := 0; i<numIn; i++ {
			v := &Value{*params[i]}
			val, err := v.GoValue()
			if err != nil {
				panic(err)
			}
			go_params[i] = reflect.ValueOf(val)
		}
		ret = callback.Call(go_params)
	} else {
		panic(fmt.Sprintf("not enough arguments to call closure; it expects %d, but we only have %d", numIn, numParams))
	}
	if return_value != nil && len(ret) > 0 {
		g, err := GValue(ret[0].Interface())
		if err != nil {
			panic(err)
		}
		(*return_value) = *g.Native()
	}
}

/*
 * Source support
 */

type Source struct {
	ptr unsafe.Pointer
}

type SourceHandle uint

// Native() returns a pointer to the underlying GSource.
func (v *Source) Native() *C.GSource {
	if v == nil || v.ptr == nil {
		return nil
	}
	return (*C.GSource)(v.ptr)
}

// IdleAdd() adds an idle source to the default main context.
func IdleAdd(f func() bool) (SourceHandle, error) {
	return idleAdd(nil, f)
}

// idleAdd() adds an idle source to the provided main context. If the
// function returns false, then it is invalidated, which should also free it.
func idleAdd(context *MainContext, f func() bool) (SourceHandle, error) {
	c := C.g_idle_source_new()
	if c == nil {
		return 0, nilPtrErr
	}
	var ctx *C.GMainContext = nil
	if context != nil {
		ctx = (*C.GMainContext)(context.ptr)
	}
	var closure *C.GClosure
	closure = closureNew(func() bool {
		ok := f()
		if !ok {
			C.g_closure_invalidate(closure)
		}
		return ok
	})
	C.g_source_set_closure(c, closure)
	cid := C.g_source_attach(c, ctx)
	return SourceHandle(cid), nil
}

/*
 * Main event loop
 */

type MainContext struct {
	ptr unsafe.Pointer
}

func (v *MainContext) IdleAdd(f func() bool) (SourceHandle, error) {
	return idleAdd(v, f)
}

type MainLoop struct {
	ptr unsafe.Pointer
}

func MainLoopNew(context *MainContext) (*MainLoop, error) {
	var ctx *C.GMainContext = nil
	if context != nil {
		ctx = (*C.GMainContext)(context.ptr)
	}
	c := C.g_main_loop_new(ctx, gbool(true))
	if c == nil {
		return nil, nilPtrErr
	}
	return &MainLoop{unsafe.Pointer(c)}, nil
}

// Native() returns a pointer to the underlying GMainLoop.
func (v *MainLoop) Native() *C.GMainLoop {
	if v == nil || v.ptr == nil {
		return nil
	}
	return (*C.GMainLoop)(v.ptr)
}

func (v *MainLoop) Run() {
	C.g_main_loop_run(v.Native())
}

func (v *MainLoop) Quit() {
	C.g_main_loop_quit(v.Native())
}

func (v *MainLoop) Context() *MainContext {
	return &MainContext{unsafe.Pointer(C.g_main_loop_get_context(v.Native()))}
}

/*
 * Miscellaneous Utility Functions
 */

// GetUserSpecialDir() is a wrapper around g_get_user_special_dir().  A
// non-nil error is returned in the case that g_get_user_special_dir()
// returns NULL to differentiate between NULL and an empty string.
func GetUserSpecialDir(directory UserDirectory) (string, error) {
	c := C.g_get_user_special_dir(C.GUserDirectory(directory))
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

/*
 * GObject
 */

// IObject is an interface type implemented by Object and all types which embed
// an Object.  It is meant to be used as a type for function arguments which
// require GObjects or any subclasses thereof.
type IObject interface {
	toGObject() *C.GObject
	ToObject() *Object
}

// Object is a representation of GLib's GObject.
type Object struct {
	ptr unsafe.Pointer
}

func ObjectNew(p unsafe.Pointer) *Object {
	return &Object{p}
}

func (v *Object) Ptr() unsafe.Pointer {
	return v.ptr
}

// Native() returns a pointer to the underlying GObject.
func (v *Object) Native() *C.GObject {
	if v == nil || v.ptr == nil {
		return nil
	}
	return (*C.GObject)(v.ptr)
}

func (v *Object) ToObject() *Object {
	return v
}

func (v *Object) toGObject() *C.GObject {
	if v == nil {
		return nil
	}
	return v.Native()
}

func (v *Object) typeFromInstance() Type {
	c := C._g_type_from_instance(C.gpointer(v.ptr))
	return Type(c)
}

// ToGObject() type converts an unsafe.Pointer as a native C GObject.
// This function is exported for visibility in other gotk3 packages and
// is not meant to be used by applications.
func ToGObject(p unsafe.Pointer) *C.GObject {
	return (*C.GObject)(p)
}

// Ref() is a wrapper around g_object_ref().
func (v *Object) Ref() {
	C.g_object_ref(C.gpointer(v.ptr))
}

// Unref() is a wrapper around g_object_unref().
func (v *Object) Unref() {
	C.g_object_unref(C.gpointer(v.ptr))
}

// RefSink() is a wrapper around g_object_ref_sink().
func (v *Object) RefSink() {
	C.g_object_ref_sink(C.gpointer(v.ptr))
}

// IsFloating() is a wrapper around g_object_is_floating().
func (v *Object) IsFloating() bool {
	c := C.g_object_is_floating(C.gpointer(v.ptr))
	return gobool(c)
}

// ForceFloating() is a wrapper around g_object_force_floating().
func (v *Object) ForceFloating() {
	C.g_object_force_floating((*C.GObject)(v.ptr))
}

// StopEmission() is a wrapper around g_signal_stop_emission_by_name().
func (v *Object) StopEmission(s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.g_signal_stop_emission_by_name(C.gpointer(v.ptr),
		(*C.gchar)(cstr))
}

// Set() is a wrapper around g_object_set().  However, unlike
// g_object_set(), this function only sets one name value pair.  Make
// multiple calls to this function to set multiple properties.
func (v *Object) Set(name string, value interface{}) error {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))

	if _, ok := value.(Object); ok {
		value = value.(Object).ptr
	}

	var p unsafe.Pointer = nil
	switch value.(type) {
	case bool:
		c := gbool(value.(bool))
		p = unsafe.Pointer(&c)
	case int8:
		c := C.gint8(value.(int8))
		p = unsafe.Pointer(&c)
	case int16:
		c := C.gint16(value.(int16))
		p = unsafe.Pointer(&c)
	case int32:
		c := C.gint32(value.(int32))
		p = unsafe.Pointer(&c)
	case int64:
		c := C.gint64(value.(int64))
		p = unsafe.Pointer(&c)
	case int:
		c := C.gint(value.(int))
		p = unsafe.Pointer(&c)
	case uint8:
		c := C.guchar(value.(uint8))
		p = unsafe.Pointer(&c)
	case uint16:
		c := C.guint16(value.(uint16))
		p = unsafe.Pointer(&c)
	case uint32:
		c := C.guint32(value.(uint32))
		p = unsafe.Pointer(&c)
	case uint64:
		c := C.guint64(value.(uint64))
		p = unsafe.Pointer(&c)
	case uint:
		c := C.guint(value.(uint))
		p = unsafe.Pointer(&c)
	case uintptr:
		p = unsafe.Pointer(C.gpointer(value.(uintptr)))
	case float32:
		c := C.gfloat(value.(float32))
		p = unsafe.Pointer(&c)
	case float64:
		c := C.gdouble(value.(float64))
		p = unsafe.Pointer(&c)
	case string:
		cstr := C.CString(value.(string))
		defer C.free(unsafe.Pointer(cstr))
		p = unsafe.Pointer(cstr)
	default:
		if pv, ok := value.(unsafe.Pointer); ok {
			p = pv
		} else {
			// Constants with separate types are not type asserted
			// above, so do a runtime check here instead.
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16,
				reflect.Int32, reflect.Int64:
				c := C.int(val.Int())
				p = unsafe.Pointer(&c)
			case reflect.Uintptr:
				p = unsafe.Pointer(C.gpointer(val.Pointer()))
			}
		}
	}
	// Can't call g_object_set() as it uses a variable arg list, use a
	// wrapper instead
	if p != nil {
		C._g_object_set_one(C.gpointer(v.ptr), (*C.gchar)(cstr), p)
		return nil
	} else {
		return errors.New("Unable to perform type conversion")
	}
}

/*
 * GObject Signals
 */

// Emit() is a wrapper around g_signal_emitv() and emits the signal
// specified by the string s to an Object.  Arguments to callback
// functions connected to this signal must be specified in args.  Emit()
// returns an interface{} which must be type asserted as the Go
// equivalent type to the return value for native C callback.
//
// Note that this code is unsafe in that the types of values in args are
// not checked against whether they are suitable for the callback.
func (v *Object) Emit(s string, args ...interface{}) (interface{}, error) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	// Create array of this instance and arguments
	valv := C.alloc_gvalue_list(C.int(len(args)) + 1)
	defer C.free(unsafe.Pointer(valv))

	// Add args and valv
	val, err := GValue(v)
	if err != nil {
		return nil, errors.New("Error converting Object to GValue: " + err.Error())
	}
	C.val_list_insert(valv, C.int(0), val.Native())
	for i := range args {
		val, err := GValue(args[i])
		if err != nil {
			return nil, fmt.Errorf("Error converting arg %d to GValue: %s", i, err.Error())
		}
		C.val_list_insert(valv, C.int(i+1), val.Native())
	}

	t := v.typeFromInstance()
	id := C.g_signal_lookup((*C.gchar)(cstr), C.GType(t))

	ret, err := ValueAlloc()
	if err != nil {
		return nil, errors.New("Error creating Value for return value")
	}
	C.g_signal_emitv(valv, id, C.GQuark(0), ret.Native())

	return ret.GoValue()
}

// HandlerBlock() is a wrapper around g_signal_handler_block().
func (v *Object) HandlerBlock(handle SignalHandle) {
	C.g_signal_handler_block(C.gpointer(v.ptr), C.gulong(handle))
}

// HandlerUnblock() is a wrapper around g_signal_handler_unblock().
func (v *Object) HandlerUnblock(handle SignalHandle) {
	C.g_signal_handler_unblock(C.gpointer(v.ptr), C.gulong(handle))
}

// HandlerDisconnect() is a wrapper around g_signal_handler_disconnect().
func (v *Object) HandlerDisconnect(handle SignalHandle) {
	C.g_signal_handler_disconnect(C.gpointer(v.ptr), C.gulong(handle))
	C.g_closure_invalidate(signals[handle])
	delete(closures.m, signals[handle])
	delete(signals, handle)
}

/*
 * GInitiallyUnowned
 */

// InitiallyUnowned is a representation of GLib's GInitiallyUnowned.
type InitiallyUnowned struct {
	Object
}

/*
 * GValue
 */

// Value is a representation of GLib's GValue.
//
// Don't allocate Values on the stack or heap manually as they may not
// be properly unset when going out of scope. Instead, use ValueAlloc(),
// which will set the runtime finalizer to unset the Value after it has
// left scope.
type Value struct {
	GValue C.GValue
}

// Native() returns a pointer to the underlying GValue.
func (v *Value) Native() *C.GValue {
	return &v.GValue
}

// ValueAlloc() allocates a Value and sets a runtime finalizer to call
// g_value_unset() on the underlying GValue after leaving scope.
// ValueAlloc() returns a non-nil error if the allocation failed.
func ValueAlloc() (*Value, error) {
	c := C._g_value_alloc()
	if c == nil {
		return nil, nilPtrErr
	}
	v := &Value{*c}
	runtime.SetFinalizer(v, (*Value).unset)
	return v, nil
}

// ValueInit() is a wrapper around g_value_init() and allocates and
// initializes a new Value with the Type t.  A runtime finalizer is set
// to call g_value_unset() on the underlying GValue after leaving scope.
// ValueInit() returns a non-nil error if the allocation failed.
func ValueInit(t Type) (*Value, error) {
	c := C._g_value_init(C.GType(t))
	if c == nil {
		return nil, nilPtrErr
	}
	v := &Value{*c}
	runtime.SetFinalizer(v, (*Value).unset)
	return v, nil
}

func (v *Value) unset() {
	C.g_value_unset(v.Native())
}

// GetType() is a wrappr around the G_VALUE_HOLDS_GTYPE() macro and
// the g_value_get_gtype() function.  GetType() returns TYPE_INVALID if v
// does not hold a Type, or otherwise returns the Type of v.
func (v *Value) GetType() (actual Type, fundamental Type) {
	/*
	c := C._g_value_holds_gtype(C.gpointer(unsafe.Pointer(v.Native())))
	if gobool(c) {
		c := C.g_value_get_gtype(v.Native())
		return Type(c)
	}
	return TYPE_INVALID
	*/
	c_actual := C._g_value_type(v.Native())
	c_fundamental := C._g_value_fundamental(c_actual)
	return Type(c_actual), Type(c_fundamental)
}

// GValue() converts a Go type to a comparable GValue.  GValue()
// returns a non-nil error if the conversion was unsuccessful.
func GValue(v interface{}) (gvalue *Value, err error) {
	if v == nil {
		val, err := ValueInit(TYPE_POINTER)
		if err != nil {
			return nil, err
		}
		val.SetPointer(uintptr(0)) // technically not portable
		return val, nil
	}

	switch v.(type) {
	case bool:
		val, err := ValueInit(TYPE_BOOLEAN)
		if err != nil {
			return nil, err
		}
		val.SetBool(v.(bool))
		return val, nil
	case int8:
		val, err := ValueInit(TYPE_CHAR)
		if err != nil {
			return nil, err
		}
		val.SetSChar(v.(int8))
		return val, nil
	case int64:
		val, err := ValueInit(TYPE_INT64)
		if err != nil {
			return nil, err
		}
		val.SetInt64(v.(int64))
		return val, nil
	case int:
		val, err := ValueInit(TYPE_INT)
		if err != nil {
			return nil, err
		}
		val.SetInt(v.(int))
		return val, nil
	case uint8:
		val, err := ValueInit(TYPE_UCHAR)
		if err != nil {
			return nil, err
		}
		val.SetUChar(v.(uint8))
		return val, nil
	case uint64:
		val, err := ValueInit(TYPE_UINT64)
		if err != nil {
			return nil, err
		}
		val.SetUInt64(v.(uint64))
		return val, nil
	case uint:
		val, err := ValueInit(TYPE_UINT)
		if err != nil {
			return nil, err
		}
		val.SetUInt(v.(uint))
		return val, nil
	case float32:
		val, err := ValueInit(TYPE_FLOAT)
		if err != nil {
			return nil, err
		}
		val.SetFloat(v.(float32))
		return val, nil
	case float64:
		val, err := ValueInit(TYPE_DOUBLE)
		if err != nil {
			return nil, err
		}
		val.SetDouble(v.(float64))
		return val, nil
	case string:
		val, err := ValueInit(TYPE_STRING)
		if err != nil {
			return nil, err
		}
		val.SetString(v.(string))
		return val, nil
	default:
		if obj, ok := v.(*Object); ok {
			val, err := ValueInit(TYPE_OBJECT)
			if err != nil {
				return nil, err
			}
			val.SetInstance(uintptr(obj.ptr))
			return val, nil
		}

		/* Try this since above doesn't catch constants under other types */
		rval := reflect.ValueOf(v)
		switch rval.Kind() {
		case reflect.Int8:
			val, err := ValueInit(TYPE_CHAR)
			if err != nil {
				return nil, err
			}
			val.SetSChar(int8(rval.Int()))
			return val, nil
		case reflect.Int16:
			return nil, errors.New("Type not implemented")
		case reflect.Int32:
			return nil, errors.New("Type not implemented")
		case reflect.Int64:
			val, err := ValueInit(TYPE_INT64)
			if err != nil {
				return nil, err
			}
			val.SetInt64(rval.Int())
			return val, nil
		case reflect.Int:
			val, err := ValueInit(TYPE_INT)
			if err != nil {
				return nil, err
			}
			val.SetInt(int(rval.Int()))
			return val, nil
		case reflect.Uintptr:
			val, err := ValueInit(TYPE_POINTER)
			if err != nil {
				return nil, err
			}
			val.SetPointer(rval.Pointer())
			return val, nil
		}
	}
	return nil, errors.New("Type not implemented")
}

// GoValue() converts a Value to comparable Go type.  GoValue()
// returns a non-nil error if the conversion was unsuccessful.  The
// returned interface{} must be type asserted as the actual Go
// representation of the Value.
//
// This function is a wrapper around the many g_value_get_*()
// functions, depending on the type of the Value.
func (v *Value) GoValue() (interface{}, error) {
	actual, fundamental := v.GetType()
	switch fundamental {
	case TYPE_INVALID:
		return nil, errors.New("Invalid type")
	case TYPE_NONE:
		return nil, nil
	case TYPE_BOOLEAN:
		c := C.g_value_get_boolean(v.Native())
		return gobool(c), nil
	case TYPE_CHAR:
		c := C.g_value_get_schar(v.Native())
		return int8(c), nil
	case TYPE_UCHAR:
		c := C.g_value_get_uchar(v.Native())
		return uint8(c), nil
	case TYPE_INT64:
		c := C.g_value_get_int64(v.Native())
		return int64(c), nil
	case TYPE_INT:
		c := C.g_value_get_int(v.Native())
		return int(c), nil
	case TYPE_UINT64:
		c := C.g_value_get_uint64(v.Native())
		return uint64(c), nil
	case TYPE_UINT:
		c := C.g_value_get_uint(v.Native())
		return uint(c), nil
	case TYPE_FLOAT:
		c := C.g_value_get_float(v.Native())
		return float32(c), nil
	case TYPE_DOUBLE:
		c := C.g_value_get_double(v.Native())
		return float64(c), nil
	case TYPE_STRING:
		c := C.g_value_get_string(v.Native())
		return C.GoString((*C.char)(c)), nil
	case TYPE_OBJECT:
		c := C.g_value_get_object(v.Native())
		// TODO: need to try and return an actual pointer to the correct object type
		// this may require an additional cast()-like method for each module
		return unsafe.Pointer(c), nil
	default:
		return nil, errors.New("Type conversion not supported for type: " + actual.Name())
	}
}

// SetBool() is a wrapper around g_value_set_boolean().
func (v *Value) SetBool(val bool) {
	C.g_value_set_boolean(v.Native(), gbool(val))
}

// SetSChar() is a wrapper around g_value_set_schar().
func (v *Value) SetSChar(val int8) {
	C.g_value_set_schar(v.Native(), C.gint8(val))
}

// SetInt64() is a wrapper around g_value_set_int64().
func (v *Value) SetInt64(val int64) {
	C.g_value_set_int64(v.Native(), C.gint64(val))
}

// SetInt() is a wrapper around g_value_set_int().
func (v *Value) SetInt(val int) {
	C.g_value_set_int(v.Native(), C.gint(val))
}

// SetUChar() is a wrapper around g_value_set_uchar().
func (v *Value) SetUChar(val uint8) {
	C.g_value_set_uchar(v.Native(), C.guchar(val))
}

// SetUInt64() is a wrapper around g_value_set_uint64().
func (v *Value) SetUInt64(val uint64) {
	C.g_value_set_uint64(v.Native(), C.guint64(val))
}

// SetUInt() is a wrapper around g_value_set_uint().
func (v *Value) SetUInt(val uint) {
	C.g_value_set_uint(v.Native(), C.guint(val))
}

// SetFloat() is a wrapper around g_value_set_float().
func (v *Value) SetFloat(val float32) {
	C.g_value_set_float(v.Native(), C.gfloat(val))
}

// SetDouble() is a wrapper around g_value_set_double().
func (v *Value) SetDouble(val float64) {
	C.g_value_set_double(v.Native(), C.gdouble(val))
}

// SetString() is a wrapper around g_value_set_string().
func (v *Value) SetString(val string) {
	cstr := C.CString(val)
	defer C.free(unsafe.Pointer(cstr))
	C.g_value_set_string(v.Native(), (*C.gchar)(cstr))
}

// SetInstance() is a wrapper around g_value_set_instance().
func (v *Value) SetInstance(instance uintptr) {
	C.g_value_set_instance(v.Native(), C.gpointer(instance))
}

// SetPointer() is a wrapper around g_value_set_pointer().
func (v *Value) SetPointer(p uintptr) {
	C.g_value_set_pointer(v.Native(), C.gpointer(p))
}

// GetString() is a wrapper around g_value_get_string().  GetString()
// returns a non-nil error if g_value_get_string() returned a NULL
// pointer to distinguish between returning a NULL pointer and returning
// an empty string.
func (v *Value) GetString() (string, error) {
	c := C.g_value_get_string(v.Native())
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

func (v *Value) PeekPointer() interface{} {
	return C.g_value_peek_pointer(v.Native())
}

// valueSlice() converts a C array of GValues to a Go slice.
func valueSlice(n_values int, values *C.GValue) (slice []*C.GValue) {
	header := (*reflect.SliceHeader)((unsafe.Pointer(&slice)))
	header.Cap = n_values
	header.Len = n_values
	header.Data = uintptr(unsafe.Pointer(&values))
	return
}

