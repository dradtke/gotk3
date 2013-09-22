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
Go bindings for GTK+ 3.  Supports version 3.8 and later.

Functions use the same names as the native C function calls, but use
CamelCase.  In cases where native GTK uses pointers to values to
simulate multiple return values, Go's native multiple return values
are used instead.  Whenever a native GTK call could return an
unexpected NULL pointer, an additonal error is returned in the Go
binding.

GTK's C API documentation can be very useful for understanding how the
functions in this package work and what each type is for.  This
documentation can be found at https://developer.gnome.org/gtk3/.

In addition to Go versions of the C GTK functions, every struct type
includes a function called Native(), taking itself as a receiver,
which returns the native C type or a pointer (in the case of
GObjects).  The returned C types are scoped to this gtk package and
must be converted to a local package before they can be used as
arguments to native GTK calls using cgo.

Memory management is handled in proper Go fashion, using runtime
finalizers to properly free memory when it is no longer needed.  Each
time a Go type is created with a pointer to a GObject, a reference is
added for Go, sinking the floating reference when necessary.  After
going out of scope and the next time Go's garbage collector is run, a
finalizer is run to remove Go's reference to the GObject.  When this
reference count hits zero (when neither Go nor GTK holds ownership)
the object will be freed internally by GTK.
*/
package gtk

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// #include "gtk.go.h"
import "C"
import (
	"errors"
	"fmt"
	"github.com/dradtke/gotk3/gdk"
	"github.com/dradtke/gotk3/glib"
	"os"
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

// Wrapper function for TestBoolConvs since cgo can't be used with
// testing package
func testBoolConvs() error {
	b := gobool(gbool(true))
	if b != true {
		return errors.New("Unexpected bool conversion result")
	}

	cb := gbool(gobool(C.gboolean(0)))
	if cb != C.gboolean(0) {
		return errors.New("Unexpected bool conversion result")
	}

	return nil
}

/*
 * Unexported vars
 */

var nilPtrErr = errors.New("cgo returned unexpected nil pointer")

/*
 * Constants
 */

// Align is a representation of GTK's GtkAlign.
type Align int

const (
	ALIGN_FILL   Align = C.GTK_ALIGN_FILL
	ALIGN_START        = C.GTK_ALIGN_START
	ALIGN_END          = C.GTK_ALIGN_END
	ALIGN_CENTER       = C.GTK_ALIGN_CENTER
)

// ButtonsType is a representation of GTK's GtkButtonsType.
type ButtonsType int

const (
	BUTTONS_NONE      ButtonsType = C.GTK_BUTTONS_NONE
	BUTTONS_OK                    = C.GTK_BUTTONS_OK
	BUTTONS_CLOSE                 = C.GTK_BUTTONS_CLOSE
	BUTTONS_CANCEL                = C.GTK_BUTTONS_CANCEL
	BUTTONS_YES_NO                = C.GTK_BUTTONS_YES_NO
	BUTTONS_OK_CANCEL             = C.GTK_BUTTONS_OK_CANCEL
)

// DialogFlags is a representation of GTK's GtkDialogFlags.
type DialogFlags int

const (
	DIALOG_MODAL               DialogFlags = C.GTK_DIALOG_MODAL
	DIALOG_DESTROY_WITH_PARENT             = C.GTK_DIALOG_DESTROY_WITH_PARENT
)

// EntryIconPosition is a representation of GTK's GtkEntryIconPosition.
type EntryIconPosition int

const (
	ENTRY_ICON_PRIMARY   EntryIconPosition = C.GTK_ENTRY_ICON_PRIMARY
	ENTRY_ICON_SECONDARY                   = C.GTK_ENTRY_ICON_SECONDARY
)

// IconSize is a representation of GTK's GtkIconSize.
type IconSize int

const (
	ICON_SIZE_INVALID       IconSize = C.GTK_ICON_SIZE_INVALID
	ICON_SIZE_MENU                   = C.GTK_ICON_SIZE_MENU
	ICON_SIZE_SMALL_TOOLBAR          = C.GTK_ICON_SIZE_SMALL_TOOLBAR
	ICON_SIZE_LARGE_TOOLBAR          = C.GTK_ICON_SIZE_LARGE_TOOLBAR
	ICON_SIZE_BUTTON                 = C.GTK_ICON_SIZE_BUTTON
	ICON_SIZE_DND                    = C.GTK_ICON_SIZE_DND
	ICON_SIZE_DIALOG                 = C.GTK_ICON_SIZE_DIALOG
)

// ImageType is a representation of GTK's GtkImageType.
type ImageType int

const (
	IMAGE_EMPTY     ImageType = C.GTK_IMAGE_EMPTY
	IMAGE_PIXBUF              = C.GTK_IMAGE_PIXBUF
	IMAGE_STOCK               = C.GTK_IMAGE_STOCK
	IMAGE_ICON_SET            = C.GTK_IMAGE_ICON_SET
	IMAGE_ANIMATION           = C.GTK_IMAGE_ANIMATION
	IMAGE_ICON_NAME           = C.GTK_IMAGE_ICON_NAME
	IMAGE_GICON               = C.GTK_IMAGE_GICON
)

// InputHints is a representation of GTK's GtkInputHints.
type InputHints int

const (
	INPUT_HINT_NONE                InputHints = C.GTK_INPUT_HINT_NONE
	INPUT_HINT_SPELLCHECK                     = C.GTK_INPUT_HINT_SPELLCHECK
	INPUT_HINT_NO_SPELLCHECK                  = C.GTK_INPUT_HINT_NO_SPELLCHECK
	INPUT_HINT_WORD_COMPLETION                = C.GTK_INPUT_HINT_WORD_COMPLETION
	INPUT_HINT_LOWERCASE                      = C.GTK_INPUT_HINT_LOWERCASE
	INPUT_HINT_UPPERCASE_CHARS                = C.GTK_INPUT_HINT_UPPERCASE_CHARS
	INPUT_HINT_UPPERCASE_WORDS                = C.GTK_INPUT_HINT_UPPERCASE_WORDS
	INPUT_HINT_UPPERCASE_SENTENCES            = C.GTK_INPUT_HINT_UPPERCASE_SENTENCES
	INPUT_HINT_INHIBIT_OSK                    = C.GTK_INPUT_HINT_INHIBIT_OSK
)

// InputPurpose is a representation of GTK's GtkInputPurpose.
type InputPurpose int

const (
	INPUT_PURPOSE_FREE_FORM InputPurpose = C.GTK_INPUT_PURPOSE_FREE_FORM
	INPUT_PURPOSE_ALPHA                  = C.GTK_INPUT_PURPOSE_ALPHA
	INPUT_PURPOSE_DIGITS                 = C.GTK_INPUT_PURPOSE_DIGITS
	INPUT_PURPOSE_NUMBER                 = C.GTK_INPUT_PURPOSE_NUMBER
	INPUT_PURPOSE_PHONE                  = C.GTK_INPUT_PURPOSE_PHONE
	INPUT_PURPOSE_URL                    = C.GTK_INPUT_PURPOSE_URL
	INPUT_PURPOSE_EMAIL                  = C.GTK_INPUT_PURPOSE_EMAIL
	INPUT_PURPOSE_NAME                   = C.GTK_INPUT_PURPOSE_NAME
	INPUT_PURPOSE_PASSWORD               = C.GTK_INPUT_PURPOSE_PASSWORD
	INPUT_PURPOSE_PIN                    = C.GTK_INPUT_PURPOSE_PIN
)

// MessageType is a representation of GTK's GtkMessageType.
type MessageType int

const (
	MESSAGE_INFO     MessageType = C.GTK_MESSAGE_INFO
	MESSAGE_WARNING              = C.GTK_MESSAGE_WARNING
	MESSAGE_QUESTION             = C.GTK_MESSAGE_QUESTION
	MESSAGE_ERROR                = C.GTK_MESSAGE_ERROR
	MESSAGE_OTHER                = C.GTK_MESSAGE_OTHER
)

// Orientation is a representation of GTK's GtkOrientation.
type Orientation int

const (
	ORIENTATION_HORIZONTAL Orientation = C.GTK_ORIENTATION_HORIZONTAL
	ORIENTATION_VERTICAL               = C.GTK_ORIENTATION_VERTICAL
)

// PackType is a representation of GTK's GtkPackType.
type PackType int

const (
	PACK_START PackType = C.GTK_PACK_START
	PACK_END            = C.GTK_PACK_END
)

// PolicyType is a representation of GTK's GtkPolicyType.
type PolicyType int

const (
	POLICY_ALWAYS    PolicyType = C.GTK_POLICY_ALWAYS
	POLICY_AUTOMATIC            = C.GTK_POLICY_AUTOMATIC
	POLICY_NEVER                = C.GTK_POLICY_NEVER
)

// PositionType is a representation of GTK's GtkPositionType.
type PositionType int

const (
	POS_LEFT   PositionType = C.GTK_POS_LEFT
	POS_RIGHT               = C.GTK_POS_RIGHT
	POS_TOP                 = C.GTK_POS_TOP
	POS_BOTTOM              = C.GTK_POS_BOTTOM
)

// ReliefStyle is a representation of GTK's GtkReliefStyle.
type ReliefStyle int

const (
	RELIEF_NORMAL ReliefStyle = C.GTK_RELIEF_NORMAL
	RELIEF_HALF               = C.GTK_RELIEF_HALF
	RELIEF_NONE               = C.GTK_RELIEF_NONE
)

// ResponseType is a representation of GTK's GtkResponseType.
type ResponseType int

const (
	RESPONSE_NONE         ResponseType = C.GTK_RESPONSE_NONE
	RESPONSE_REJECT                    = C.GTK_RESPONSE_REJECT
	RESPONSE_ACCEPT                    = C.GTK_RESPONSE_ACCEPT
	RESPONSE_DELETE_EVENT              = C.GTK_RESPONSE_DELETE_EVENT
	RESPONSE_OK                        = C.GTK_RESPONSE_OK
	RESPONSE_CANCEL                    = C.GTK_RESPONSE_CANCEL
	RESPONSE_CLOSE                     = C.GTK_RESPONSE_CLOSE
	RESPONSE_YES                       = C.GTK_RESPONSE_YES
	RESPONSE_NO                        = C.GTK_RESPONSE_NO
	RESPONSE_APPLY                     = C.GTK_RESPONSE_APPLY
	RESPONSE_HELP                      = C.GTK_RESPONSE_HELP
)

// Stock is a special type that does not have an equivalent type in
// GTK.  It is the type used as a parameter anytime an identifier for
// stock icons are needed.  A Stock must be type converted to string when
// function parameters may take a Stock, but when other string values are
// valid as well.
type Stock string

const (
	STOCK_ABOUT                         Stock = C.GTK_STOCK_ABOUT
	STOCK_ADD                                 = C.GTK_STOCK_ADD
	STOCK_APPLY                               = C.GTK_STOCK_APPLY
	STOCK_BOLD                                = C.GTK_STOCK_BOLD
	STOCK_CANCEL                              = C.GTK_STOCK_CANCEL
	STOCK_CAPS_LOCK_WARNING                   = C.GTK_STOCK_CAPS_LOCK_WARNING
	STOCK_CDROM                               = C.GTK_STOCK_CDROM
	STOCK_CLEAR                               = C.GTK_STOCK_CLEAR
	STOCK_CLOSE                               = C.GTK_STOCK_CLOSE
	STOCK_COLOR_PICKER                        = C.GTK_STOCK_COLOR_PICKER
	STOCK_CONNECT                             = C.GTK_STOCK_CONNECT
	STOCK_CONVERT                             = C.GTK_STOCK_CONVERT
	STOCK_COPY                                = C.GTK_STOCK_COPY
	STOCK_CUT                                 = C.GTK_STOCK_CUT
	STOCK_DELETE                              = C.GTK_STOCK_DELETE
	STOCK_DIALOG_AUTHENTICATION               = C.GTK_STOCK_DIALOG_AUTHENTICATION
	STOCK_DIALOG_INFO                         = C.GTK_STOCK_DIALOG_INFO
	STOCK_DIALOG_WARNING                      = C.GTK_STOCK_DIALOG_WARNING
	STOCK_DIALOG_ERROR                        = C.GTK_STOCK_DIALOG_ERROR
	STOCK_DIALOG_QUESTION                     = C.GTK_STOCK_DIALOG_QUESTION
	STOCK_DIRECTORY                           = C.GTK_STOCK_DIRECTORY
	STOCK_DISCARD                             = C.GTK_STOCK_DISCARD
	STOCK_DISCONNECT                          = C.GTK_STOCK_DISCONNECT
	STOCK_DND                                 = C.GTK_STOCK_DND
	STOCK_DND_MULTIPLE                        = C.GTK_STOCK_DND_MULTIPLE
	STOCK_EDIT                                = C.GTK_STOCK_EDIT
	STOCK_EXECUTE                             = C.GTK_STOCK_EXECUTE
	STOCK_FILE                                = C.GTK_STOCK_FILE
	STOCK_FIND                                = C.GTK_STOCK_FIND
	STOCK_FIND_AND_REPLACE                    = C.GTK_STOCK_FIND_AND_REPLACE
	STOCK_FLOPPY                              = C.GTK_STOCK_FLOPPY
	STOCK_FULLSCREEN                          = C.GTK_STOCK_FULLSCREEN
	STOCK_GOTO_BOTTOM                         = C.GTK_STOCK_GOTO_BOTTOM
	STOCK_GOTO_FIRST                          = C.GTK_STOCK_GOTO_FIRST
	STOCK_GOTO_LAST                           = C.GTK_STOCK_GOTO_LAST
	STOCK_GOTO_TOP                            = C.GTK_STOCK_GOTO_TOP
	STOCK_GO_BACK                             = C.GTK_STOCK_GO_BACK
	STOCK_GO_DOWN                             = C.GTK_STOCK_GO_DOWN
	STOCK_GO_FORWARD                          = C.GTK_STOCK_GO_FORWARD
	STOCK_GO_UP                               = C.GTK_STOCK_GO_UP
	STOCK_HARDDISK                            = C.GTK_STOCK_HARDDISK
	STOCK_HELP                                = C.GTK_STOCK_HELP
	STOCK_HOME                                = C.GTK_STOCK_HOME
	STOCK_INDEX                               = C.GTK_STOCK_INDEX
	STOCK_INDENT                              = C.GTK_STOCK_INDENT
	STOCK_INFO                                = C.GTK_STOCK_INFO
	STOCK_ITALIC                              = C.GTK_STOCK_ITALIC
	STOCK_JUMP_TO                             = C.GTK_STOCK_JUMP_TO
	STOCK_JUSTIFY_CENTER                      = C.GTK_STOCK_JUSTIFY_CENTER
	STOCK_JUSTIFY_FILL                        = C.GTK_STOCK_JUSTIFY_FILL
	STOCK_JUSTIFY_LEFT                        = C.GTK_STOCK_JUSTIFY_LEFT
	STOCK_JUSTIFY_RIGHT                       = C.GTK_STOCK_JUSTIFY_RIGHT
	STOCK_LEAVE_FULLSCREEN                    = C.GTK_STOCK_LEAVE_FULLSCREEN
	STOCK_MISSING_IMAGE                       = C.GTK_STOCK_MISSING_IMAGE
	STOCK_MEDIA_FORWARD                       = C.GTK_STOCK_MEDIA_FORWARD
	STOCK_MEDIA_NEXT                          = C.GTK_STOCK_MEDIA_NEXT
	STOCK_MEDIA_PAUSE                         = C.GTK_STOCK_MEDIA_PAUSE
	STOCK_MEDIA_PLAY                          = C.GTK_STOCK_MEDIA_PLAY
	STOCK_MEDIA_PREVIOUS                      = C.GTK_STOCK_MEDIA_PREVIOUS
	STOCK_MEDIA_RECORD                        = C.GTK_STOCK_MEDIA_RECORD
	STOCK_MEDIA_REWIND                        = C.GTK_STOCK_MEDIA_REWIND
	STOCK_MEDIA_STOP                          = C.GTK_STOCK_MEDIA_STOP
	STOCK_NETWORK                             = C.GTK_STOCK_NETWORK
	STOCK_NEW                                 = C.GTK_STOCK_NEW
	STOCK_NO                                  = C.GTK_STOCK_NO
	STOCK_OK                                  = C.GTK_STOCK_OK
	STOCK_OPEN                                = C.GTK_STOCK_OPEN
	STOCK_ORIENTATION_PORTRAIT                = C.GTK_STOCK_ORIENTATION_PORTRAIT
	STOCK_ORIENTATION_LANDSCAPE               = C.GTK_STOCK_ORIENTATION_LANDSCAPE
	STOCK_ORIENTATION_REVERSE_LANDSCAPE       = C.GTK_STOCK_ORIENTATION_REVERSE_LANDSCAPE
	STOCK_ORIENTATION_REVERSE_PORTRAIT        = C.GTK_STOCK_ORIENTATION_REVERSE_PORTRAIT
	STOCK_PAGE_SETUP                          = C.GTK_STOCK_PAGE_SETUP
	STOCK_PASTE                               = C.GTK_STOCK_PASTE
	STOCK_PREFERENCES                         = C.GTK_STOCK_PREFERENCES
	STOCK_PRINT                               = C.GTK_STOCK_PRINT
	STOCK_PRINT_ERROR                         = C.GTK_STOCK_PRINT_ERROR
	STOCK_PRINT_PAUSED                        = C.GTK_STOCK_PRINT_PAUSED
	STOCK_PRINT_PREVIEW                       = C.GTK_STOCK_PRINT_PREVIEW
	STOCK_PRINT_REPORT                        = C.GTK_STOCK_PRINT_REPORT
	STOCK_PRINT_WARNING                       = C.GTK_STOCK_PRINT_WARNING
	STOCK_PROPERTIES                          = C.GTK_STOCK_PROPERTIES
	STOCK_QUIT                                = C.GTK_STOCK_QUIT
	STOCK_REDO                                = C.GTK_STOCK_REDO
	STOCK_REFRESH                             = C.GTK_STOCK_REFRESH
	STOCK_REMOVE                              = C.GTK_STOCK_REMOVE
	STOCK_REVERT_TO_SAVED                     = C.GTK_STOCK_REVERT_TO_SAVED
	STOCK_SAVE                                = C.GTK_STOCK_SAVE
	STOCK_SAVE_AS                             = C.GTK_STOCK_SAVE_AS
	STOCK_SELECT_ALL                          = C.GTK_STOCK_SELECT_ALL
	STOCK_SELECT_COLOR                        = C.GTK_STOCK_SELECT_COLOR
	STOCK_SELECT_FONT                         = C.GTK_STOCK_SELECT_FONT
	STOCK_SORT_ASCENDING                      = C.GTK_STOCK_SORT_ASCENDING
	STOCK_SORT_DESCENDING                     = C.GTK_STOCK_SORT_DESCENDING
	STOCK_SPELL_CHECK                         = C.GTK_STOCK_SPELL_CHECK
	STOCK_STOP                                = C.GTK_STOCK_STOP
	STOCK_STRIKETHROUGH                       = C.GTK_STOCK_STRIKETHROUGH
	STOCK_UNDELETE                            = C.GTK_STOCK_UNDELETE
	STOCK_UNDERLINE                           = C.GTK_STOCK_UNDERLINE
	STOCK_UNDO                                = C.GTK_STOCK_UNDO
	STOCK_UNINDENT                            = C.GTK_STOCK_UNINDENT
	STOCK_YES                                 = C.GTK_STOCK_YES
	STOCK_ZOOM_100                            = C.GTK_STOCK_ZOOM_100
	STOCK_ZOOM_FIT                            = C.GTK_STOCK_ZOOM_FIT
	STOCK_ZOOM_IN                             = C.GTK_STOCK_ZOOM_IN
	STOCK_ZOOM_OUT                            = C.GTK_STOCK_ZOOM_OUT
)

// TreeModelFlags is a representation of GTK's GtkTreeModelFlags.
type TreeModelFlags int

const (
	TREE_MODEL_ITERS_PERSIST TreeModelFlags = C.GTK_TREE_MODEL_ITERS_PERSIST
	TREE_MODEL_LIST_ONLY                    = C.GTK_TREE_MODEL_LIST_ONLY
)

// WindowPosition is a representation of GTK's GtkWindowPosition.
type WindowPosition int

const (
	WIN_POS_NONE             WindowPosition = C.GTK_WIN_POS_NONE
	WIN_POS_CENTER                          = C.GTK_WIN_POS_CENTER
	WIN_POS_MOUSE                           = C.GTK_WIN_POS_MOUSE
	WIN_POS_CENTER_ALWAYS                   = C.GTK_WIN_POS_CENTER_ALWAYS
	WIN_POS_CENTER_ON_PARENT                = C.GTK_WIN_POS_CENTER_ON_PARENT
)

// WindowType is a representation of GTK's GtkWindowType.
type WindowType int

const (
	WINDOW_TOPLEVEL WindowType = C.GTK_WINDOW_TOPLEVEL
	WINDOW_POPUP               = C.GTK_WINDOW_POPUP
)

/*
 * Init and main event loop
 */

/*
Init() is a wrapper around gtk_init() and must be called before any
other GTK calls and is used to initialize everything necessary.

In addition to setting up GTK for usage, a pointer to a slice of
strings may be passed in to parse standard GTK command line arguments.
args will be modified to remove any flags that were handled.
Alternatively, nil may be passed in to not perform any command line
parsing.
*/
func Init(args *[]string) {
	if args != nil {
		argc := C.int(len(*args))
		argv := make([]*C.char, argc)
		for i, arg := range *args {
			argv[i] = C.CString(arg)
		}
		C.gtk_init((*C.int)(unsafe.Pointer(&argc)),
			(***C.char)(unsafe.Pointer(&argv)))
		unhandled := make([]string, argc)
		for i := 0; i < int(argc); i++ {
			unhandled[i] = C.GoString(argv[i])
			C.free(unsafe.Pointer(argv[i]))
		}
		*args = unhandled
	} else {
		C.gtk_init(nil, nil)
	}
}

// Main() is a wrapper around gtk_main() and runs the GTK main loop,
// blocking until MainQuit() is called.
func Main() {
	C.gtk_main()
}

// MainQuit() is a wrapper around gtk_main_quit() is used to terminate
// the GTK main loop (started by Main()).
func MainQuit() {
	C.gtk_main_quit()
}

func MainLevel() uint {
	return uint(C.gtk_main_level())
}

func MainIteration() bool {
	return gobool(C.gtk_main_iteration())
}

/*
 * GtkAdjustment
 */

// Adjustment is a representation of GTK's GtkAdjustment.
type Adjustment struct {
	glib.InitiallyUnowned
}

var adjustmentType = glib.Type(C.gtk_adjustment_get_type())

func GetAdjustmentType() glib.Type {
	return adjustmentType
}

// Native() returns a pointer to the underlying GtkAdjustment.
func (v *Adjustment) Native() *C.GtkAdjustment {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(adjustmentType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkAdjustment)(v.Ptr())
}

func wrapAdjustment(obj *glib.Object) Adjustment {
	return Adjustment{glib.InitiallyUnowned{*obj}}
}

/*
 * GtkBin
 */

// Bin is a representation of GTK's GtkBin.
type Bin struct {
	Container
}

var binType = glib.Type(C.gtk_bin_get_type())

func GetBinType() glib.Type {
	return binType
}

// Native() returns a pointer to the underlying GtkBin.
func (v *Bin) Native() *C.GtkBin {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(binType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkBin)(v.Ptr())
}

func wrapBin(obj *glib.Object) (b Bin) {
	b.Container = wrapContainer(obj)
	return
}

// Child() is a wrapper around gtk_bin_get_child().
func (v *Bin) Child() (*Widget, error) {
	c := C.gtk_bin_get_child(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

/*
 * GtkBuildable
 */

// Buildable is a representation of GTK's GtkBuildable GInterface.
type Buildable struct {
	ptr unsafe.Pointer
}

var buildableType = glib.Type(C.gtk_buildable_get_type())

func GetBuildableType() glib.Type {
	return buildableType
}

func (v *Buildable) Native() *C.GtkBuildable {
	if v == nil {
		fmt.Println("nil object, not getting native buildable")
		return nil
	}
	return (*C.GtkBuildable)(v.ptr)
}

func (v *Buildable) BuildableName() string {
	return C.GoString((*C.char)(C.gtk_buildable_get_name(v.Native())))
}

/*
 * GtkBuilder
 */

// Builder is a representation of GTK's GtkBuilder.
type Builder struct {
	*glib.Object
}

var builderType = glib.Type(C.gtk_builder_get_type())

func GetBuilderType() glib.Type {
	return glib.Type(builderType)
}

// Native() returns a pointer to the underlying GtkBuilder.
func (v *Builder) Native() *C.GtkBuilder {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(builderType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkBuilder)(v.Ptr())
}

// BuilderNew() is a wrapper around gtk_builder_new().
func BuilderNew() (*Builder, error) {
	c := C.gtk_builder_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	b := &Builder{obj}
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return b, nil
}

// AddFromFile() is a wrapper around gtk_builder_add_from_file().
func (b *Builder) AddFromFile(filename string) error {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	var err *C.GError = nil
	res := C.gtk_builder_add_from_file(b.Native(), (*C.gchar)(cstr), &err)
	if res == 0 {
		defer C.g_error_free(err)
		return errors.New(C.GoString((*C.char)(C.error_get_message(err))))
	}
	return nil
}

// AddFromResource() is a wrapper around gtk_builder_add_from_resource().
func (b *Builder) AddFromResource(path string) error {
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))
	var err *C.GError = nil
	res := C.gtk_builder_add_from_resource(b.Native(), (*C.gchar)(cstr), &err)
	if res == 0 {
		defer C.g_error_free(err)
		return errors.New(C.GoString((*C.char)(C.error_get_message(err))))
	}
	return nil
}

// AddFromString() is a wrapper around gtk_builder_add_from_string().
func (b *Builder) AddFromString(str string) error {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	length := (C.gsize)(len(str))
	var err *C.GError = nil
	res := C.gtk_builder_add_from_string(b.Native(), (*C.gchar)(cstr), length, &err)
	if res == 0 {
		defer C.g_error_free(err)
		return errors.New(C.GoString((*C.char)(C.error_get_message(err))))
	}
	return nil
}

// Object() is a wrapper around gtk_builder_get_object(). The returned result
// is an IObject, so it will need to be type-asserted to the appropriate type before
// being used. For example, to get an object and type assert it as a window:
//
//   obj, err := builder.GetObject("window")
//   if err != nil {
//       // object not found
//       return
//   }
//   if w, ok := obj.(*gtk.Window); ok {
//       // do stuff with w here
//   } else {
//       // not a *gtk.Window
//   }
//
func (b *Builder) GetObject(name string) (glib.IObject, error) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_builder_get_object(b.Native(), (*C.gchar)(cstr))
	if c == nil {
		return nil, errors.New("object '" + name + "' not found")
	}
	obj, err := cast(c)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

/*
 * GtkButton
 */

// Button is a representation of GTK's GtkButton.
type Button struct {
	Bin
}

var buttonType = glib.Type(C.gtk_button_get_type())

func GetButtonType() glib.Type {
	return glib.Type(buttonType)
}

// Native() returns a pointer to the underlying GtkButton.
func (v *Button) Native() *C.GtkButton {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(buttonType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkButton)(v.Ptr())
}

func wrapButton(obj *glib.Object) (b Button) {
	b.Bin = wrapBin(obj)
	return b
}

// ButtonNew() is a wrapper around gtk_button_new().
func ButtonNew() (*Button, error) {
	c := C.gtk_button_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	b := wrapButton(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &b, nil
}

// ButtonNewWithLabel() is a wrapper around gtk_button_new_with_label().
func ButtonNewWithLabel(label string) (*Button, error) {
	cstr := C.CString(label)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_button_new_with_label((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	b := wrapButton(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &b, nil
}

// ButtonNewFromStock() is a wrapper around gtk_button_new_from_stock().
func ButtonNewFromStock(stock Stock) (*Button, error) {
	cstr := C.CString(string(stock))
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_button_new_from_stock((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	b := wrapButton(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &b, nil
}

// ButtonNewWithMnemonic() is a wrapper around gtk_button_new_with_mnemonic().
func ButtonNewWithMnemonic(label string) (*Button, error) {
	cstr := C.CString(label)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_button_new_with_mnemonic((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	b := wrapButton(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &b, nil
}

// Clicked() is a wrapper around gtk_button_clicked().
func (v *Button) Clicked() {
	C.gtk_button_clicked(v.Native())
}

// SetRelief() is a wrapper around gtk_button_set_relief().
func (v *Button) SetRelief(newStyle ReliefStyle) {
	C.gtk_button_set_relief(v.Native(), C.GtkReliefStyle(newStyle))
}

// Relief() is a wrapper around gtk_button_get_relief().
func (v *Button) Relief() ReliefStyle {
	c := C.gtk_button_get_relief(v.Native())
	return ReliefStyle(c)
}

// SetLabel() is a wrapper around gtk_button_set_label().
func (v *Button) SetLabel(label string) {
	cstr := C.CString(label)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_button_set_label(v.Native(), (*C.gchar)(cstr))
}

// Label() is a wrapper around gtk_button_get_label().
func (v *Button) Label() (string, error) {
	c := C.gtk_button_get_label(v.Native())
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// SetUseUnderline() is a wrapper around gtk_button_set_use_underline().
func (v *Button) SetUseUnderline(useUnderline bool) {
	C.gtk_button_set_use_underline(v.Native(), gbool(useUnderline))
}

// UseUnderline() is a wrapper around gtk_button_get_use_underline().
func (v *Button) UseUnderline() bool {
	c := C.gtk_button_get_use_underline(v.Native())
	return gobool(c)
}

// SetUseStock() is a wrapper around gtk_button_set_use_stock().
func (v *Button) SetUseStock(useStock bool) {
	C.gtk_button_set_use_stock(v.Native(), gbool(useStock))
}

// UseStock() is a wrapper around gtk_button_get_use_stock().
func (v *Button) UseStock() bool {
	c := C.gtk_button_get_use_stock(v.Native())
	return gobool(c)
}

// SetFocusOnClick() is a wrapper around gtk_button_set_focus_on_click().
func (v *Button) SetFocusOnClick(focusOnClick bool) {
	C.gtk_button_set_focus_on_click(v.Native(), gbool(focusOnClick))
}

// FocusOnClick() is a wrapper around gtk_button_get_focus_on_click().
func (v *Button) FocusOnClick() bool {
	c := C.gtk_button_get_focus_on_click(v.Native())
	return gobool(c)
}

// SetAlignment() is a wrapper around gtk_button_set_alignment().
func (v *Button) SetAlignment(xalign, yalign float32) {
	C.gtk_button_set_alignment(v.Native(), (C.gfloat)(xalign),
		(C.gfloat)(yalign))
}

// Alignment() is a wrapper around gtk_button_get_alignment().
func (v *Button) Alignment() (xalign, yalign float32) {
	var x, y C.gfloat
	C.gtk_button_get_alignment(v.Native(), &x, &y)
	return float32(x), float32(y)
}

// SetImage() is a wrapper around gtk_button_set_image().
func (v *Button) SetImage(image IWidget) {
	C.gtk_button_set_image(v.Native(), image.toWidget())
}

// Image() is a wrapper around gtk_button_get_image().
func (v *Button) Image() (*Widget, error) {
	c := C.gtk_button_get_image(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// SetImagePosition() is a wrapper around gtk_button_set_image_position().
func (v *Button) SetImagePosition(position PositionType) {
	C.gtk_button_set_image_position(v.Native(), C.GtkPositionType(position))
}

// ImagePosition() is a wrapper around gtk_button_get_image_position().
func (v *Button) ImagePosition() PositionType {
	c := C.gtk_button_get_image_position(v.Native())
	return PositionType(c)
}

// SetAlwaysShowImage() is a wrapper around gtk_button_set_always_show_image().
func (v *Button) SetAlwaysShowImage(alwaysShow bool) {
	C.gtk_button_set_always_show_image(v.Native(), gbool(alwaysShow))
}

// AlwaysShowImage() is a wrapper around gtk_button_get_always_show_image().
func (v *Button) AlwaysShowImage() bool {
	c := C.gtk_button_get_always_show_image(v.Native())
	return gobool(c)
}

// EventWindow() is a wrapper around gtk_button_get_event_window().
func (v *Button) EventWindow() (*gdk.Window, error) {
	c := C.gtk_button_get_event_window(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := &gdk.Window{obj}
	w.Ref()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return w, nil
}

/*
 * GtkBox
 */

// Box is a representation of GTK's GtkBox.
type Box struct {
	Container
}

var boxType = glib.Type(C.gtk_box_get_type())

func GetBoxType() glib.Type {
	return glib.Type(boxType)
}

// Native() returns a pointer to the underlying GtkBox.
func (v *Box) Native() *C.GtkBox {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(boxType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkBox)(v.Ptr())
}

func wrapBox(obj *glib.Object) (b Box) {
	b.Container = wrapContainer(obj)
	return
}

// BoxNew() is a wrapper around gtk_box_new().
func BoxNew(orientation Orientation, spacing int) (*Box, error) {
	c := C.gtk_box_new(C.GtkOrientation(orientation), C.gint(spacing))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	b := wrapBox(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &b, nil
}

// PackStart() is a wrapper around gtk_box_pack_start().
func (v *Box) PackStart(child IWidget, expand, fill bool, padding uint) {
	C.gtk_box_pack_start(v.Native(), child.toWidget(), gbool(expand),
		gbool(fill), C.guint(padding))
}

// PackEnd() is a wrapper around gtk_box_pack_end().
func (v *Box) PackEnd(child IWidget, expand, fill bool, padding uint) {
	C.gtk_box_pack_end(v.Native(), child.toWidget(), gbool(expand),
		gbool(fill), C.guint(padding))
}

// Homogeneous() is a wrapper around gtk_box_get_homogeneous().
func (v *Box) Homogeneous() bool {
	c := C.gtk_box_get_homogeneous(v.Native())
	return gobool(c)
}

// SetHomogeneous() is a wrapper around gtk_box_set_homogeneous().
func (v *Box) SetHomogeneous(homogeneous bool) {
	C.gtk_box_set_homogeneous(v.Native(), gbool(homogeneous))
}

// Spacing() is a wrapper around gtk_box_get_spacing().
func (v *Box) Spacing() int {
	c := C.gtk_box_get_spacing(v.Native())
	return int(c)
}

// SetSpacing() is a wrapper around gtk_box_set_spacing()
func (v *Box) SetSpacing(spacing int) {
	C.gtk_box_set_spacing(v.Native(), C.gint(spacing))
}

// ReorderChild() is a wrapper around gtk_box_reorder_child().
func (v *Box) ReorderChild(child IWidget, position int) {
	C.gtk_box_reorder_child(v.Native(), child.toWidget(), C.gint(position))
}

// QueryChildPacking() is a wrapper around gtk_box_query_child_packing().
func (v *Box) QueryChildPacking(child IWidget) (expand, fill bool, padding uint, packType PackType) {
	var cexpand, cfill C.gboolean
	var cpadding C.guint
	var cpackType C.GtkPackType

	C.gtk_box_query_child_packing(v.Native(), child.toWidget(), &cexpand,
		&cfill, &cpadding, &cpackType)
	return gobool(cexpand), gobool(cfill), uint(cpadding), PackType(cpackType)
}

// SetChildPacking() is a wrapper around gtk_box_set_child_packing().
func (v *Box) SetChildPacking(child IWidget, expand, fill bool, padding uint, packType PackType) {
	C.gtk_box_set_child_packing(v.Native(), child.toWidget(), gbool(expand),
		gbool(fill), C.guint(padding), C.GtkPackType(packType))
}

/*
 * GtkCellLayout
 */

// CellLayout is a representation of GTK's GtkCellLayout GInterface.
type CellLayout struct {
	*glib.Object
}

var cellLayoutType = glib.Type(C.gtk_cell_layout_get_type())

func GetCellLayoutType() glib.Type {
	return glib.Type(cellLayoutType)
}

// ICellLayout is an interface type implemented by all structs
// embedding a CellLayout.  It is meant to be used as an argument type
// for wrapper functions that wrap around a C GTK function taking a
// GtkCellLayout.
type ICellLayout interface {
	toCellLayout() *C.GtkCellLayout
}

// Native() returns a pointer to the underlying GObject as a GtkCellLayout.
func (v *CellLayout) Native() *C.GtkCellLayout {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(cellLayoutType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkCellLayout)(v.Ptr())
}

func wrapCellLayout(obj *glib.Object) CellLayout {
	return CellLayout{obj}
}

func (v *CellLayout) toCellLayout() *C.GtkCellLayout {
	if v == nil {
		return nil
	}
	return v.Native()
}

// PackStart() is a wrapper around gtk_cell_layout_pack_start().
func (v *CellLayout) PackStart(cell ICellRenderer, expand bool) {
	C.gtk_cell_layout_pack_start(v.Native(), cell.toCellRenderer(),
		gbool(expand))
}

// AddAttribute() is a wrapper around gtk_cell_layout_add_attribute().
func (v *CellLayout) AddAttribute(cell ICellRenderer, attribute string, column int) {
	cstr := C.CString(attribute)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_cell_layout_add_attribute(v.Native(), cell.toCellRenderer(),
		(*C.gchar)(cstr), C.gint(column))
}

/*
 * GtkCellRenderer
 */

// CellRenderer is a representation of GTK's GtkCellRenderer.
type CellRenderer struct {
	glib.InitiallyUnowned
}

var cellRendererType = glib.Type(C.gtk_cell_renderer_get_type())

func GetCellRendererType() glib.Type {
	return glib.Type(cellRendererType)
}

// ICellRenderer is an interface type implemented by all structs
// embedding a CellRenderer.  It is meant to be used as an argument type
// for wrapper functions that wrap around a C GTK function taking a
// GtkCellRenderer.
type ICellRenderer interface {
	toCellRenderer() *C.GtkCellRenderer
}

// Native() returns a pointer to the underlying GtkCellRenderer.
func (v *CellRenderer) Native() *C.GtkCellRenderer {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(cellRendererType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkCellRenderer)(v.Ptr())
}

func (v *CellRenderer) toCellRenderer() *C.GtkCellRenderer {
	if v == nil {
		return nil
	}
	return v.Native()
}

func wrapCellRenderer(obj *glib.Object) CellRenderer {
	return CellRenderer{glib.InitiallyUnowned{*obj}}
}

/*
 * GtkCellRendererText
 */

// CellRendererText is a representation of GTK's GtkCellRendererText.
type CellRendererText struct {
	CellRenderer
}

var cellRendererTextType = glib.Type(C.gtk_cell_renderer_text_get_type())

func GetCellRendererTextType() glib.Type {
	return glib.Type(cellRendererTextType)
}

// Native() returns a pointer to the underlying GtkCellRendererText.
func (v *CellRendererText) Native() *C.GtkCellRendererText {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(cellRendererTextType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkCellRendererText)(v.Ptr())
}

func (v *CellRendererText) toCellRenderer() *C.GtkCellRenderer {
	if v == nil {
		return nil
	}
	return v.CellRenderer.Native()
}

func wrapCellRendererText(obj *glib.Object) CellRendererText {
	return CellRendererText{CellRenderer{glib.InitiallyUnowned{*obj}}}
}

// CellRendererTextNew() is a wrapper around gtk_cell_renderer_text_new().
func CellRendererTextNew() (*CellRendererText, error) {
	c := C.gtk_cell_renderer_text_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	crt := wrapCellRendererText(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &crt, nil
}

/*
 * GtkClipboard
 */

// Clipboard is a wrapper around GTK's GtkClipboard.
type Clipboard struct {
	*glib.Object
}

var clipboardType = glib.Type(C.gtk_clipboard_get_type())

func GetClipboardType() glib.Type {
	return glib.Type(clipboardType)
}

// Native() returns a pointer to the underlying GtkClipboard.
func (v *Clipboard) Native() *C.GtkClipboard {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(clipboardType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkClipboard)(v.Ptr())
}

func wrapClipboard(obj *glib.Object) Clipboard {
	return Clipboard{obj}
}

// Clipboard() is a wrapper around gtk_clipboard_get().
func ClipboardGet(atom gdk.Atom) (*Clipboard, error) {
	c := C.gtk_clipboard_get(atom.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	cb := &Clipboard{obj}
	obj.Ref()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return cb, nil
}

// ClipboardForDisplay() is a wrapper around gtk_clipboard_get_for_display().
func ClipboardGetForDisplay(display *gdk.Display, atom gdk.Atom) (*Clipboard, error) {
	c := C.gtk_clipboard_get_for_display(display.Native(), atom.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	cb := &Clipboard{obj}
	obj.Ref()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return cb, nil
}

// SetText() is a wrapper around gtk_clipboard_set_text().
func (v *Clipboard) SetText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_clipboard_set_text(v.Native(), (*C.gchar)(cstr),
		C.gint(len(text)))
}

/*
 * GtkComboBox
 */

// ComboBox is a representation of GTK's GtkComboBox.
type ComboBox struct {
	Bin

	// Interfaces
	CellLayout
}

var comboBoxType = glib.Type(C.gtk_combo_box_get_type())

func GetComboBoxType() glib.Type {
	return glib.Type(comboBoxType)
}

// Native() returns a pointer to the underlying GtkComboBox.
func (v *ComboBox) Native() *C.GtkComboBox {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(comboBoxType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkComboBox)(v.Ptr())
}

func (v *ComboBox) toCellLayout() *C.GtkCellLayout {
	if v == nil {
		return nil
	}
	return (*C.GtkCellLayout)(v.Ptr())
}

func wrapComboBox(obj *glib.Object) (c ComboBox) {
	c.Bin = wrapBin(obj)
	c.CellLayout = CellLayout{obj}
	return
}

// ComboBoxNew() is a wrapper around gtk_combo_box_new().
func ComboBoxNew() (*ComboBox, error) {
	c := C.gtk_combo_box_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	cb := wrapComboBox(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &cb, nil
}

// ComboBoxNewWithEntry() is a wrapper around gtk_combo_box_new_with_entry().
func ComboBoxNewWithEntry() (*ComboBox, error) {
	c := C.gtk_combo_box_new_with_entry()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	cb := wrapComboBox(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &cb, nil
}

// ComboBoxNewWithModel() is a wrapper around gtk_combo_box_new_with_model().
func ComboBoxNewWithModel(model ITreeModel) (*ComboBox, error) {
	c := C.gtk_combo_box_new_with_model(model.toTreeModel())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	cb := wrapComboBox(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &cb, nil
}

// Active() is a wrapper around gtk_combo_box_get_active().
func (v *ComboBox) Active() int {
	c := C.gtk_combo_box_get_active(v.Native())
	return int(c)
}

// SetActive() is a wrapper around gtk_combo_box_set_active().
func (v *ComboBox) SetActive(index int) {
	C.gtk_combo_box_set_active(v.Native(), C.gint(index))
}

/*
 * GtkContainer
 */

// Container is a representation of GTK's GtkContainer.
type Container struct {
	Widget
}

var containerType = glib.Type(C.gtk_container_get_type())

func GetContainerType() glib.Type {
	return glib.Type(containerType)
}

// Native() returns a pointer to the underlying GtkContainer.
func (v *Container) Native() *C.GtkContainer {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(containerType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkContainer)(v.Ptr())
}

func wrapContainer(obj *glib.Object) (c Container) {
	c.Widget = wrapWidget(obj)
	return
}

// Add() is a wrapper around gtk_container_add().
func (v *Container) Add(w IWidget) {
	C.gtk_container_add(v.Native(), w.toWidget())
}

// Remove() is a wrapper around gtk_container_remove().
func (v *Container) Remove(w IWidget) {
	C.gtk_container_remove(v.Native(), w.toWidget())
}

/*
 * GtkDialog
 */

// Dialog is a representation of GTK's GtkDialog.
type Dialog struct {
	Window
}

var dialogType = glib.Type(C.gtk_dialog_get_type())

func GetDialogType() glib.Type {
	return glib.Type(dialogType)
}

// Native() returns a pointer to the underlying GtkDialog.
func (v *Dialog) Native() *C.GtkDialog {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(dialogType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkDialog)(v.Ptr())
}

func wrapDialog(obj *glib.Object) (d Dialog) {
	d.Window = wrapWindow(obj)
	return
}

// DialogNew() is a wrapper around gtk_dialog_new().
func DialogNew() (*Dialog, error) {
	c := C.gtk_dialog_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	d := wrapDialog(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &d, nil
}

// Run() is a wrapper around gtk_dialog_run().
func (v *Dialog) Run() int {
	c := C.gtk_dialog_run(v.Native())
	return int(c)
}

// Response() is a wrapper around gtk_dialog_response().
func (v *Dialog) Response(response ResponseType) {
	C.gtk_dialog_response(v.Native(), C.gint(response))
}

// AddButton() is a wrapper around gtk_dialog_add_button().  text may
// be either the literal button text, or a Stock type converted to a
// string.
func (v *Dialog) AddButton(text string, id ResponseType) (*Button, error) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_dialog_add_button(v.Native(), (*C.gchar)(cstr), C.gint(id))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	b := wrapButton(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &b, nil
}

// AddActionWidget() is a wrapper around gtk_dialog_add_action_widget().
func (v *Dialog) AddActionWidget(child IWidget, id ResponseType) {
	C.gtk_dialog_add_action_widget(v.Native(), child.toWidget(), C.gint(id))
}

// SetDefaultResponse() is a wrapper around gtk_dialog_set_default_response().
func (v *Dialog) SetDefaultResponse(id ResponseType) {
	C.gtk_dialog_set_default_response(v.Native(), C.gint(id))
}

// SetResponseSensitive() is a wrapper around
// gtk_dialog_set_response_sensitive().
func (v *Dialog) SetResponseSensitive(id ResponseType, setting bool) {
	C.gtk_dialog_set_response_sensitive(v.Native(), C.gint(id),
		gbool(setting))
}

// ResponseForWidget() is a wrapper around
// gtk_dialog_get_response_for_widget().
func (v *Dialog) GetResponseForWidget(widget IWidget) ResponseType {
	c := C.gtk_dialog_get_response_for_widget(v.Native(), widget.toWidget())
	return ResponseType(c)
}

// WidgetForResponse() is a wrapper around
// gtk_dialog_get_widget_for_response().
func (v *Dialog) GetWidgetForResponse(id ResponseType) (*Widget, error) {
	c := C.gtk_dialog_get_widget_for_response(v.Native(), C.gint(id))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// ActionArea() is a wrapper around gtk_dialog_get_action_area().
func (v *Dialog) ActionArea() (*Widget, error) {
	c := C.gtk_dialog_get_action_area(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// ContentArea() is a wrapper around gtk_dialog_get_content_area().
func (v *Dialog) ContentArea() (*Box, error) {
	c := C.gtk_dialog_get_content_area(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	b := wrapBox(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &b, nil
}

// TODO(jrick)
/*
func (v *gdk.Screen) AlternativeDialogButtonOrder() bool {
	c := C.gtk_alternative_dialog_button_order(v.Native())
	return gobool(c)
}
*/

// TODO(jrick)
/*
func SetAlternativeButtonOrder(ids ...ResponseType) {
}
*/

/*
 * GtkEntry
 */

// Entry is a representation of GTK's GtkEntry.
type Entry struct {
	Widget
}

var entryType = glib.Type(C.gtk_entry_get_type())

func GetEntryType() glib.Type {
	return glib.Type(entryType)
}

// Native() returns a pointer to the underlying GtkEntry.
func (v *Entry) Native() *C.GtkEntry {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(entryType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkEntry)(v.Ptr())
}

func wrapEntry(obj *glib.Object) (e Entry) {
	e.Widget = wrapWidget(obj)
	return
}

// EntryNew() is a wrapper around gtk_entry_new().
func EntryNew() (*Entry, error) {
	c := C.gtk_entry_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	e := wrapEntry(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &e, nil
}

// EntryNewWithBuffer() is a wrapper around gtk_entry_new_with_buffer().
func EntryNewWithBuffer(buffer *EntryBuffer) (*Entry, error) {
	c := C.gtk_entry_new_with_buffer(buffer.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	e := wrapEntry(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &e, nil
}

// Buffer() is a wrapper around gtk_entry_get_buffer().
func (v *Entry) Buffer() (*EntryBuffer, error) {
	c := C.gtk_entry_get_buffer(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	e := EntryBuffer{obj}
	obj.Ref()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &e, nil
}

// SetBuffer() is a wrapper around gtk_entry_set_buffer().
func (v *Entry) SetBuffer(buffer *EntryBuffer) {
	C.gtk_entry_set_buffer(v.Native(), buffer.Native())
}

// SetText() is a wrapper around gtk_entry_set_text().
func (v *Entry) SetText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_entry_set_text(v.Native(), (*C.gchar)(cstr))
}

// Text() is a wrapper around gtk_entry_get_text().
func (v *Entry) Text() (string, error) {
	c := C.gtk_entry_get_text(v.Native())
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// TextLength() is a wrapper around gtk_entry_get_text_length().
func (v *Entry) TextLength() uint16 {
	c := C.gtk_entry_get_text_length(v.Native())
	return uint16(c)
}

// TODO(jrick) GdkRectangle
/*
func (v *Entry) TextArea() {
}
*/

// SetVisibility() is a wrapper around gtk_entry_set_visibility().
func (v *Entry) SetVisibility(visible bool) {
	C.gtk_entry_set_visibility(v.Native(), gbool(visible))
}

// SetInvisibleChar() is a wrapper around gtk_entry_set_invisible_char().
func (v *Entry) SetInvisibleChar(ch rune) {
	C.gtk_entry_set_invisible_char(v.Native(), C.gunichar(ch))
}

// UnsetInvisibleChar() is a wrapper around gtk_entry_unset_invisible_char().
func (v *Entry) UnsetInvisibleChar() {
	C.gtk_entry_unset_invisible_char(v.Native())
}

// SetMaxLength() is a wrapper around gtk_entry_set_max_length().
func (v *Entry) SetMaxLength(len int) {
	C.gtk_entry_set_max_length(v.Native(), C.gint(len))
}

// ActivatesDefault() is a wrapper around gtk_entry_get_activates_default().
func (v *Entry) ActivatesDefault() bool {
	c := C.gtk_entry_get_activates_default(v.Native())
	return gobool(c)
}

// HasFrame() is a wrapper around gtk_entry_get_has_frame().
func (v *Entry) HasFrame() bool {
	c := C.gtk_entry_get_has_frame(v.Native())
	return gobool(c)
}

// WidthChars() is a wrapper around gtk_entry_get_width_chars().
func (v *Entry) WidthChars() int {
	c := C.gtk_entry_get_width_chars(v.Native())
	return int(c)
}

// SetActivatesDefault() is a wrapper around gtk_entry_set_activates_default().
func (v *Entry) SetActivatesDefault(setting bool) {
	C.gtk_entry_set_activates_default(v.Native(), gbool(setting))
}

// SetHasFrame() is a wrapper around gtk_entry_set_has_frame().
func (v *Entry) SetHasFrame(setting bool) {
	C.gtk_entry_set_has_frame(v.Native(), gbool(setting))
}

// SetWidthChars() is a wrapper around gtk_entry_set_width_chars().
func (v *Entry) SetWidthChars(nChars int) {
	C.gtk_entry_set_width_chars(v.Native(), C.gint(nChars))
}

// InvisibleChar() is a wrapper around gtk_entry_get_invisible_char().
func (v *Entry) InvisibleChar() rune {
	c := C.gtk_entry_get_invisible_char(v.Native())
	return rune(c)
}

// SetAlignment() is a wrapper around gtk_entry_set_alignment().
func (v *Entry) SetAlignment(xalign float32) {
	C.gtk_entry_set_alignment(v.Native(), C.gfloat(xalign))
}

// Alignment() is a wrapper around gtk_entry_get_alignment().
func (v *Entry) Alignment() float32 {
	c := C.gtk_entry_get_alignment(v.Native())
	return float32(c)
}

// SetPlaceholderText() is a wrapper around gtk_entry_set_placeholder_text().
func (v *Entry) SetPlaceholderText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_entry_set_placeholder_text(v.Native(), (*C.gchar)(cstr))
}

// PlaceholderText() is a wrapper around gtk_entry_get_placeholder_text().
func (v *Entry) PlaceholderText() (string, error) {
	c := C.gtk_entry_get_placeholder_text(v.Native())
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// SetOverwriteMode() is a wrapper around gtk_entry_set_overwrite_mode().
func (v *Entry) SetOverwriteMode(overwrite bool) {
	C.gtk_entry_set_overwrite_mode(v.Native(), gbool(overwrite))
}

// OverwriteMode() is a wrapper around gtk_entry_get_overwrite_mode().
func (v *Entry) OverwriteMode() bool {
	c := C.gtk_entry_get_overwrite_mode(v.Native())
	return gobool(c)
}

// TODO(jrick) Pangolayout
/*
func (v *Entry) Layout() {
}
*/

// LayoutOffsets() is a wrapper around gtk_entry_get_layout_offsets().
func (v *Entry) LayoutOffsets() (x, y int) {
	var gx, gy C.gint
	C.gtk_entry_get_layout_offsets(v.Native(), &gx, &gy)
	return int(gx), int(gy)
}

// LayoutIndexToTextIndex() is a wrapper around
// gtk_entry_layout_index_to_text_index().
func (v *Entry) LayoutIndexToTextIndex(layoutIndex int) int {
	c := C.gtk_entry_layout_index_to_text_index(v.Native(),
		C.gint(layoutIndex))
	return int(c)
}

// TextIndexToLayoutIndex() is a wrapper around
// gtk_entry_text_index_to_layout_index().
func (v *Entry) TextIndexToLayoutIndex(textIndex int) int {
	c := C.gtk_entry_text_index_to_layout_index(v.Native(),
		C.gint(textIndex))
	return int(c)
}

// TODO(jrick) PandoAttrList
/*
func (v *Entry) SetAttributes() {
}
*/

// TODO(jrick) PandoAttrList
/*
func (v *Entry) Attributes() {
}
*/

// MaxLength() is a wrapper around gtk_entry_get_max_length().
func (v *Entry) MaxLength() int {
	c := C.gtk_entry_get_max_length(v.Native())
	return int(c)
}

// Visibility() is a wrapper around gtk_entry_get_visibility().
func (v *Entry) Visibility() bool {
	c := C.gtk_entry_get_visibility(v.Native())
	return gobool(c)
}

// SetCompletion() is a wrapper around gtk_entry_set_completion().
func (v *Entry) SetCompletion(completion *EntryCompletion) {
	C.gtk_entry_set_completion(v.Native(), completion.Native())
}

// Completion() is a wrapper around gtk_entry_get_completion().
func (v *Entry) Completion() (*EntryCompletion, error) {
	c := C.gtk_entry_get_completion(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	e := &EntryCompletion{obj}
	obj.Ref()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return e, nil
}

// SetCursorHAdjustment() is a wrapper around
// gtk_entry_set_cursor_hadjustment().
func (v *Entry) SetCursorHAdjustment(adjustment *Adjustment) {
	C.gtk_entry_set_cursor_hadjustment(v.Native(), adjustment.Native())
}

// CursorHAdjustment() is a wrapper around
// gtk_entry_get_cursor_hadjustment().
func (v *Entry) CursorHAdjustment() (*Adjustment, error) {
	c := C.gtk_entry_get_cursor_hadjustment(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	a := &Adjustment{glib.InitiallyUnowned{*obj}}
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return a, nil
}

// SetProgressFraction() is a wrapper around gtk_entry_set_progress_fraction().
func (v *Entry) SetProgressFraction(fraction float64) {
	C.gtk_entry_set_progress_fraction(v.Native(), C.gdouble(fraction))
}

// ProgressFraction() is a wrapper around gtk_entry_get_progress_fraction().
func (v *Entry) ProgressFraction() float64 {
	c := C.gtk_entry_get_progress_fraction(v.Native())
	return float64(c)
}

// SetProgressPulseStep() is a wrapper around
// gtk_entry_set_progress_pulse_step().
func (v *Entry) SetProgressPulseStep(fraction float64) {
	C.gtk_entry_set_progress_pulse_step(v.Native(), C.gdouble(fraction))
}

// ProgressPulseStep() is a wrapper around
// gtk_entry_get_progress_pulse_step().
func (v *Entry) ProgressPulseStep() float64 {
	c := C.gtk_entry_get_progress_pulse_step(v.Native())
	return float64(c)
}

// ProgressPulse() is a wrapper around gtk_entry_progress_pulse().
func (v *Entry) ProgressPulse() {
	C.gtk_entry_progress_pulse(v.Native())
}

// TODO(jrick) GdkEventKey
/*
func (v *Entry) IMContextFilterKeypress() {
}
*/

// ResetIMContext() is a wrapper around gtk_entry_reset_im_context().
func (v *Entry) ResetIMContext() {
	C.gtk_entry_reset_im_context(v.Native())
}

// TODO(jrick) GdkPixbuf
/*
func (v *Entry) SetIconFromPixbuf() {
}
*/

// SetIconFromStock() is a wrapper around gtk_entry_set_icon_from_stock().
func (v *Entry) SetIconFromStock(iconPos EntryIconPosition, stockID string) {
	cstr := C.CString(stockID)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_entry_set_icon_from_stock(v.Native(),
		C.GtkEntryIconPosition(iconPos), (*C.gchar)(cstr))
}

// SetIconFromIconName() is a wrapper around
// gtk_entry_set_icon_from_icon_name().
func (v *Entry) SetIconFromIconName(iconPos EntryIconPosition, name string) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_entry_set_icon_from_icon_name(v.Native(),
		C.GtkEntryIconPosition(iconPos), (*C.gchar)(cstr))
}

// TODO(jrick) GIcon
/*
func (v *Entry) SetIconFromGIcon() {
}
*/

// IconStorageType() is a wrapper around gtk_entry_get_icon_storage_type().
func (v *Entry) GetIconStorageType(iconPos EntryIconPosition) ImageType {
	c := C.gtk_entry_get_icon_storage_type(v.Native(),
		C.GtkEntryIconPosition(iconPos))
	return ImageType(c)
}

// TODO(jrick) GdkPixbuf
/*
func (v *Entry) IconPixbuf() {
}
*/

// IconStock() is a wrapper around gtk_entry_get_icon_stock().
func (v *Entry) GetIconStock(iconPos EntryIconPosition) (string, error) {
	c := C.gtk_entry_get_icon_stock(v.Native(),
		C.GtkEntryIconPosition(iconPos))
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// IconName() is a wrapper around gtk_entry_get_icon_name().
func (v *Entry) GetIconName(iconPos EntryIconPosition) (string, error) {
	c := C.gtk_entry_get_icon_name(v.Native(),
		C.GtkEntryIconPosition(iconPos))
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// TODO(jrick) GIcon
/*
func (v *Entry) IconGIcon() {
}
*/

// SetIconActivatable() is a wrapper around gtk_entry_set_icon_activatable().
func (v *Entry) SetIconActivatable(iconPos EntryIconPosition, activatable bool) {
	C.gtk_entry_set_icon_activatable(v.Native(),
		C.GtkEntryIconPosition(iconPos), gbool(activatable))
}

// IconActivatable() is a wrapper around gtk_entry_get_icon_activatable().
func (v *Entry) GetIconActivatable(iconPos EntryIconPosition) bool {
	c := C.gtk_entry_get_icon_activatable(v.Native(),
		C.GtkEntryIconPosition(iconPos))
	return gobool(c)
}

// SetIconSensitive() is a wrapper around gtk_entry_set_icon_sensitive().
func (v *Entry) SetIconSensitive(iconPos EntryIconPosition, sensitive bool) {
	C.gtk_entry_set_icon_sensitive(v.Native(),
		C.GtkEntryIconPosition(iconPos), gbool(sensitive))
}

// IconSensitive() is a wrapper around gtk_entry_get_icon_sensitive().
func (v *Entry) GetIconSensitive(iconPos EntryIconPosition) bool {
	c := C.gtk_entry_get_icon_sensitive(v.Native(),
		C.GtkEntryIconPosition(iconPos))
	return gobool(c)
}

// IconAtPos() is a wrapper around gtk_entry_get_icon_at_pos().
func (v *Entry) GetIconAtPos(x, y int) int {
	c := C.gtk_entry_get_icon_at_pos(v.Native(), C.gint(x), C.gint(y))
	return int(c)
}

// SetIconTooltipText() is a wrapper around gtk_entry_set_icon_tooltip_text().
func (v *Entry) SetIconTooltipText(iconPos EntryIconPosition, tooltip string) {
	cstr := C.CString(tooltip)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_entry_set_icon_tooltip_text(v.Native(),
		C.GtkEntryIconPosition(iconPos), (*C.gchar)(cstr))
}

// IconTooltipText() is a wrapper around gtk_entry_get_icon_tooltip_text().
func (v *Entry) GetIconTooltipText(iconPos EntryIconPosition) (string, error) {
	c := C.gtk_entry_get_icon_tooltip_text(v.Native(),
		C.GtkEntryIconPosition(iconPos))
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// SetIconTooltipMarkup() is a wrapper around
// gtk_entry_set_icon_tooltip_markup().
func (v *Entry) SetIconTooltipMarkup(iconPos EntryIconPosition, tooltip string) {
	cstr := C.CString(tooltip)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_entry_set_icon_tooltip_markup(v.Native(),
		C.GtkEntryIconPosition(iconPos), (*C.gchar)(cstr))
}

// IconTooltipMarkup() is a wrapper around
// gtk_entry_get_icon_tooltip_markup().
func (v *Entry) GetIconTooltipMarkup(iconPos EntryIconPosition) (string, error) {
	c := C.gtk_entry_get_icon_tooltip_markup(v.Native(),
		C.GtkEntryIconPosition(iconPos))
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// TODO(jrick) GdkDragAction
/*
func (v *Entry) SetIconDragSource() {
}
*/

// CurrentIconDragSource() is a wrapper around
// gtk_entry_get_current_icon_drag_source().
func (v *Entry) CurrentIconDragSource() int {
	c := C.gtk_entry_get_current_icon_drag_source(v.Native())
	return int(c)
}

// TODO(jrick) GdkRectangle
/*
func (v *Entry) IconArea() {
}
*/

// SetInputPurpose() is a wrapper around gtk_entry_set_input_purpose().
func (v *Entry) SetInputPurpose(purpose InputPurpose) {
	C.gtk_entry_set_input_purpose(v.Native(), C.GtkInputPurpose(purpose))
}

// InputPurpose() is a wrapper around gtk_entry_get_input_purpose().
func (v *Entry) InputPurpose() InputPurpose {
	c := C.gtk_entry_get_input_purpose(v.Native())
	return InputPurpose(c)
}

// SetInputHints() is a wrapper around gtk_entry_set_input_hints().
func (v *Entry) SetInputHints(hints InputHints) {
	C.gtk_entry_set_input_hints(v.Native(), C.GtkInputHints(hints))
}

// InputHints() is a wrapper around gtk_entry_get_input_hints().
func (v *Entry) InputHints() InputHints {
	c := C.gtk_entry_get_input_hints(v.Native())
	return InputHints(c)
}

/*
 * GtkEntryBuffer
 */

// EntryBuffer is a representation of GTK's GtkEntryBuffer.
type EntryBuffer struct {
	*glib.Object
}

var entryBufferType = glib.Type(C.gtk_entry_buffer_get_type())

func GetEntryBufferType() glib.Type {
	return glib.Type(entryBufferType)
}

// Native() returns a pointer to the underlying GtkEntryBuffer.
func (v *EntryBuffer) Native() *C.GtkEntryBuffer {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(entryBufferType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkEntryBuffer)(v.Ptr())
}

func wrapEntryBuffer(obj *glib.Object) EntryBuffer {
	return EntryBuffer{obj}
}

// EntryBufferNew() is a wrapper around gtk_entry_buffer_new().
func EntryBufferNew(initialChars string, nInitialChars int) (*EntryBuffer, error) {
	cstr := C.CString(initialChars)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_entry_buffer_new((*C.gchar)(cstr), C.gint(nInitialChars))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	e := wrapEntryBuffer(obj)
	obj.Ref()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &e, nil
}

// Text() is a wrapper around gtk_entry_buffer_get_text().  A
// non-nil error is returned in the case that gtk_entry_buffer_get_text
// returns NULL to differentiate between NULL and an empty string.
func (v *EntryBuffer) Text() (string, error) {
	c := C.gtk_entry_buffer_get_text(v.Native())
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// SetText() is a wrapper around gtk_entry_buffer_set_text().
func (v *EntryBuffer) SetText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_entry_buffer_set_text(v.Native(), (*C.gchar)(cstr),
		C.gint(len(text)))
}

// Bytes() is a wrapper around gtk_entry_buffer_get_bytes().
func (v *EntryBuffer) Bytes() uint {
	c := C.gtk_entry_buffer_get_bytes(v.Native())
	return uint(c)
}

// Length() is a wrapper around gtk_entry_buffer_get_length().
func (v *EntryBuffer) Length() uint {
	c := C.gtk_entry_buffer_get_length(v.Native())
	return uint(c)
}

// MaxLength() is a wrapper around gtk_entry_buffer_get_max_length().
func (v *EntryBuffer) MaxLength() int {
	c := C.gtk_entry_buffer_get_max_length(v.Native())
	return int(c)
}

// SetMaxLength() is a wrapper around gtk_entry_buffer_set_max_length().
func (v *EntryBuffer) SetMaxLength(maxLength int) {
	C.gtk_entry_buffer_set_max_length(v.Native(), C.gint(maxLength))
}

// InsertText() is a wrapper around gtk_entry_buffer_insert_text().
func (v *EntryBuffer) InsertText(position uint, text string) uint {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_entry_buffer_insert_text(v.Native(), C.guint(position),
		(*C.gchar)(cstr), C.gint(len(text)))
	return uint(c)
}

// DeleteText() is a wrapper around gtk_entry_buffer_delete_text().
func (v *EntryBuffer) DeleteText(position uint, nChars int) uint {
	c := C.gtk_entry_buffer_delete_text(v.Native(), C.guint(position),
		C.gint(nChars))
	return uint(c)
}

// EmitDeletedText() is a wrapper around gtk_entry_buffer_emit_deleted_text().
func (v *EntryBuffer) EmitDeletedText(pos, nChars uint) {
	C.gtk_entry_buffer_emit_deleted_text(v.Native(), C.guint(pos),
		C.guint(nChars))
}

// EmitInsertedText() is a wrapper around gtk_entry_buffer_emit_inserted_text().
func (v *EntryBuffer) EmitInsertedText(pos uint, text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_entry_buffer_emit_inserted_text(v.Native(), C.guint(pos),
		(*C.gchar)(cstr), C.guint(len(text)))
}

/*
 * GtkEntryCompletion
 */

// EntryCompletion is a representation of GTK's GtkEntryCompletion.
type EntryCompletion struct {
	*glib.Object
}

var entryCompletionType = glib.Type(C.gtk_entry_completion_get_type())

func GetEntryCompletionType() glib.Type {
	return glib.Type(entryCompletionType)
}

// Native() returns a pointer to the underlying GtkEntryCompletion.
func (v *EntryCompletion) Native() *C.GtkEntryCompletion {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(entryCompletionType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkEntryCompletion)(v.Ptr())
}

func wrapEntryCompletion(obj *glib.Object) EntryCompletion {
	return EntryCompletion{obj}
}

/*
 * GtkFileChooser
 */

// FileChooser is a representation of GTK's GtkFileChooser GInterface.
type FileChooser struct {
	ptr unsafe.Pointer
}

var fileChooserType = glib.Type(C.gtk_file_chooser_get_type())

func GetFileChooserType() glib.Type {
	return glib.Type(fileChooserType)
}

func (v *FileChooser) Native() *C.GtkFileChooser {
	if v == nil {
		fmt.Println("nil object, not getting native file chooser")
		return nil
	}
	return (*C.GtkFileChooser)(v.ptr)
}

func (f *FileChooser) SetCurrentFolder(filename string) {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_file_chooser_set_filename(f.Native(), cstr)
}

func (f *FileChooser) Filename() string {
	c := C.gtk_file_chooser_get_filename(f.Native())
	if c == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(c))
	str := C.GoString((*C.char)(c))
	return str
}

/*
 * GtkFileChooserButton
 */
type FileChooserButton struct {
	Box

	// Interfaces
	FileChooser
}

var fileChooserButtonType = glib.Type(C.gtk_file_chooser_button_get_type())

func GetFileChooserButtonType() glib.Type {
	return glib.Type(fileChooserButtonType)
}

// Native() returns a pointer to the underlying GtkGrid.
func (v *FileChooserButton) Native() *C.GtkFileChooserButton {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(fileChooserButtonType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkFileChooserButton)(v.Ptr())
}

func wrapFileChooserButton(obj *glib.Object) (b FileChooserButton) {
	b.FileChooser = FileChooser{obj.Ptr()}
	b.Box = wrapBox(obj)
	return
}

/*
 * GtkGrid
 */

// Grid is a representation of GTK's GtkGrid.
type Grid struct {
	Container

	// Interfaces
	Orientable
}

var gridType = glib.Type(C.gtk_grid_get_type())

func GetGridType() glib.Type {
	return glib.Type(gridType)
}

// Native() returns a pointer to the underlying GtkGrid.
func (v *Grid) Native() *C.GtkGrid {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(gridType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkGrid)(v.Ptr())
}

func (v *Grid) toOrientable() *C.GtkOrientable {
	if v == nil {
		return nil
	}
	return (*C.GtkOrientable)(v.Ptr())
}

func wrapGrid(obj *glib.Object) (g Grid) {
	g.Orientable = Orientable{obj}
	g.Container = wrapContainer(obj)
	return
}

// GridNew() is a wrapper around gtk_grid_new().
func GridNew() (*Grid, error) {
	c := C.gtk_grid_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	g := wrapGrid(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &g, nil
}

// Attach() is a wrapper around gtk_grid_attach().
func (v *Grid) Attach(child IWidget, left, top, width, height int) {
	C.gtk_grid_attach(v.Native(), child.toWidget(), C.gint(left),
		C.gint(top), C.gint(width), C.gint(height))
}

// AttachNextTo() is a wrapper around gtk_grid_attach_next_to().
func (v *Grid) AttachNextTo(child, sibling IWidget, side PositionType, width, height int) {
	C.gtk_grid_attach_next_to(v.Native(), child.toWidget(),
		sibling.toWidget(), C.GtkPositionType(side), C.gint(width),
		C.gint(height))
}

// ChildAt() is a wrapper around gtk_grid_get_child_at().
func (v *Grid) GetChildAt(left, top int) (*Widget, error) {
	c := C.gtk_grid_get_child_at(v.Native(), C.gint(left), C.gint(top))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// InsertRow() is a wrapper around gtk_grid_insert_row().
func (v *Grid) InsertRow(position int) {
	C.gtk_grid_insert_row(v.Native(), C.gint(position))
}

// InsertColumn() is a wrapper around gtk_grid_insert_column().
func (v *Grid) InsertColumn(position int) {
	C.gtk_grid_insert_column(v.Native(), C.gint(position))
}

// InsertNextTo() is a wrapper around gtk_grid_insert_next_to()
func (v *Grid) InsertNextTo(sibling IWidget, side PositionType) {
	C.gtk_grid_insert_next_to(v.Native(), sibling.toWidget(),
		C.GtkPositionType(side))
}

// SetRowHomogeneous() is a wrapper around gtk_grid_set_row_homogeneous().
func (v *Grid) SetRowHomogeneous(homogeneous bool) {
	C.gtk_grid_set_row_homogeneous(v.Native(), gbool(homogeneous))
}

// RowHomogeneous() is a wrapper around gtk_grid_get_row_homogeneous().
func (v *Grid) RowHomogeneous() bool {
	c := C.gtk_grid_get_row_homogeneous(v.Native())
	return gobool(c)
}

// SetRowSpacing() is a wrapper around gtk_grid_set_row_spacing().
func (v *Grid) SetRowSpacing(spacing uint) {
	C.gtk_grid_set_row_spacing(v.Native(), C.guint(spacing))
}

// RowSpacing() is a wrapper around gtk_grid_get_row_spacing().
func (v *Grid) RowSpacing() uint {
	c := C.gtk_grid_get_row_spacing(v.Native())
	return uint(c)
}

// SetColumnHomogeneous() is a wrapper around gtk_grid_set_column_homogeneous().
func (v *Grid) SetColumnHomogeneous(homogeneous bool) {
	C.gtk_grid_set_column_homogeneous(v.Native(), gbool(homogeneous))
}

// ColumnHomogeneous() is a wrapper around gtk_grid_get_column_homogeneous().
func (v *Grid) ColumnHomogeneous() bool {
	c := C.gtk_grid_get_column_homogeneous(v.Native())
	return gobool(c)
}

// SetColumnSpacing() is a wrapper around gtk_grid_set_column_spacing().
func (v *Grid) SetColumnSpacing(spacing uint) {
	C.gtk_grid_set_column_spacing(v.Native(), C.guint(spacing))
}

// ColumnSpacing() is a wrapper around gtk_grid_get_column_spacing().
func (v *Grid) ColumnSpacing() uint {
	c := C.gtk_grid_get_column_spacing(v.Native())
	return uint(c)
}

/*
 * GtkImage
 */

// Image is a representation of GTK's GtkImage.
type Image struct {
	Misc
}

var imageType = glib.Type(C.gtk_image_get_type())

func GetImageType() glib.Type {
	return glib.Type(imageType)
}

// Native() returns a pointer to the underlying GtkImage.
func (v *Image) Native() *C.GtkImage {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(imageType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkImage)(v.Ptr())
}

func wrapImage(obj *glib.Object) (i Image) {
	i.Misc = wrapMisc(obj)
	return
}

// ImageNew() is a wrapper around gtk_image_new().
func ImageNew() (*Image, error) {
	c := C.gtk_image_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	i := wrapImage(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &i, nil
}

// ImageNewFromFile() is a wrapper around gtk_image_new_from_file().
func ImageNewFromFile(filename string) (*Image, error) {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_image_new_from_file((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	i := wrapImage(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &i, nil
}

// ImageNewFromResource() is a wrapper around gtk_image_new_from_resource().
func ImageNewFromResource(resourcePath string) (*Image, error) {
	cstr := C.CString(resourcePath)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_image_new_from_resource((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	i := wrapImage(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &i, nil
}

// TODO(jrick) GdkPixbuf
/*
func ImageNewFromPixbuf() {
}
*/

// ImageNewFromStock() is a wrapper around gtk_image_new_from_stock().
func ImageNewFromStock(stock Stock, size IconSize) (*Image, error) {
	cstr := C.CString(string(stock))
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_image_new_from_stock((*C.gchar)(cstr), C.GtkIconSize(size))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	i := wrapImage(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &i, nil
}

// TODO(jrick) GtkIconSet
/*
func ImageNewFromIconSet() {
}
*/

// TODO(jrick) GdkPixbufAnimation
/*
func ImageNewFromAnimation() {
}
*/

// ImageNewFromIconName() is a wrapper around gtk_image_new_from_icon_name().
func ImageNewFromIconName(iconName string, size IconSize) (*Image, error) {
	cstr := C.CString(iconName)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_image_new_from_icon_name((*C.gchar)(cstr),
		C.GtkIconSize(size))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	i := wrapImage(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &i, nil
}

// TODO(jrick) GIcon
/*
func ImageNewFromGIcon() {
}
*/

// Clear() is a wrapper around gtk_image_clear().
func (v *Image) Clear() {
	C.gtk_image_clear(v.Native())
}

// SetFromFile() is a wrapper around gtk_image_set_from_file().
func (v *Image) SetFromFile(filename string) {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_image_set_from_file(v.Native(), (*C.gchar)(cstr))
}

// SetFromResource() is a wrapper around gtk_image_set_from_resource().
func (v *Image) SetFromResource(resourcePath string) {
	cstr := C.CString(resourcePath)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_image_set_from_resource(v.Native(), (*C.gchar)(cstr))
}

// TODO(jrick) GdkPixbuf
/*
func (v *Image) SetFromPixbuf() {
}
*/

// SetFromStock() is a wrapper around gtk_image_set_from_stock().
func (v *Image) SetFromStock(stock Stock, size IconSize) {
	cstr := C.CString(string(stock))
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_image_set_from_stock(v.Native(), (*C.gchar)(cstr),
		C.GtkIconSize(size))
}

// TODO(jrick) GtkIconSet
/*
func (v *Image) SetFromIconSet() {
}
*/

// TODO(jrick) GdkPixbufAnimation
/*
func (v *Image) SetFromAnimation() {
}
*/

// SetFromIconName() is a wrapper around gtk_image_set_from_icon_name().
func (v *Image) SetFromIconName(iconName string, size IconSize) {
	cstr := C.CString(iconName)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_image_set_from_icon_name(v.Native(), (*C.gchar)(cstr),
		C.GtkIconSize(size))
}

// TODO(jrick) GIcon
/*
func (v *Image) SetFromGIcon() {
}
*/

// SetPixelSize() is a wrapper around gtk_image_set_pixel_size().
func (v *Image) SetPixelSize(pixelSize int) {
	C.gtk_image_set_pixel_size(v.Native(), C.gint(pixelSize))
}

// StorageType() is a wrapper around gtk_image_get_storage_type().
func (v *Image) StorageType() ImageType {
	c := C.gtk_image_get_storage_type(v.Native())
	return ImageType(c)
}

// TODO(jrick) GdkPixbuf
/*
func (v *Image) Pixbuf() {
}
*/

// TODO(jrick) GtkIconSet
/*
func (v *Image) IconSet() {
}
*/

// TODO(jrick) GdkPixbufAnimation
/*
func (v *Image) Animation() {
}
*/

// IconName() is a wrapper around gtk_image_get_icon_name().
func (v *Image) IconName() (string, IconSize) {
	var iconName *C.gchar
	var size C.GtkIconSize
	C.gtk_image_get_icon_name(v.Native(), &iconName, &size)
	return C.GoString((*C.char)(iconName)), IconSize(size)
}

// TODO(jrick) GIcon
/*
func (v *Image) GIcon() {
}
*/

// PixelSize() is a wrapper around gtk_image_get_pixel_size().
func (v *Image) PixelSize() int {
	c := C.gtk_image_get_pixel_size(v.Native())
	return int(c)
}

/*
 * GtkImageMenuItem
 */

// ImageMenuItem is a representation of GTK's GtkImageMenuItem.
type ImageMenuItem struct {
	MenuItem
}

var imageMenuItemType = glib.Type(C.gtk_image_menu_item_get_type())

func GetImageMenuItemType() glib.Type {
	return glib.Type(imageMenuItemType)
}

// Native() returns a pointer to the underlying GtkImageMenuItem.
func (v *ImageMenuItem) Native() *C.GtkImageMenuItem {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(imageMenuItemType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkImageMenuItem)(v.Ptr())
}

func wrapImageMenuItem(obj *glib.Object) (m ImageMenuItem) {
	m.MenuItem = wrapMenuItem(obj)
	return
}

/*
 * GtkLabel
 */

// Label is a representation of GTK's GtkLabel.
type Label struct {
	Misc
}

var labelType = glib.Type(C.gtk_label_get_type())

func GetLabelType() glib.Type {
	return glib.Type(labelType)
}

// Native() returns a pointer to the underlying GtkLabel.
func (v *Label) Native() *C.GtkLabel {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(labelType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkLabel)(v.Ptr())
}

func wrapLabel(obj *glib.Object) (l Label) {
	l.Misc = wrapMisc(obj)
	return
}

// LabelNew() is a wrapper around gtk_label_new().
func LabelNew(str string) (*Label, error) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_label_new((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	l := wrapLabel(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &l, nil
}

// SetText() is a wrapper around gtk_label_set_text().
func (v *Label) SetText(str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_label_set_text(v.Native(), (*C.gchar)(cstr))
}

// SetMarkup() is a wrapper around gtk_label_set_markup().
func (v *Label) SetMarkup(str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_label_set_markup(v.Native(), (*C.gchar)(cstr))
}

// SetMarkupWithMnemonic() is a wrapper around
// gtk_label_set_markup_with_mnemonic().
func (v *Label) SetMarkupWithMnemonic(str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_label_set_markup_with_mnemonic(v.Native(), (*C.gchar)(cstr))
}

// SetPattern() is a wrapper around gtk_label_set_pattern().
func (v *Label) SetPattern(patern string) {
	cstr := C.CString(patern)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_label_set_pattern(v.Native(), (*C.gchar)(cstr))
}

// SetWidthChars() is a wrapper around gtk_label_set_width_chars().
func (v *Label) SetWidthChars(nChars int) {
	C.gtk_label_set_width_chars(v.Native(), C.gint(nChars))
}

// SetMaxWidthChars() is a wrapper around gtk_label_set_max_width_chars().
func (v *Label) SetMaxWidthChars(nChars int) {
	C.gtk_label_set_max_width_chars(v.Native(), C.gint(nChars))
}

// SetLineWrap() is a wrapper around gtk_label_set_line_wrap().
func (v *Label) SetLineWrap(wrap bool) {
	C.gtk_label_set_line_wrap(v.Native(), gbool(wrap))
}

// Selectable() is a wrapper around gtk_label_get_selectable().
func (v *Label) Selectable() bool {
	c := C.gtk_label_get_selectable(v.Native())
	return gobool(c)
}

// Text() is a wrapper around gtk_label_get_text().
func (v *Label) Text() (string, error) {
	c := C.gtk_label_get_text(v.Native())
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// LabelNewWithMnemonic() is a wrapper around gtk_label_new_with_mnemonic().
func LabelNewWithMnemonic(str string) (*Label, error) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_label_new_with_mnemonic((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	l := wrapLabel(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &l, nil
}

// SetSelectable() is a wrapper around gtk_label_set_selectable().
func (v *Label) SetSelectable(setting bool) {
	C.gtk_label_set_selectable(v.Native(), gbool(setting))
}

// SetLabel() is a wrapper around gtk_label_set_label().
func (v *Label) SetLabel(str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_label_set_label(v.Native(), (*C.gchar)(cstr))
}

/*
 * GtkListStore
 */

// ListStore is a representation of GTK's GtkListStore.
type ListStore struct {
	*glib.Object

	// Interfaces
	TreeModel
}

var listStoreType = glib.Type(C.gtk_list_store_get_type())

func GetListStoreType() glib.Type {
	return glib.Type(listStoreType)
}

// Native() returns a pointer to the underlying GtkListStore.
func (v *ListStore) Native() *C.GtkListStore {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(listStoreType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkListStore)(v.Ptr())
}

func wrapListStore(obj *glib.Object) (l ListStore) {
	l.TreeModel = wrapTreeModel(obj)
	l.Object = obj
	return
}

func (v *ListStore) toTreeModel() *C.GtkTreeModel {
	if v == nil {
		return nil
	}
	return (*C.GtkTreeModel)(v.Ptr())
}

func (v *ListStore) ColumnIndex(name string) int {
	return v.indexMap[name]
}

// ListStoreNew() is a wrapper around gtk_list_store_newv().
func ListStoreNew(columns map[string]glib.Type) (*ListStore, error) {
	gtypes := C.alloc_types(C.int(len(columns)))
	indexMap := make(map[string]int)
	n := 0
	for name, val := range columns {
		C.set_type(gtypes, C.int(n), C.GType(val))
		indexMap[name] = n
		n++
	}
	defer C.g_free(C.gpointer(gtypes))
	c := C.gtk_list_store_newv(C.gint(len(columns)), gtypes)
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	ls := wrapListStore(obj)
	ls.indexMap = indexMap
	obj.Ref()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &ls, nil
}

// TODO(jrick)
/*
func (v *ListStore) SetColumnTypes(types ...glib.Type) {
}
*/

// Set() is a wrapper around gtk_list_store_set_valuesv()
func (v *ListStore) Set(iter *TreeIter, values map[string]interface{}) error {
	n := len(values)
	i := 0
	ccolumns := make([]C.gint, n)
	cvalues := make([]*C.GValue, n)
	for key, val := range values {
		index, ok := v.indexMap[key]
		if !ok {
			return fmt.Errorf("unrecognized key: '%s'", key)
		}
		v, err := glib.GValue(val)
		if err != nil {
			return err
		}
		ccolumns[i] = C.gint(index)
		// for some reason, this cast is necessary since GValue is defined in the glib package
		cvalues[i] = (*C.GValue)(unsafe.Pointer(v.Native()))
		i++
	}
	var (
		cn    = C.gint(n)
		ccols = (*C.gint)(unsafe.Pointer(&ccolumns))
		cvals = (*C.GValue)(unsafe.Pointer(&cvalues))
	)
	C.gtk_list_store_set_valuesv(v.Native(), iter.Native(), ccols, cvals, cn)
	return nil
}

// InsertWithValues() is a wrapper around gtk_list_store_insert_with_valuesv()
func (v *ListStore) InsertWithValues(position int, values map[string]interface{}) (*TreeIter, error) {
	n := len(values)
	i := 0
	ccolumns := make([]C.gint, n)
	cvalues := make([]C.GValue, n)
	for key, val := range values {
		index, ok := v.indexMap[key]
		if !ok {
			return nil, fmt.Errorf("unrecognized key: '%s'", key)
		}
		v, err := glib.GValue(val)
		if err != nil {
			return nil, err
		}
		ccolumns[i] = C.gint(index)
		// for some reason, this cast is necessary since GValue is defined in the glib package
		cvalues[i] = *(*C.GValue)(unsafe.Pointer(v.Native()))
		i++
	}
	var (
		citer C.GtkTreeIter
		cpos  = C.gint(position)
		cn    = C.gint(n)
		ccols = (*C.gint)(unsafe.Pointer(&ccolumns[0]))
		cvals = (*C.GValue)(unsafe.Pointer(&cvalues[0]))
	)
	C.gtk_list_store_insert_with_valuesv(v.Native(), &citer, cpos, ccols, cvals, cn)
	if &citer == nil {
		return nil, nilPtrErr
	}
	return &TreeIter{citer}, nil
}

// Prepend() is a wrapper around gtk_list_store_prepend().
func (v *ListStore) Prepend(iter *TreeIter) {
	C.gtk_list_store_prepend(v.Native(), iter.Native())
}

// Append() is a wrapper around gtk_list_store_append().
func (v *ListStore) Append(iter *TreeIter) {
	C.gtk_list_store_append(v.Native(), iter.Native())
}

// Clear() is a wrapper around gtk_list_store_clear().
func (v *ListStore) Clear() {
	C.gtk_list_store_clear(v.Native())
}

// IterIsValid() is a wrapper around gtk_list_store_iter_is_valid().
func (v *ListStore) IterIsValid(iter *TreeIter) bool {
	c := C.gtk_list_store_iter_is_valid(v.Native(), iter.Native())
	return gobool(c)
}

// TODO(jrick)
/*
func (v *ListStore) Reorder(newOrder []int) {
}
*/

// Swap() is a wrapper around gtk_list_store_swap().
func (v *ListStore) Swap(a, b *TreeIter) {
	C.gtk_list_store_swap(v.Native(), a.Native(), b.Native())
}

// MoveBefore() is a wrapper around gtk_list_store_move_before().
func (v *ListStore) MoveBefore(iter, position *TreeIter) {
	C.gtk_list_store_move_before(v.Native(), iter.Native(),
		position.Native())
}

// MoveAfter() is a wrapper around gtk_list_store_move_after().
func (v *ListStore) MoveAfter(iter, position *TreeIter) {
	C.gtk_list_store_move_after(v.Native(), iter.Native(),
		position.Native())
}

/*
 * GtkMenu
 */

// Menu is a representation of GTK's GtkMenu.
type Menu struct {
	MenuShell
}

var menuType = glib.Type(C.gtk_menu_get_type())

func GetMenuType() glib.Type {
	return glib.Type(menuType)
}

// Native() returns a pointer to the underlying GtkMenu.
func (v *Menu) Native() *C.GtkMenu {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(menuType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkMenu)(v.Ptr())
}

func wrapMenu(obj *glib.Object) (m Menu) {
	m.MenuShell = wrapMenuShell(obj)
	return
}

// MenuNew() is a wrapper around gtk_menu_new().
func MenuNew() (*Menu, error) {
	c := C.gtk_menu_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	m := wrapMenu(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &m, nil
}

/*
 * GtkMenuBar
 */

// MenuBar is a representation of GTK's GtkMenuBar.
type MenuBar struct {
	MenuShell
}

var menuBarType = glib.Type(C.gtk_menu_bar_get_type())

func GetMenuBarType() glib.Type {
	return glib.Type(menuBarType)
}

// Native() returns a pointer to the underlying GtkMenuBar.
func (v *MenuBar) Native() *C.GtkMenuBar {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(menuBarType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkMenuBar)(v.Ptr())
}

func wrapMenuBar(obj *glib.Object) (m MenuBar) {
	m.MenuShell = wrapMenuShell(obj)
	return
}

// MenuBarNew() is a wrapper around gtk_menu_bar_new().
func MenuBarNew() (*MenuBar, error) {
	c := C.gtk_menu_bar_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	m := wrapMenuBar(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &m, nil
}

/*
 * GtkMenuItem
 */

// MenuItem is a representation of GTK's GtkMenuItem.
type MenuItem struct {
	Bin
}

var menuItemType = glib.Type(C.gtk_menu_item_get_type())

func GetMenuItemType() glib.Type {
	return glib.Type(menuItemType)
}

// Native() returns a pointer to the underlying GtkMenuItem.
func (v *MenuItem) Native() *C.GtkMenuItem {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(menuItemType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkMenuItem)(v.Ptr())
}

func wrapMenuItem(obj *glib.Object) (m MenuItem) {
	m.Bin = wrapBin(obj)
	return
}

// MenuItemNew() is a wrapper around gtk_menu_item_new().
func MenuItemNew() (*MenuItem, error) {
	c := C.gtk_menu_item_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	m := wrapMenuItem(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &m, nil
}

// MenuItemNewWithLabel() is a wrapper around gtk_menu_item_new_with_label().
func MenuItemNewWithLabel(label string) (*MenuItem, error) {
	cstr := C.CString(label)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_menu_item_new_with_label((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	m := wrapMenuItem(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &m, nil
}

// MenuItemNewWithMnemonic() is a wrapper around
// gtk_menu_item_new_with_mnemonic().
func MenuItemNewWithMnemonic(label string) (*MenuItem, error) {
	cstr := C.CString(label)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_menu_item_new_with_mnemonic((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	m := wrapMenuItem(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &m, nil
}

// SetSubmenu() is a wrapper around gtk_menu_item_set_submenu().
func (v *MenuItem) SetSubmenu(submenu IWidget) {
	C.gtk_menu_item_set_submenu(v.Native(), submenu.toWidget())
}

/*
 * GtkMenuShell
 */

// MenuShell is a representation of GTK's GtkMenuShell.
type MenuShell struct {
	Container
}

var menuShellType = glib.Type(C.gtk_menu_shell_get_type())

func GetMenuShellType() glib.Type {
	return glib.Type(menuShellType)
}

// Native() returns a pointer to the underlying GtkMenuShell.
func (v *MenuShell) Native() *C.GtkMenuShell {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(menuShellType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkMenuShell)(v.Ptr())
}

func wrapMenuShell(obj *glib.Object) (m MenuShell) {
	m.Container = wrapContainer(obj)
	return
}

// Append() is a wrapper around gtk_menu_shell_append().
func (v *MenuShell) Append(child IWidget) {
	C.gtk_menu_shell_append(v.Native(), child.toWidget())
}

/*
 * GtkMessageDialog
 */

// MessageDialog is a representation of GTK's GtkMessageDialog.
type MessageDialog struct {
	Dialog
}

var messageDialogType = glib.Type(C.gtk_message_dialog_get_type())

func GetMessageDialogType() glib.Type {
	return glib.Type(messageDialogType)
}

// Native() returns a pointer to the underlying GtkMessageDialog.
func (v *MessageDialog) Native() *C.GtkMessageDialog {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(messageDialogType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkMessageDialog)(v.Ptr())
}

func wrapMessageDialog(obj *glib.Object) (m MessageDialog) {
	m.Dialog = wrapDialog(obj)
	return
}

// MessageDialogNew() is a wrapper around gtk_message_dialog_new().
// The text is created and formatted by the format specifier and any
// additional arguments.
func MessageDialogNew(parent IWindow, flags DialogFlags, mType MessageType, buttons ButtonsType, format string, a ...interface{}) *MessageDialog {
	s := fmt.Sprintf(format, a...)
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	var w *C.GtkWindow = nil
	if parent != nil {
		w = parent.toWindow()
	}
	c := C._gtk_message_dialog_new(w,
		C.GtkDialogFlags(flags), C.GtkMessageType(mType),
		C.GtkButtonsType(buttons), cstr)
	obj := glib.ObjectNew(unsafe.Pointer(c))
	m := wrapMessageDialog(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &m
}

/*
 * GtkMisc
 */

// Misc is a representation of GTK's GtkMisc.
type Misc struct {
	Widget
}

var miscType = glib.Type(C.gtk_misc_get_type())

func GetMiscType() glib.Type {
	return glib.Type(miscType)
}

// Native() returns a pointer to the underlying GtkMisc.
func (v *Misc) Native() *C.GtkMisc {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(miscType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkMisc)(v.Ptr())
}

func wrapMisc(obj *glib.Object) (m Misc) {
	m.Widget = wrapWidget(obj)
	return
}

/*
 * GtkNotebook
 */

// Notebook is a representation of GTK's GtkNotebook.
type Notebook struct {
	Container
}

var notebookType = glib.Type(C.gtk_notebook_get_type())

func GetNotebookType() glib.Type {
	return glib.Type(notebookType)
}

// Native() returns a pointer to the underlying GtkNotebook.
func (v *Notebook) Native() *C.GtkNotebook {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(notebookType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkNotebook)(v.Ptr())
}

func wrapNotebook(obj *glib.Object) (n Notebook) {
	n.Container = wrapContainer(obj)
	return
}

// NotebookNew() is a wrapper around gtk_notebook_new().
func NotebookNew() (*Notebook, error) {
	c := C.gtk_notebook_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	n := wrapNotebook(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &n, nil
}

// AppendPage() is a wrapper around gtk_notebook_append_page().
func (v *Notebook) AppendPage(child IWidget, tabLabel IWidget) int {
	c := C.gtk_notebook_append_page(v.Native(), child.toWidget(),
		tabLabel.toWidget())
	return int(c)
}

// AppendPageMenu() is a wrapper around gtk_notebook_append_page_menu().
func (v *Notebook) AppendPageMenu(child IWidget, tabLabel IWidget, menuLabel IWidget) int {
	c := C.gtk_notebook_append_page_menu(v.Native(), child.toWidget(),
		tabLabel.toWidget(), menuLabel.toWidget())
	return int(c)
}

// PrependPage() is a wrapper around gtk_notebook_prepend_page().
func (v *Notebook) PrependPage(child IWidget, tabLabel IWidget) int {
	c := C.gtk_notebook_prepend_page(v.Native(), child.toWidget(),
		tabLabel.toWidget())
	return int(c)
}

// PrependPageMenu() is a wrapper around gtk_notebook_prepend_page_menu().
func (v *Notebook) PrependPageMenu(child IWidget, tabLabel IWidget, menuLabel IWidget) int {
	c := C.gtk_notebook_prepend_page_menu(v.Native(), child.toWidget(),
		tabLabel.toWidget(), menuLabel.toWidget())
	return int(c)
}

// InsertPage() is a wrapper around gtk_notebook_insert_page().
func (v *Notebook) InsertPage(child IWidget, tabLabel IWidget, position int) int {
	c := C.gtk_notebook_insert_page(v.Native(), child.toWidget(),
		tabLabel.toWidget(), C.gint(position))
	return int(c)
}

// InsertPageMenu() is a wrapper around gtk_notebook_insert_page_menu().
func (v *Notebook) InsertPageMenu(child IWidget, tabLabel IWidget, menuLabel IWidget, position int) int {
	c := C.gtk_notebook_insert_page_menu(v.Native(), child.toWidget(),
		tabLabel.toWidget(), menuLabel.toWidget(), C.gint(position))
	return int(c)
}

// RemovePage() is a wrapper around gtk_notebook_remove_page().
func (v *Notebook) RemovePage(pageNum int) {
	C.gtk_notebook_remove_page(v.Native(), C.gint(pageNum))
}

// PageNum() is a wrapper around gtk_notebook_page_num().
func (v *Notebook) PageNum(child IWidget) int {
	c := C.gtk_notebook_page_num(v.Native(), child.toWidget())
	return int(c)
}

// NextPage() is a wrapper around gtk_notebook_next_page().
func (v *Notebook) NextPage() {
	C.gtk_notebook_next_page(v.Native())
}

// PrevPage() is a wrapper around gtk_notebook_prev_page().
func (v *Notebook) PrevPage() {
	C.gtk_notebook_prev_page(v.Native())
}

// ReorderChild() is a wrapper around gtk_notebook_reorder_child().
func (v *Notebook) ReorderChild(child IWidget, position int) {
	C.gtk_notebook_reorder_child(v.Native(), child.toWidget(),
		C.gint(position))
}

// SetTabPos() is a wrapper around gtk_notebook_set_tab_pos().
func (v *Notebook) SetTabPos(pos PositionType) {
	C.gtk_notebook_set_tab_pos(v.Native(), C.GtkPositionType(pos))
}

// SetShowTabs() is a wrapper around gtk_notebook_set_show_tabs().
func (v *Notebook) SetShowTabs(showTabs bool) {
	C.gtk_notebook_set_show_tabs(v.Native(), gbool(showTabs))
}

// SetShowBorder() is a wrapper around gtk_notebook_set_show_border().
func (v *Notebook) SetShowBorder(showBorder bool) {
	C.gtk_notebook_set_show_border(v.Native(), gbool(showBorder))
}

// SetScrollable() is a wrapper around gtk_notebook_set_scrollable().
func (v *Notebook) SetScrollable(scrollable bool) {
	C.gtk_notebook_set_scrollable(v.Native(), gbool(scrollable))
}

// PopupEnable() is a wrapper around gtk_notebook_popup_enable().
func (v *Notebook) PopupEnable() {
	C.gtk_notebook_popup_enable(v.Native())
}

// PopupDisable() is a wrapper around gtk_notebook_popup_disable().
func (v *Notebook) PopupDisable() {
	C.gtk_notebook_popup_disable(v.Native())
}

// CurrentPage() is a wrapper around gtk_notebook_get_current_page().
func (v *Notebook) CurrentPage() int {
	c := C.gtk_notebook_get_current_page(v.Native())
	return int(c)
}

// MenuLabel() is a wrapper around gtk_notebook_get_menu_label().
func (v *Notebook) GetMenuLabel(child IWidget) (*Widget, error) {
	c := C.gtk_notebook_get_menu_label(v.Native(), child.toWidget())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// NthPage() is a wrapper around gtk_notebook_get_nth_page().
func (v *Notebook) GetNthPage(pageNum int) (*Widget, error) {
	c := C.gtk_notebook_get_nth_page(v.Native(), C.gint(pageNum))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// NPages() is a wrapper around gtk_notebook_get_n_pages().
func (v *Notebook) NPages() int {
	c := C.gtk_notebook_get_n_pages(v.Native())
	return int(c)
}

// TabLabel() is a wrapper around gtk_notebook_get_tab_label().
func (v *Notebook) GetTabLabel(child IWidget) (*Widget, error) {
	c := C.gtk_notebook_get_tab_label(v.Native(), child.toWidget())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// SetMenuLabel() is a wrapper around gtk_notebook_set_menu_label().
func (v *Notebook) SetMenuLabel(child, menuLabel IWidget) {
	C.gtk_notebook_set_menu_label(v.Native(), child.toWidget(),
		menuLabel.toWidget())
}

// SetMenuLabelText() is a wrapper around gtk_notebook_set_menu_label_text().
func (v *Notebook) SetMenuLabelText(child IWidget, menuText string) {
	cstr := C.CString(menuText)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_notebook_set_menu_label_text(v.Native(), child.toWidget(),
		(*C.gchar)(cstr))
}

// SetTabLabel() is a wrapper around gtk_notebook_set_tab_label().
func (v *Notebook) SetTabLabel(child, tabLabel IWidget) {
	C.gtk_notebook_set_tab_label(v.Native(), child.toWidget(),
		tabLabel.toWidget())
}

// SetTabLabelText() is a wrapper around gtk_notebook_set_tab_label_text().
func (v *Notebook) SetTabLabelText(child IWidget, tabText string) {
	cstr := C.CString(tabText)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_notebook_set_tab_label_text(v.Native(), child.toWidget(),
		(*C.gchar)(cstr))
}

// SetTabReorderable() is a wrapper around gtk_notebook_set_tab_reorderable().
func (v *Notebook) SetTabReorderable(child IWidget, reorderable bool) {
	C.gtk_notebook_set_tab_reorderable(v.Native(), child.toWidget(),
		gbool(reorderable))
}

// SetTabDetachable() is a wrapper around gtk_notebook_set_tab_detachable().
func (v *Notebook) SetTabDetachable(child IWidget, detachable bool) {
	C.gtk_notebook_set_tab_detachable(v.Native(), child.toWidget(),
		gbool(detachable))
}

// MenuLabelText() is a wrapper around gtk_notebook_get_menu_label_text().
func (v *Notebook) GetMenuLabelText(child IWidget) (string, error) {
	c := C.gtk_notebook_get_menu_label_text(v.Native(), child.toWidget())
	if c == nil {
		return "", errors.New("No menu label for widget")
	}
	return C.GoString((*C.char)(c)), nil
}

// Scrollable() is a wrapper around gtk_notebook_get_scrollable().
func (v *Notebook) Scrollable() bool {
	c := C.gtk_notebook_get_scrollable(v.Native())
	return gobool(c)
}

// ShowBorder() is a wrapper around gtk_notebook_get_show_border().
func (v *Notebook) ShowBorder() bool {
	c := C.gtk_notebook_get_show_border(v.Native())
	return gobool(c)
}

// ShowTabs() is a wrapper around gtk_notebook_get_show_tabs().
func (v *Notebook) ShowTabs() bool {
	c := C.gtk_notebook_get_show_tabs(v.Native())
	return gobool(c)
}

// TabLabelText() is a wrapper around gtk_notebook_get_tab_label_text().
func (v *Notebook) GetTabLabelText(child IWidget) (string, error) {
	c := C.gtk_notebook_get_tab_label_text(v.Native(), child.toWidget())
	if c == nil {
		return "", errors.New("No tab label for widget")
	}
	return C.GoString((*C.char)(c)), nil
}

// TabPos() is a wrapper around gtk_notebook_get_tab_pos().
func (v *Notebook) TabPos() PositionType {
	c := C.gtk_notebook_get_tab_pos(v.Native())
	return PositionType(c)
}

// TabReorderable() is a wrapper around gtk_notebook_get_tab_reorderable().
func (v *Notebook) GetTabReorderable(child IWidget) bool {
	c := C.gtk_notebook_get_tab_reorderable(v.Native(), child.toWidget())
	return gobool(c)
}

// TabDetachable() is a wrapper around gtk_notebook_get_tab_detachable().
func (v *Notebook) GetTabDetachable(child IWidget) bool {
	c := C.gtk_notebook_get_tab_detachable(v.Native(), child.toWidget())
	return gobool(c)
}

// SetCurrentPage() is a wrapper around gtk_notebook_set_current_page().
func (v *Notebook) SetCurrentPage(pageNum int) {
	C.gtk_notebook_set_current_page(v.Native(), C.gint(pageNum))
}

// SetGroupName() is a wrapper around gtk_notebook_set_group_name().
func (v *Notebook) SetGroupName(groupName string) {
	cstr := C.CString(groupName)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_notebook_set_group_name(v.Native(), (*C.gchar)(cstr))
}

// GroupName() is a wrapper around gtk_notebook_get_group_name().
func (v *Notebook) GroupName() (string, error) {
	c := C.gtk_notebook_get_group_name(v.Native())
	if c == nil {
		return "", errors.New("No group name")
	}
	return C.GoString((*C.char)(c)), nil
}

// SetActionWidget() is a wrapper around gtk_notebook_set_action_widget().
func (v *Notebook) SetActionWidget(widget IWidget, packType PackType) {
	C.gtk_notebook_set_action_widget(v.Native(), widget.toWidget(),
		C.GtkPackType(packType))
}

// ActionWidget() is a wrapper around gtk_notebook_get_action_widget().
func (v *Notebook) GetActionWidget(packType PackType) (*Widget, error) {
	c := C.gtk_notebook_get_action_widget(v.Native(),
		C.GtkPackType(packType))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

/*
 * GtkOffscreenWindow
 */

// OffscreenWindow is a representation of GTK's GtkOffscreenWindow.
type OffscreenWindow struct {
	Bin
}

var offscreenWindowType = glib.Type(C.gtk_offscreen_window_get_type())

func GetOffscreenWindowType() glib.Type {
	return glib.Type(offscreenWindowType)
}

// Native() returns a pointer to the underlying GtkWindow.
func (v *OffscreenWindow) Native() *C.GtkOffscreenWindow {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(offscreenWindowType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkOffscreenWindow)(v.Ptr())
}

func wrapOffscreenWindow(obj *glib.Object) (o OffscreenWindow) {
	o.Bin = wrapBin(obj)
	return
}

/*
 * GtkOrientable
 */

// Orientable is a representation of GTK's GtkOrientable GInterface.
type Orientable struct {
	*glib.Object
}

var orientableType = glib.Type(C.gtk_orientable_get_type())

func GetOrientableType() glib.Type {
	return glib.Type(orientableType)
}

// IOrientable is an interface type implemented by all structs
// embedding an Orientable.  It is meant to be used as an argument type
// for wrapper functions that wrap around a C GTK function taking a
// GtkOrientable.
type IOrientable interface {
	toOrientable() *C.GtkOrientable
}

// Native returns a pointer to the underlying GObject as a GtkOrientable.
func (v *Orientable) Native() *C.GtkOrientable {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(orientableType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkOrientable)(v.Ptr())
}

func wrapOrientable(obj *glib.Object) Orientable {
	return Orientable{obj}
}

// Orientation() is a wrapper around gtk_orientable_get_orientation().
func (v *Orientable) Orientation() Orientation {
	c := C.gtk_orientable_get_orientation(v.Native())
	return Orientation(c)
}

// SetOrientation() is a wrapper around gtk_orientable_set_orientation().
func (v *Orientable) SetOrientation(orientation Orientation) {
	C.gtk_orientable_set_orientation(v.Native(),
		C.GtkOrientation(orientation))
}

/*
 * GtkProgressBar
 */

// ProgressBar is a representation of GTK's GtkProgressBar.
type ProgressBar struct {
	Widget
}

var progressBarType = glib.Type(C.gtk_progress_bar_get_type())

func GetProgressBarType() glib.Type {
	return glib.Type(progressBarType)
}

// Native() returns a pointer to the underlying GtkProgressBar.
func (v *ProgressBar) Native() *C.GtkProgressBar {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(progressBarType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkProgressBar)(v.Ptr())
}

func wrapProgressBar(obj *glib.Object) (p ProgressBar) {
	p.Widget = wrapWidget(obj)
	return
}

// ProgressBarNew() is a wrapper around gtk_progress_bar_new().
func ProgressBarNew() (*ProgressBar, error) {
	c := C.gtk_progress_bar_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	p := wrapProgressBar(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &p, nil
}

// SetFraction() is a wrapper around gtk_progress_bar_set_fraction().
func (v *ProgressBar) SetFraction(fraction float64) {
	C.gtk_progress_bar_set_fraction(v.Native(), C.gdouble(fraction))
}

// Fraction() is a wrapper around gtk_progress_bar_get_fraction().
func (v *ProgressBar) Fraction() float64 {
	c := C.gtk_progress_bar_get_fraction(v.Native())
	return float64(c)
}

// SetText() is a wrapper around gtk_progress_bar_set_text().
func (v *ProgressBar) SetText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_progress_bar_set_text(v.Native(), (*C.gchar)(cstr))
}

/*
 * GtkScrolledWindow
 */

// ScrolledWindow is a representation of GTK's GtkScrolledWindow.
type ScrolledWindow struct {
	Bin
}

var scrolledWindowType = glib.Type(C.gtk_scrolled_window_get_type())

func GetScrolledWindowType() glib.Type {
	return glib.Type(scrolledWindowType)
}

// Native() returns a pointer to the underlying GtkScrolledWindow.
func (v *ScrolledWindow) Native() *C.GtkScrolledWindow {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(scrolledWindowType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkScrolledWindow)(v.Ptr())
}

func wrapScrolledWindow(obj *glib.Object) (s ScrolledWindow) {
	s.Bin = wrapBin(obj)
	return
}

// ScrolledWindowNew() is a wrapper around gtk_scrolled_window_new().
func ScrolledWindowNew(hadjustment, vadjustment *Adjustment) (*ScrolledWindow, error) {
	c := C.gtk_scrolled_window_new(hadjustment.Native(),
		vadjustment.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	s := wrapScrolledWindow(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &s, nil
}

// SetPolicy() is a wrapper around gtk_scrolled_window_set_policy().
func (v *ScrolledWindow) SetPolicy(hScrollbarPolicy, vScrollbarPolicy PolicyType) {
	C.gtk_scrolled_window_set_policy(v.Native(),
		C.GtkPolicyType(hScrollbarPolicy),
		C.GtkPolicyType(vScrollbarPolicy))
}

/*
 * GtkSpinButton
 */

// SpinButton is a representation of GTK's GtkSpinButton.
type SpinButton struct {
	Entry
}

var spinButtonType = glib.Type(C.gtk_spin_button_get_type())

func GetSpinButtonType() glib.Type {
	return glib.Type(spinButtonType)
}

// Native() returns a pointer to the underlying GtkSpinButton.
func (v *SpinButton) Native() *C.GtkSpinButton {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(spinButtonType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkSpinButton)(v.Ptr())
}

func wrapSpinButton(obj *glib.Object) (s SpinButton) {
	s.Entry = wrapEntry(obj)
	return
}

// Configure() is a wrapper around gtk_spin_button_configure().
func (v *SpinButton) Configure(adjustment *Adjustment, climbRate float64, digits uint) {
	C.gtk_spin_button_configure(v.Native(), adjustment.Native(),
		C.gdouble(climbRate), C.guint(digits))
}

// SpinButtonNew() is a wrapper around gtk_spin_button_new().
func SpinButtonNew(adjustment *Adjustment, climbRate float64, digits uint) (*SpinButton, error) {
	c := C.gtk_spin_button_new(adjustment.Native(),
		C.gdouble(climbRate), C.guint(digits))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	s := wrapSpinButton(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &s, nil
}

// SpinButtonNewWithRange() is a wrapper around
// gtk_spin_button_new_with_range().
func SpinButtonNewWithRange(min, max, step float64) (*SpinButton, error) {
	c := C.gtk_spin_button_new_with_range(C.gdouble(min), C.gdouble(max),
		C.gdouble(step))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	s := wrapSpinButton(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &s, nil
}

// ValueAsInt() is a wrapper around gtk_spin_button_get_value_as_int().
func (v *SpinButton) ValueAsInt() int {
	c := C.gtk_spin_button_get_value_as_int(v.Native())
	return int(c)
}

// SetValue() is a wrapper around gtk_spin_button_set_value().
func (v *SpinButton) SetValue(value float64) {
	C.gtk_spin_button_set_value(v.Native(), C.gdouble(value))
}

// Value() is a wrapper around gtk_spin_button_get_value().
func (v *SpinButton) Value() float64 {
	c := C.gtk_spin_button_get_value(v.Native())
	return float64(c)
}

/*
 * GtkStatusbar
 */

// Statusbar is a representation of GTK's GtkStatusbar
type Statusbar struct {
	Box
}

var statusbarType = glib.Type(C.gtk_statusbar_get_type())

func GetStatusbarType() glib.Type {
	return glib.Type(statusbarType)
}

// Native() returns a pointer to the underlying GtkStatusbar
func (v *Statusbar) Native() *C.GtkStatusbar {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(statusbarType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkStatusbar)(v.Ptr())
}

func wrapStatusbar(obj *glib.Object) (s Statusbar) {
	s.Box = wrapBox(obj)
	return
}

// StatusbarNew() is a wrapper around gtk_statusbar_new().
func StatusbarNew() (*Statusbar, error) {
	c := C.gtk_statusbar_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	s := wrapStatusbar(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &s, nil
}

// ContextId() is a wrapper around gtk_statusbar_get_context_id().
func (v *Statusbar) GetContextId(contextDescription string) uint {
	cstr := C.CString(contextDescription)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_statusbar_get_context_id(v.Native(), (*C.gchar)(cstr))
	return uint(c)
}

// Push() is a wrapper around gtk_statusbar_push().
func (v *Statusbar) Push(contextID uint, text string) uint {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_statusbar_push(v.Native(), C.guint(contextID),
		(*C.gchar)(cstr))
	return uint(c)
}

// Pop() is a wrapper around gtk_statusbar_pop().
func (v *Statusbar) Pop(contextID uint) {
	C.gtk_statusbar_pop(v.Native(), C.guint(contextID))
}

// MessageArea() is a wrapper around gtk_statusbar_get_message_area().
func (v *Statusbar) MessageArea() (*Box, error) {
	c := C.gtk_statusbar_get_message_area(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	b := wrapBox(obj)
	return &b, nil
}

/*
 * GtkTextBuffer
 */
// TextBuffer is a representation of GTK's GtkTextBuffer.
type TextBuffer struct {
	*glib.Object
}

var textBufferType = glib.Type(C.gtk_text_buffer_get_type())

func GetTextBufferType() glib.Type {
	return glib.Type(textBufferType)
}

func (v *TextBuffer) Native() *C.GtkTextBuffer {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(textBufferType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkTextBuffer)(v.Ptr())
}

func (t *TextBuffer) Insert(iter *TextIter, text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_text_buffer_insert(t.Native(), iter.c, (*C.gchar)(cstr), C.gint(len(text)))
}

func (t *TextBuffer) GetIterAtOffset(offset int) *TextIter {
	var c C.GtkTextIter
	C.gtk_text_buffer_get_iter_at_offset(t.Native(), &c, C.gint(offset))
	iter := &TextIter{&c}
	//runtime.SetFinalizer(iter, freeTextIter)
	return iter
}

func (t *TextBuffer) SetText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_text_buffer_set_text(t.Native(), (*C.gchar)(cstr), C.gint(len(text)))
}

/*
 * GtkTextIter
 */
// TextIter is a representation of GTK's GtkTextIter.
type TextIter struct {
	c *C.GtkTextIter
}

var textIterType = glib.Type(C.gtk_text_iter_get_type())

func GetTextIterType() glib.Type {
	return glib.Type(textIterType)
}

/*
 * GtkTextView
 */
// TextView is a representation of GTK's GtkTextView.
type TextView struct {
	Container
}

var textViewType = glib.Type(C.gtk_text_view_get_type())

func GetTextViewType() glib.Type {
	return glib.Type(textViewType)
}

func (v *TextView) Native() *C.GtkTextView {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(textViewType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkTextView)(v.Ptr())
}

func wrapTextView(obj *glib.Object) (t TextView) {
	t.Container = wrapContainer(obj)
	return
}

func (t *TextView) Buffer() (*TextBuffer, error) {
	c := C.gtk_text_view_get_buffer(t.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	buffer := &TextBuffer{obj}
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return buffer, nil
}

/*
 * GtkTreeIter
 */

// TreeIter is a representation of GTK's GtkTreeIter.
type TreeIter struct {
	GtkTreeIter C.GtkTreeIter
}

var treeIterType = glib.Type(C.gtk_tree_iter_get_type())

func GetTreeIterType() glib.Type {
	return glib.Type(treeIterType)
}

// Native() returns a pointer to the underlying GtkTreeIter.
func (v *TreeIter) Native() *C.GtkTreeIter {
	if v == nil {
		return nil
	}
	return &v.GtkTreeIter
}

func (v *TreeIter) free() {
	C.gtk_tree_iter_free(v.Native())
}

// Copy() is a wrapper around gtk_tree_iter_copy().
func (v *TreeIter) Copy() (*TreeIter, error) {
	c := C.gtk_tree_iter_copy(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	t := &TreeIter{*c}
	runtime.SetFinalizer(t, (*TreeIter).free)
	return t, nil
}

/*
 * GtkTreeModel
 */

// TreeModel is a representation of GTK's GtkTreeModel GInterface.
type TreeModel struct {
	*glib.Object

	// other useful data
	indexMap map[string]int
}

var treeModelType = glib.Type(C.gtk_tree_model_get_type())

func GetTreeModelType() glib.Type {
	return glib.Type(treeModelType)
}

// ITreeModel is an interface type implemented by all structs
// embedding a TreeModel.  It is meant to be used as an argument type
// for wrapper functions that wrap around a C GTK function taking a
// GtkTreeModel.
type ITreeModel interface {
	toTreeModel() *C.GtkTreeModel
	columnIndex(string) (int, bool)
}

// Native() returns a pointer to the underlying GObject as a GtkTreeModel.
func (v *TreeModel) Native() *C.GtkTreeModel {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(treeModelType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkTreeModel)(v.Ptr())
}

func (v *TreeModel) toTreeModel() *C.GtkTreeModel {
	if v == nil {
		return nil
	}
	return v.Native()
}

func (v *TreeModel) columnIndex(name string) (int, bool) {
	i, ok := v.indexMap[name]
	return i, ok
}

func wrapTreeModel(obj *glib.Object) (t TreeModel) {
	t.Object = obj
	return
}

// Flags() is a wrapper around gtk_tree_model_get_flags().
func (v *TreeModel) Flags() TreeModelFlags {
	c := C.gtk_tree_model_get_flags(v.Native())
	return TreeModelFlags(c)
}

// NColumns() is a wrapper around gtk_tree_model_get_n_columns().
func (v *TreeModel) NColumns() int {
	c := C.gtk_tree_model_get_n_columns(v.Native())
	return int(c)
}

// ColumnType() is a wrapper around gtk_tree_model_get_column_type().
func (v *TreeModel) GetColumnType(index int) glib.Type {
	c := C.gtk_tree_model_get_column_type(v.Native(), C.gint(index))
	return glib.Type(c)
}

// Iter() is a wrapper around gtk_tree_model_get_iter().
func (v *TreeModel) GetIter(path *TreePath) (*TreeIter, error) {
	var iter C.GtkTreeIter
	c := C.gtk_tree_model_get_iter(v.Native(), &iter, path.Native())
	if !gobool(c) {
		return nil, errors.New("Unable to set iterator")
	}
	t := &TreeIter{iter}
	runtime.SetFinalizer(t, (*TreeIter).free)
	return t, nil
}

// IterFromString() is a wrapper around
// gtk_tree_model_get_iter_from_string().
func (v *TreeModel) GetIterFromString(path string) (*TreeIter, error) {
	var iter C.GtkTreeIter
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_tree_model_get_iter_from_string(v.Native(), &iter,
		(*C.gchar)(cstr))
	if !gobool(c) {
		return nil, errors.New("Unable to set iterator")
	}
	t := &TreeIter{iter}
	runtime.SetFinalizer(t, (*TreeIter).free)
	return t, nil
}

// IterFirst() is a wrapper around gtk_tree_model_get_iter_first().
func (v *TreeModel) IterFirst() (*TreeIter, error) {
	var iter C.GtkTreeIter
	c := C.gtk_tree_model_get_iter_first(v.Native(), &iter)
	if !gobool(c) {
		return nil, errors.New("Unable to set iterator")
	}
	t := &TreeIter{iter}
	runtime.SetFinalizer(t, (*TreeIter).free)
	return t, nil
}

// Path() is a wrapper around gtk_tree_model_get_path().
func (v *TreeModel) GetPath(iter *TreeIter) (*TreePath, error) {
	c := C.gtk_tree_model_get_path(v.Native(), iter.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	p := &TreePath{c}
	runtime.SetFinalizer(p, (*TreePath).free)
	return p, nil
}

// Value() is a wrapper around gtk_tree_model_get_value().
func (v *TreeModel) GetValue(iter *TreeIter, column int) (*glib.Value, error) {
	val, err := glib.ValueAlloc()
	if err != nil {
		return nil, err
	}
	C.gtk_tree_model_get_value(
		(*C.GtkTreeModel)(unsafe.Pointer(v.Native())),
		iter.Native(),
		C.gint(column),
		(*C.GValue)(unsafe.Pointer(val.Native())))
	return val, nil
}

/*
 * GtkTreePath
 */

// TreePath is a representation of GTK's GtkTreePath.
type TreePath struct {
	GtkTreePath *C.GtkTreePath
}

var treePathType = glib.Type(C.gtk_tree_path_get_type())

func GetTreePathType() glib.Type {
	return glib.Type(treePathType)
}

// Native() returns a pointer to the underlying GtkTreePath.
func (v *TreePath) Native() *C.GtkTreePath {
	if v == nil {
		return nil
	}
	return v.GtkTreePath
}

func (v *TreePath) free() {
	C.gtk_tree_path_free(v.Native())
}

/*
 * GtkTreeSelection
 */

// TreeSelection is a representation of GTK's GtkTreeSelection.
type TreeSelection struct {
	*glib.Object
}

var treeSelectionType = glib.Type(C.gtk_tree_selection_get_type())

func GetTreeSelectionType() glib.Type {
	return glib.Type(treeSelectionType)
}

// Native() returns a pointer to the underlying GtkTreeSelection.
func (v *TreeSelection) Native() *C.GtkTreeSelection {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(treeSelectionType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkTreeSelection)(v.Ptr())
}

func wrapTreeSelection(obj *glib.Object) TreeSelection {
	return TreeSelection{obj}
}

// Selected() is a wrapper around gtk_tree_selection_get_selected().
func (v *TreeSelection) GetSelected(model *ITreeModel, iter *TreeIter) bool {
	var pcmodel **C.GtkTreeModel
	if pcmodel != nil {
		cmodel := (*model).toTreeModel()
		pcmodel = &cmodel
	} else {
		pcmodel = nil
	}
	c := C.gtk_tree_selection_get_selected(v.Native(),
		pcmodel, iter.Native())
	return gobool(c)
}

/*
 * GtkTreeView
 */

// TreeView is a representation of GTK's GtkTreeView.
type TreeView struct {
	Container
}

var treeViewType = glib.Type(C.gtk_tree_view_get_type())

func GetTreeViewType() glib.Type {
	return glib.Type(treeViewType)
}

// Native() returns a pointer to the underlying GtkTreeView.
func (v *TreeView) Native() *C.GtkTreeView {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(treeViewType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkTreeView)(v.Ptr())
}

func wrapTreeView(obj *glib.Object) (t TreeView) {
	t.Container = wrapContainer(obj)
	return
}

// TreeViewNew() is a wrapper around gtk_tree_view_new().
func TreeViewNew() (*TreeView, error) {
	c := C.gtk_tree_view_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	t := wrapTreeView(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &t, nil
}

// TreeViewNewWithModel() is a wrapper around gtk_tree_view_new_with_model().
func TreeViewNewWithModel(model ITreeModel) (*TreeView, error) {
	c := C.gtk_tree_view_new_with_model(model.toTreeModel())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	t := wrapTreeView(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &t, nil
}

// Model() is a wrapper around gtk_tree_view_get_model().
func (v *TreeView) Model() (*TreeModel, error) {
	c := C.gtk_tree_view_get_model(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	t := wrapTreeModel(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &t, nil
}

// SetModel() is a wrapper around gtk_tree_view_set_model().
func (v *TreeView) SetModel(model ITreeModel) {
	C.gtk_tree_view_set_model(v.Native(), model.toTreeModel())
}

// Selection() is a wrapper around gtk_tree_view_get_selection().
func (v *TreeView) Selection() (*TreeSelection, error) {
	c := C.gtk_tree_view_get_selection(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	s := wrapTreeSelection(obj)
	obj.Ref()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &s, nil
}

// AppendColumn() is a wrapper around gtk_tree_view_append_column().
func (v *TreeView) AppendColumn(column *TreeViewColumn) int {
	c := C.gtk_tree_view_append_column(v.Native(), column.Native())
	return int(c)
}

/*
 * GtkTreeViewColumn
 */

// TreeViewColumns is a representation of GTK's GtkTreeViewColumn.
type TreeViewColumn struct {
	glib.InitiallyUnowned
}

var treeViewColumnType = glib.Type(C.gtk_tree_view_column_get_type())

func GetTreeViewColumnType() glib.Type {
	return glib.Type(treeViewColumnType)
}

// Native() returns a pointer to the underlying GtkTreeViewColumn.
func (v *TreeViewColumn) Native() *C.GtkTreeViewColumn {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(treeViewColumnType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkTreeViewColumn)(v.Ptr())
}

func wrapTreeViewColumn(obj *glib.Object) TreeViewColumn {
	return TreeViewColumn{glib.InitiallyUnowned{*obj}}
}

// TreeViewColumnNew() is a wrapper around gtk_tree_view_column_new().
func TreeViewColumnNew() (*TreeViewColumn, error) {
	c := C.gtk_tree_view_column_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	t := wrapTreeViewColumn(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &t, nil
}

// TreeViewColumnNewWithAttribute() creates a new column with the specified attributes.
// Rather than wrap around any C functions, it uses existing gotk3 calls to create the column,
// set its title, pack the renderer, and set attributes in separate steps.
// The model is needed to be able to map attribute values to the correct column indexes.
func TreeViewColumnNewWithAttributes(title string, renderer ICellRenderer, model ITreeModel, attributes map[string]string) (*TreeViewColumn, error) {
	v, err := TreeViewColumnNew()
	if err != nil {
		return nil, err
	}
	v.SetTitle(title)
	v.PackStart(renderer, true)
	for key, value := range attributes {
		k_cstr := C.CString(key)
		col, ok := model.columnIndex(value)
		if !ok {
			return nil, fmt.Errorf("unknown column '%s'", value)
		}
		C._gtk_tree_view_column_set_attribute(v.Native(), renderer.toCellRenderer(), (*C.gchar)(k_cstr), C.gint(col))
		C.free(unsafe.Pointer(k_cstr))
	}
	return v, nil
}

// AddAttribute() is a wrapper around gtk_tree_view_column_add_attribute().
func (v *TreeViewColumn) AddAttribute(renderer ICellRenderer, attribute string, column int) {
	cstr := C.CString(attribute)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_tree_view_column_add_attribute(v.Native(),
		renderer.toCellRenderer(), (*C.gchar)(cstr), C.gint(column))
}

// SetExpand() is a wrapper around gtk_tree_view_column_set_expand().
func (v *TreeViewColumn) SetExpand(expand bool) {
	C.gtk_tree_view_column_set_expand(v.Native(), gbool(expand))
}

// Expand() is a wrapper around gtk_tree_view_column_get_expand().
func (v *TreeViewColumn) Expand() bool {
	c := C.gtk_tree_view_column_get_expand(v.Native())
	return gobool(c)
}

// SetResizable() is a wrapper around gtk_tree_view_column_set_resizable().
func (v *TreeViewColumn) SetResizable(resizable bool) {
	C.gtk_tree_view_column_set_resizable(v.Native(), gbool(resizable))
}

// Resizable() is a wrapper around gtk_tree_view_column_get_resizable().
func (v *TreeViewColumn) Resizable() bool {
	c := C.gtk_tree_view_column_get_resizable(v.Native())
	return gobool(c)
}

// SetMinWidth() is a wrapper around gtk_tree_view_column_set_min_width().
func (v *TreeViewColumn) SetMinWidth(minWidth int) {
	C.gtk_tree_view_column_set_min_width(v.Native(), C.gint(minWidth))
}

// MinWidth() is a wrapper around gtk_tree_view_column_get_min_width().
func (v *TreeViewColumn) MinWidth() int {
	c := C.gtk_tree_view_column_get_min_width(v.Native())
	return int(c)
}

// SetTitle() is a wrapper around gtk_tree_view_column_set_title().
func (v *TreeViewColumn) SetTitle(title string) {
	t_cstr := C.CString(title)
	defer C.free(unsafe.Pointer(t_cstr))
	C.gtk_tree_view_column_set_title(v.Native(), (*C.gchar)(t_cstr))
}

// PackStart() is a wrapper around gtk_tree_view_column_pack_start().
func (v *TreeViewColumn) PackStart(renderer ICellRenderer, expand bool) {
	C.gtk_tree_view_column_pack_start(v.Native(), renderer.toCellRenderer(), gbool(expand))
}

/*
 * GtkWidget
 */

// Widget is a representation of GTK's GtkWidget.
type Widget struct {
	glib.InitiallyUnowned

	// Interfaces
	Buildable
}

var widgetType = glib.Type(C.gtk_widget_get_type())

func GetWidgetType() glib.Type {
	return glib.Type(widgetType)
}

// IWidget is an interface type implemented by all structs
// embedding a Widget.  It is meant to be used as an argument type
// for wrapper functions that wrap around a C GTK function taking a
// GtkWidget.
type IWidget interface {
	toWidget() *C.GtkWidget
}

// Native() returns a pointer to the underlying GtkWidget.
func (v *Widget) Native() *C.GtkWidget {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(widgetType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkWidget)(v.Ptr())
}

func (v *Widget) toWidget() *C.GtkWidget {
	if v == nil {
		return nil
	}
	return v.Native()
}

func wrapWidget(obj *glib.Object) (w Widget) {
	w.InitiallyUnowned = glib.InitiallyUnowned{*obj}
	w.Buildable = Buildable{obj.Ptr()}
	return
}

// Destroy() is a wrapper around gtk_widget_destroy().
func (v *Widget) Destroy() {
	C.gtk_widget_destroy(v.Native())
}

// InDestruction() is a wrapper around gtk_widget_in_destruction().
func (v *Widget) InDestruction() bool {
	return gobool(C.gtk_widget_in_destruction(v.Native()))
}

// TODO(jrick) this may require some rethinking
/*
func (v *Widget) Destroyed(widgetPointer **Widget) {
}
*/

// Unparent() is a wrapper around gtk_widget_unparent().
func (v *Widget) Unparent() {
	C.gtk_widget_unparent(v.Native())
}

// Show() is a wrapper around gtk_widget_show().
func (v *Widget) Show() {
	C.gtk_widget_show(v.Native())
}

// Hide() is a wrapper around gtk_widget_hide().
func (v *Widget) Hide() {
	C.gtk_widget_hide(v.Native())
}

// SetVisible() is a wrapper around gtk_widget_set_visible().
func (v *Widget) SetVisible(visible bool) {
	C.gtk_widget_set_visible(v.Native(), gbool(visible))
}

// ShowNow() is a wrapper around gtk_widget_show_now().
func (v *Widget) ShowNow() {
	C.gtk_widget_show_now(v.Native())
}

// ShowAll() is a wrapper around gtk_widget_show_all().
func (v *Widget) ShowAll() {
	C.gtk_widget_show_all(v.Native())
}

// SetNoShowAll() is a wrapper around gtk_widget_set_no_show_all().
func (v *Widget) SetNoShowAll(noShowAll bool) {
	C.gtk_widget_set_no_show_all(v.Native(), gbool(noShowAll))
}

// NoShowAll() is a wrapper around gtk_widget_get_no_show_all().
func (v *Widget) NoShowAll() bool {
	c := C.gtk_widget_get_no_show_all(v.Native())
	return gobool(c)
}

func (v *Widget) HideOnDelete() bool {
	return gobool(C.gtk_widget_hide_on_delete(v.Native()))
}

// Map() is a wrapper around gtk_widget_map().
func (v *Widget) Map() {
	C.gtk_widget_map(v.Native())
}

// Unmap() is a wrapper around gtk_widget_unmap().
func (v *Widget) Unmap() {
	C.gtk_widget_unmap(v.Native())
}

//void gtk_widget_realize(GtkWidget *widget);
//void gtk_widget_unrealize(GtkWidget *widget);
//void gtk_widget_draw(GtkWidget *widget, cairo_t *cr);
//void gtk_widget_queue_resize(GtkWidget *widget);
//void gtk_widget_queue_resize_no_redraw(GtkWidget *widget);
//GdkFrameClock *gtk_widget_get_frame_clock(GtkWidget *widget);
//guint gtk_widget_add_tick_callback (GtkWidget *widget,
//                                    GtkTickCallback callback,
//                                    gpointer user_data,
//                                    GDestroyNotify notify);
//void gtk_widget_remove_tick_callback(GtkWidget *widget, guint id);

// TODO(jrick) GtkAllocation
/*
func (v *Widget) SizeAllocate() {
}
*/

// TODO(jrick) GtkAccelGroup GdkModifierType GtkAccelFlags
/*
func (v *Widget) AddAccelerator() {
}
*/

// TODO(jrick) GtkAccelGroup GdkModifierType
/*
func (v *Widget) RemoveAccelerator() {
}
*/

// TODO(jrick) GtkAccelGroup
/*
func (v *Widget) SetAccelPath() {
}
*/

// TODO(jrick) GList
/*
func (v *Widget) ListAccelClosures() {
}
*/

//gboolean gtk_widget_can_activate_accel(GtkWidget *widget, guint signal_id);

// Event() is a wrapper around gtk_widget_event().
func (v *Widget) Event(event *gdk.Event) bool {
	c := C.gtk_widget_event(v.Native(),
		(*C.GdkEvent)(unsafe.Pointer(event.Native())))
	return gobool(c)
}

// Activate() is a wrapper around gtk_widget_activate().
func (v *Widget) Activate() bool {
	return gobool(C.gtk_widget_activate(v.Native()))
}

// Reparent() is a wrapper around gtk_widget_reparent().
func (v *Widget) Reparent(newParent IWidget) {
	C.gtk_widget_reparent(v.Native(), newParent.toWidget())
}

// TODO(jrick) GdkRectangle
/*
func (v *Widget) Intersect() {
}
*/

// IsFocus() is a wrapper around gtk_widget_is_focus().
func (v *Widget) IsFocus() bool {
	return gobool(C.gtk_widget_is_focus(v.Native()))
}

// GrabFocus() is a wrapper around gtk_widget_grab_focus().
func (v *Widget) GrabFocus() {
	C.gtk_widget_grab_focus(v.Native())
}

// GrabDefault() is a wrapper around gtk_widget_grab_default().
func (v *Widget) GrabDefault() {
	C.gtk_widget_grab_default(v.Native())
}

// SetName() is a wrapper around gtk_widget_set_name().
func (v *Widget) SetName(name string) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_widget_set_name(v.Native(), (*C.gchar)(cstr))
}

// Name() is a wrapper around gtk_widget_get_name().  A non-nil
// error is returned in the case that gtk_widget_get_name returns NULL to
// differentiate between NULL and an empty string.
func (v *Widget) Name() (string, error) {
	c := C.gtk_widget_get_name(v.Native())
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// SetSensitive() is a wrapper around gtk_widget_set_sensitive().
func (v *Widget) SetSensitive(sensitive bool) {
	C.gtk_widget_set_sensitive(v.Native(), gbool(sensitive))
}

// SetParent() is a wrapper around gtk_widget_set_parent().
func (v *Widget) SetParent(parent IWidget) {
	C.gtk_widget_set_parent(v.Native(), parent.toWidget())
}

// Parent() is a wrapper around gtk_widget_get_parent().
func (v *Widget) Parent() (*Widget, error) {
	c := C.gtk_widget_get_parent(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// SetSizeRequest() is a wrapper around gtk_widget_set_size_request().
func (v *Widget) SetSizeRequest(width, height int) {
	C.gtk_widget_set_size_request(v.Native(), C.gint(width), C.gint(height))
}

// SizeRequest() is a wrapper around gtk_widget_get_size_request().
func (v *Widget) SizeRequest() (width, height int) {
	var w, h C.gint
	C.gtk_widget_get_size_request(v.Native(), &w, &h)
	return int(w), int(h)
}

// SetParentWindow() is a wrapper around gtk_widget_set_parent_window().
func (v *Widget) SetParentWindow(parentWindow *gdk.Window) {
	C.gtk_widget_set_parent_window(v.Native(),
		(*C.GdkWindow)(unsafe.Pointer(parentWindow.Native())))
}

// ParentWindow() is a wrapper around gtk_widget_get_parent_window().
func (v *Widget) ParentWindow() (*gdk.Window, error) {
	c := C.gtk_widget_get_parent_window(v.Native())
	if v == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := gdk.Window{obj}
	w.Ref()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// SetEvents() is a wrapper around gtk_widget_set_events().
func (v *Widget) SetEvents(events int) {
	C.gtk_widget_set_events(v.Native(), C.gint(events))
}

// Events() is a wrapper around gtk_widget_get_events().
func (v *Widget) Events() int {
	return int(C.gtk_widget_get_events(v.Native()))
}

// AddEvents() is a wrapper around gtk_widget_add_events().
func (v *Widget) AddEvents(events int) {
	C.gtk_widget_add_events(v.Native(), C.gint(events))
}

// TODO(jrick) GdkEventMask
/*
func (v *Widget) SetDeviceEvents() {
}
*/

// TODO(jrick) GdkEventMask
/*
func (v *Widget) DeviceEvents() {
}
*/

// TODO(jrick) GdkEventMask
/*
func (v *Widget) AddDeviceEvents() {
}
*/

// SetDeviceEnabled() is a wrapper around gtk_widget_set_device_enabled().
func (v *Widget) SetDeviceEnabled(device *gdk.Device, enabled bool) {
	C.gtk_widget_set_device_enabled(v.Native(),
		(*C.GdkDevice)(unsafe.Pointer(device.Native())), gbool(enabled))
}

// DeviceEnabled() is a wrapper around gtk_widget_get_device_enabled().
func (v *Widget) GetDeviceEnabled(device *gdk.Device) bool {
	c := C.gtk_widget_get_device_enabled(v.Native(),
		(*C.GdkDevice)(unsafe.Pointer(device.Native())))
	return gobool(c)
}

// Toplevel() is a wrapper around gtk_widget_get_toplevel().
func (v *Widget) Toplevel() (*Widget, error) {
	c := C.gtk_widget_get_toplevel(v.Native())
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWidget(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// TooltipText() is a wrapper around gtk_widget_get_tooltip_text().
// A non-nil error is returned in the case that
// gtk_widget_get_tooltip_text returns NULL to differentiate between NULL
// and an empty string.
func (v *Widget) TooltipText() (string, error) {
	c := C.gtk_widget_get_tooltip_text(v.Native())
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}

// SetTooltipText() is a wrapper around gtk_widget_set_tooltip_text().
func (v *Widget) SetTooltipText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_widget_set_tooltip_text(v.Native(), (*C.gchar)(cstr))
}

// OverrideFont() is a wrapper around gtk_widget_override_font().
func (v *Widget) OverrideFont(description string) {
	cstr := C.CString(description)
	defer C.free(unsafe.Pointer(cstr))
	c := C.pango_font_description_from_string(cstr)
	C.gtk_widget_override_font(v.Native(), c)
}

// HAlign() is a wrapper around gtk_widget_get_halign().
func (v *Widget) HAlign() Align {
	c := C.gtk_widget_get_halign(v.Native())
	return Align(c)
}

// SetHAlign() is a wrapper around gtk_widget_set_halign().
func (v *Widget) SetHAlign(align Align) {
	C.gtk_widget_set_halign(v.Native(), C.GtkAlign(align))
}

// VAlign() is a wrapper around gtk_widget_get_valign().
func (v *Widget) VAlign() Align {
	c := C.gtk_widget_get_valign(v.Native())
	return Align(c)
}

// SetVAlign() is a wrapper around gtk_widget_set_valign().
func (v *Widget) SetVAlign(align Align) {
	C.gtk_widget_set_valign(v.Native(), C.GtkAlign(align))
}

// HExpand() is a wrapper around gtk_widget_get_hexpand().
func (v *Widget) HExpand() bool {
	c := C.gtk_widget_get_hexpand(v.Native())
	return gobool(c)
}

// SetHExpand() is a wrapper around gtk_widget_set_hexpand().
func (v *Widget) SetHExpand(expand bool) {
	C.gtk_widget_set_hexpand(v.Native(), gbool(expand))
}

// VExpand() is a wrapper around gtk_widget_get_vexpand().
func (v *Widget) VExpand() bool {
	c := C.gtk_widget_get_vexpand(v.Native())
	return gobool(c)
}

// SetVExpand() is a wrapper around gtk_widget_set_vexpand().
func (v *Widget) SetVExpand(expand bool) {
	C.gtk_widget_set_vexpand(v.Native(), gbool(expand))
}

/*
 * GtkWindow
 */

// Window is a representation of GTK's GtkWindow.
type Window struct {
	Bin
}

var windowType = glib.Type(C.gtk_window_get_type())

func GetWindowType() glib.Type {
	return glib.Type(windowType)
}

// IWindow is an interface type implemented by all structs embedding a
// Window.  It is meant to be used as an argument type for wrapper
// functions that wrap around a C GTK function taking a GtkWindow.
type IWindow interface {
	toWindow() *C.GtkWindow
}

// Native() returns a pointer to the underlying GtkWindow.
func (v *Window) Native() *C.GtkWindow {
	if v == nil {
		return nil
	}
	if warn := v.Typecheck(windowType); warn != nil {
		fmt.Fprintln(os.Stderr, warn)
	}
	return (*C.GtkWindow)(v.Ptr())
}

func (v *Window) toWindow() *C.GtkWindow {
	if v == nil {
		return nil
	}
	return v.Native()
}

func wrapWindow(obj *glib.Object) (w Window) {
	w.Bin = wrapBin(obj)
	return
}

// WindowNew() is a wrapper around gtk_window_new().
func WindowNew(t WindowType) (*Window, error) {
	c := C.gtk_window_new(C.GtkWindowType(t))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.ObjectNew(unsafe.Pointer(c))
	w := wrapWindow(obj)
	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	return &w, nil
}

// SetTitle() is a wrapper around gtk_window_set_title().
func (v *Window) SetTitle(title string) {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_window_set_title(v.Native(), (*C.gchar)(cstr))
}

// SetDefaultSize() is a wrapper around gtk_window_set_default_size().
func (v *Window) SetDefaultSize(width, height int) {
	C.gtk_window_set_default_size(v.Native(), C.gint(width), C.gint(height))
}

// SetDefaultGeometry() is a wrapper around gtk_window_set_default_geometry().
func (v *Window) SetDefaultGeometry(width, height int) {
	C.gtk_window_set_default_geometry(v.Native(), C.gint(width),
		C.gint(height))
}

// TODO(jrick) GdkGeometry GdkWindowHints
/*
func (v *Window) SetGeometryHints() {
}
*/

// TODO(jrick) GdkGravity
/*
func (v *Window) SetGravity() {
}
*/

// TODO(jrick) GdkGravity
/*
func (v *Window) Gravity() {
}
*/

// SetPosition() is a wrapper around gtk_window_set_position()
func (v *Window) SetPosition(position WindowPosition) {
	C.gtk_window_set_position(v.Native(), C.GtkWindowPosition(position))
}

// SetTransientFor() is a wrapper around gtk_window_set_transient_for().
func (v *Window) SetTransientFor(parent IWindow) {
	var pw *C.GtkWindow = nil
	if parent != nil {
		pw = parent.toWindow()
	}
	C.gtk_window_set_transient_for(v.Native(), pw)
}

// cast() takes a native GObject and casts it to the appropriate Go struct.
func cast(c *C.GObject) (glib.IObject, error) {
	var (
		className = C.GoString((*C.char)(C.object_get_class_name(c)))
		obj       = glib.ObjectNew(unsafe.Pointer(c))
	)
	runtime.SetFinalizer(obj, (*glib.Object).Unref)
	switch className {
	case "GtkAdjustment":
		a := wrapAdjustment(obj)
		return &a, nil
	case "GtkBin":
		b := wrapBin(obj)
		return &b, nil
	case "GtkBox":
		b := wrapBox(obj)
		return &b, nil
	case "GtkButton":
		b := wrapButton(obj)
		return &b, nil
	case "GtkCellRenderer":
		c := wrapCellRenderer(obj)
		return &c, nil
	case "GtkCellRendererText":
		c := wrapCellRendererText(obj)
		return &c, nil
	case "GtkClipboard":
		c := wrapClipboard(obj)
		return &c, nil
	case "GtkComboBox":
		c := wrapComboBox(obj)
		return &c, nil
	case "GtkContainer":
		c := wrapContainer(obj)
		return &c, nil
	case "GtkDialog":
		d := wrapDialog(obj)
		return &d, nil
	case "GtkEntry":
		e := wrapEntry(obj)
		return &e, nil
	case "GtkEntryBuffer":
		e := wrapEntryBuffer(obj)
		return &e, nil
	case "GtkEntryCompletion":
		e := wrapEntryCompletion(obj)
		return &e, nil
	case "GtkFileChooserButton":
		f := wrapFileChooserButton(obj)
		return &f, nil
	case "GtkGrid":
		g := wrapGrid(obj)
		return &g, nil
	case "GtkImage":
		i := wrapImage(obj)
		return &i, nil
	case "GtkImageMenuItem":
		i := wrapImageMenuItem(obj)
		return &i, nil
	case "GtkLabel":
		l := wrapLabel(obj)
		return &l, nil
	case "GtkListStore":
		l := wrapListStore(obj)
		return &l, nil
	case "GtkMenu":
		m := wrapMenu(obj)
		return &m, nil
	case "GtkMenuBar":
		m := wrapMenuBar(obj)
		return &m, nil
	case "GtkMenuItem":
		m := wrapMenuItem(obj)
		return &m, nil
	case "GtkMenuShell":
		m := wrapMenuShell(obj)
		return &m, nil
	case "GtkMessageDialog":
		m := wrapMessageDialog(obj)
		return &m, nil
	case "GtkMisc":
		m := wrapMisc(obj)
		return &m, nil
	case "GtkNotebook":
		n := wrapNotebook(obj)
		return &n, nil
	case "GtkOffscreenWindow":
		o := wrapOffscreenWindow(obj)
		return &o, nil
	case "GtkProgressBar":
		p := wrapProgressBar(obj)
		return &p, nil
	case "GtkScrolledWindow":
		s := wrapScrolledWindow(obj)
		return &s, nil
	case "GtkSpinButton":
		s := wrapSpinButton(obj)
		return &s, nil
	case "GtkStatusbar":
		s := wrapStatusbar(obj)
		return &s, nil
	case "GtkTextView":
		t := wrapTextView(obj)
		return &t, nil
	case "GtkTreeModel":
		t := wrapTreeModel(obj)
		return &t, nil
	case "GtkTreeSelection":
		t := wrapTreeSelection(obj)
		return &t, nil
	case "GtkTreeView":
		t := wrapTreeView(obj)
		return &t, nil
	case "GtkTreeViewColumn":
		t := wrapTreeViewColumn(obj)
		return &t, nil
	case "GtkWidget":
		w := wrapWidget(obj)
		return &w, nil
	case "GtkWindow":
		w := wrapWindow(obj)
		return &w, nil
	default:
		return nil, errors.New("unrecognized class name '" + className + "'")
	}
}
