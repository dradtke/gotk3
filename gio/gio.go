package gio

// #cgo pkg-config: gio-2.0
// #include <gio/gio.h>
// #include "gio.go.h"
import "C"
import (
	"errors"
	"github.com/dradtke/gotk3/glib"
	"runtime"
	"unsafe"
)

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

var nilPtrErr = errors.New("cgo returned unexpected nil pointer")

/*
 * Constants
 */

type ApplicationFlags int

const (
	FLAGS_NONE           ApplicationFlags = C.G_APPLICATION_FLAGS_NONE
	IS_SERVICE                            = C.G_APPLICATION_IS_SERVICE
	IS_LAUNCHER                           = C.G_APPLICATION_IS_LAUNCHER
	HANDLES_OPEN                          = C.G_APPLICATION_HANDLES_OPEN
	HANDLES_COMMAND_LINE                  = C.G_APPLICATION_HANDLES_COMMAND_LINE
	SEND_ENVIRONMENT                      = C.G_APPLICATION_SEND_ENVIRONMENT
	NON_UNIQUE                            = C.G_APPLICATION_NON_UNIQUE
)

/*
 * Static methods
 */

// ApplicationIdIsValid() is a wrapper around g_application_id_is_valid().
func ApplicationIdIsValid(id string) bool {
	cstr := C.CString(id)
	defer C.free(unsafe.Pointer(cstr))
	return gobool(C.g_application_id_is_valid((*C.gchar)(cstr)))
}

/*
 * Application
 */

type Application struct {
	glib.Object
}

func wrapApplication(obj *glib.Object) Application {
	return Application{*obj}
}

func ApplicationNew(id string, flags ApplicationFlags) (*Application, error) {
	cstr := C.CString(id)
	defer C.free(unsafe.Pointer(cstr))
	c := C.g_application_new((*C.gchar)(cstr), C.GApplicationFlags(flags))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	a := wrapApplication(obj)
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &a, nil
}
