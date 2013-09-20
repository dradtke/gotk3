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

// Native() returns a pointer to the underlying GApplication.
func (v *Application) Native() *C.GApplication {
	if v == nil || v.Ptr() == nil {
		return nil
	}
	return (*C.GApplication)(v.Ptr())
}

// ApplicationNew() is a wrapper around g_application_new().
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

// ApplicationId() is a wrapper around g_application_get_application_id().
func (v *Application) ApplicationId() string {
	c := C.g_application_get_application_id(v.Native())
	return C.GoString((*C.char)(c))
}

// SetApplicationId() is a wrapper around g_application_set_application_id().
func (v *Application) SetApplicationId(id string) {
	cstr := C.CString(id)
	defer C.free(unsafe.Pointer(cstr))
	C.g_application_set_application_id(v.Native(), (*C.gchar)(cstr))
}

// InactivityTimeout() is a wrapper around g_application_get_inactivity_timeout().
func (v *Application) InactivityTimeout() uint {
	c := C.g_application_get_inactivity_timeout(v.Native())
	return uint(c)
}

// SetInactivityTimeout() is a wrapper around g_application_set_inactivity_timeout().
func (v *Application) SetInactivityTimeout(timeout uint) {
	C.g_application_set_inactivity_timeout(v.Native(), C.guint(timeout))
}

// Flags() is a wrapper around g_application_get_flags().
func (v *Application) Flags() ApplicationFlags {
	c := C.g_application_get_flags(v.Native())
	return ApplicationFlags(c)
}

// SetFlags() is a wrapper around g_application_set_flags().
func (v *Application) SetFlags(flags ApplicationFlags) {
	C.g_application_set_flags(v.Native(), C.GApplicationFlags(flags))
}

// Hold() is a wrapper around g_application_hold().
func (v *Application) Hold() {
	C.g_application_hold(v.Native())
}

// Release() is a wrapper around g_application_release().
func (v *Application) Release() {
	C.g_application_release(v.Native())
}

func (v *Application) Quit() {
	C.g_application_quit(v.Native())
}

// Run() is a wrapper around g_application_run().
func (v *Application) Run(args []string) int {
	var c C.int
	if args != nil {
		argc := len(args)
		argv := make([]*C.char, argc)
		for i, arg := range args {
			argv[i] = C.CString(arg)
		}
		c = C.g_application_run(v.Native(), C.int(argc),
			(**C.char)(unsafe.Pointer(&argv)))
	} else {
		c = C.g_application_run(v.Native(), 0, nil)
	}
	return int(c)
}

// Need at least GIO 2.38
/*
// MarkBusy() is a wrapper around g_application_mark_busy().
func (v *Application) MarkBusy() {
	C.g_application_mark_busy(v.Native())
}

// UnmarkBusy() is a wrapper around g_application_unmark_busy().
func (v *Application) UnmarkBusy() {
	C.g_application_unmark_busy(v.Native())
}
*/

// TODO: add GDBusConnection