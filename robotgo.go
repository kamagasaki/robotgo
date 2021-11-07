// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

/*

Package robotgo Go native cross-platform system automation.

Please make sure Golang, GCC is installed correctly before installing RobotGo;

See Requirements:
	https://github.com/go-vgo/robotgo#requirements

Installation:
	go get -u github.com/go-vgo/robotgo
*/
package robotgo

/*
//#if defined(IS_MACOSX)
	#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations
	#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit
	#cgo darwin LDFLAGS: -framework Carbon -framework CoreFoundation
//#elif defined(USE_X11)
	#cgo linux CFLAGS: -I/usr/src
	#cgo linux LDFLAGS: -L/usr/src -lX11 -lXtst -lm
	// #cgo linux LDFLAGS: -lX11-xcb -lxcb -lxcb-xkb -lxkbcommon -lxkbcommon-x11
//#endif
	// #cgo windows LDFLAGS: -lgdi32 -luser32 -lpng -lz
	#cgo windows LDFLAGS: -lgdi32 -luser32
// #include <AppKit/NSEvent.h>
#include "screen/goScreen.h"
#include "mouse/goMouse.h"
#include "key/goKey.h"
//#include "event/goEvent.h"
#include "window/goWindow.h"
*/
import "C"

import (
	"fmt"
	"image"

	// "os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	// "syscall"
	"math/rand"

	"github.com/go-vgo/robotgo/clipboard"
	"github.com/vcaesar/tt"
)

const (
	// Version get the robotgo version
	Version = "v0.100.0.1189, MT. Baker!"
)

// GetVersion get the robotgo version
func GetVersion() string {
	return Version
}

var (
	// MouseSleep set the mouse default millisecond sleep time
	MouseSleep = 0
	// KeySleep set the key default millisecond sleep time
	KeySleep = 0
)

type (
	// Map a map[string]interface{}
	Map map[string]interface{}
	// CHex define CHex as c rgb Hex type (C.MMRGBHex)
	CHex C.MMRGBHex
	// CBitmap define CBitmap as C.MMBitmapRef type
	CBitmap C.MMBitmapRef
)

// Bitmap is Bitmap struct
type Bitmap struct {
	ImgBuf        *uint8
	Width, Height int

	Bytewidth     int
	BitsPixel     uint8
	BytesPerPixel uint8
}

// Point is point struct
type Point struct {
	X int
	Y int
}

// Size is size structure
type Size struct {
	W, H int
}

// Rect is rect structure
type Rect struct {
	Point
	Size
}

// Try handler(err)
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// MilliSleep sleep tm milli second
func MilliSleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Millisecond)
}

// Sleep time.Sleep tm second
func Sleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Second)
}

// MicroSleep time C.microsleep(tm)
func MicroSleep(tm float64) {
	C.microsleep(C.double(tm))
}

// GoString trans C.char to string
func GoString(char *C.char) string {
	return C.GoString(char)
}

/*
      _______.  ______ .______       _______  _______ .__   __.
    /       | /      ||   _  \     |   ____||   ____||  \ |  |
   |   (----`|  ,----'|  |_)  |    |  |__   |  |__   |   \|  |
    \   \    |  |     |      /     |   __|  |   __|  |  . `  |
.----)   |   |  `----.|  |\  \----.|  |____ |  |____ |  |\   |
|_______/     \______|| _| `._____||_______||_______||__| \__|
*/

// ToMMRGBHex trans CHex to C.MMRGBHex
func ToMMRGBHex(hex CHex) C.MMRGBHex {
	return C.MMRGBHex(hex)
}

// UintToHex trans uint32 to robotgo.CHex
func UintToHex(u uint32) CHex {
	hex := U32ToHex(C.uint32_t(u))
	return CHex(hex)
}

// U32ToHex trans C.uint32_t to C.MMRGBHex
func U32ToHex(hex C.uint32_t) C.MMRGBHex {
	return C.MMRGBHex(hex)
}

// U8ToHex trans *C.uint8_t to C.MMRGBHex
func U8ToHex(hex *C.uint8_t) C.MMRGBHex {
	return C.MMRGBHex(*hex)
}

// PadHex trans C.MMRGBHex to string
func PadHex(hex C.MMRGBHex) string {
	color := C.pad_hex(hex)
	gcolor := C.GoString(color)
	C.free(unsafe.Pointer(color))

	return gcolor
}

// HexToRgb trans hex to rgb
func HexToRgb(hex uint32) *C.uint8_t {
	return C.color_hex_to_rgb(C.uint32_t(hex))
}

// RgbToHex trans rgb to hex
func RgbToHex(r, g, b uint8) C.uint32_t {
	return C.color_rgb_to_hex(C.uint8_t(r), C.uint8_t(g), C.uint8_t(b))
}

// GetPxColor get the pixel color return C.MMRGBHex
func GetPxColor(x, y int) C.MMRGBHex {
	cx := C.int32_t(x)
	cy := C.int32_t(y)

	color := C.get_px_color(cx, cy)
	return color
}

// GetPixelColor get the pixel color return string
func GetPixelColor(x, y int) string {
	cx := C.int32_t(x)
	cy := C.int32_t(y)

	color := C.get_pixel_color(cx, cy)
	gcolor := C.GoString(color)
	C.free(unsafe.Pointer(color))

	return gcolor
}

// GetMouseColor get the mouse pos's color
func GetMouseColor() string {
	x, y := GetMousePos()
	return GetPixelColor(x, y)
}

// ScaleX get the primary display horizontal DPI scale factor
func ScaleX() int {
	return int(C.scale_x())
}

// SysScale get the sys scale
func SysScale() float64 {
	s := C.sys_scale()
	return float64(s)
}

// Scaled return x * sys_scale
func Scaled(x int) int {
	return int(float64(x) * SysScale())
}

// ScaleY get primary display vertical DPI scale factor
func ScaleY() int {
	return int(C.scale_y())
}

// GetScreenSize get the screen size
func GetScreenSize() (int, int) {
	size := C.get_screen_size()
	// fmt.Println("...", size, size.width)
	return int(size.w), int(size.h)
}

// GetScreenRect get the screen rect
func GetScreenRect(displayId ...int) Rect {
	display := 0
	if len(displayId) > 0 {
		display = displayId[0]
	}

	rect := C.getScreenRect(C.int32_t(display))
	return Rect{
		Point{
			X: int(rect.origin.x),
			Y: int(rect.origin.y),
		},
		Size{
			W: int(rect.size.w),
			H: int(rect.size.h),
		},
	}
}

// Scale get the screen scale
func Scale() int {
	dpi := map[int]int{
		0: 100,
		// DPI Scaling Level
		96:  100,
		120: 125,
		144: 150,
		168: 175,
		192: 200,
		216: 225,
		// Custom DPI
		240: 250,
		288: 300,
		384: 400,
		480: 500,
	}

	x := ScaleX()
	return dpi[x]
}

// Scale0 return ScaleX() / 0.96
func Scale0() int {
	return int(float64(ScaleX()) / 0.96)
}

// Mul mul the scale
func Mul(x int) int {
	s := Scale()
	return x * s / 100
}

// GetScaleSize get the screen scale size
func GetScaleSize() (int, int) {
	x, y := GetScreenSize()
	s := Scale()
	return x * s / 100, y * s / 100
}

// SetXDisplayName set XDisplay name (Linux)
func SetXDisplayName(name string) string {
	cname := C.CString(name)
	str := C.set_XDisplay_name(cname)

	gstr := C.GoString(str)
	C.free(unsafe.Pointer(cname))

	return gstr
}

// GetXDisplayName get XDisplay name (Linux)
func GetXDisplayName() string {
	name := C.get_XDisplay_name()
	gname := C.GoString(name)
	C.free(unsafe.Pointer(name))

	return gname
}

// CaptureScreen capture the screen return bitmap(c struct),
// use `defer robotgo.FreeBitmap(bitmap)` to free the bitmap
//
// robotgo.CaptureScreen(x, y, w, h int)
func CaptureScreen(args ...int) C.MMBitmapRef {
	var x, y, w, h C.int32_t

	if len(args) > 3 {
		x = C.int32_t(args[0])
		y = C.int32_t(args[1])
		w = C.int32_t(args[2])
		h = C.int32_t(args[3])
	} else {
		x = 0
		y = 0

		// Get the main screen size.
		displaySize := C.getMainDisplaySize()
		w = displaySize.w
		h = displaySize.h
	}

	bit := C.capture_screen(x, y, w, h)
	return bit
}

// GoCaptureScreen capture the screen and return bitmap(go struct)
func GoCaptureScreen(args ...int) Bitmap {
	bit := CaptureScreen(args...)
	defer FreeBitmap(bit)

	return ToBitmap(bit)
}

// CaptureImg capture the screen and return image.Image
func CaptureImg(args ...int) image.Image {
	bit := CaptureScreen(args...)
	defer FreeBitmap(bit)

	return ToImage(bit)
}

// FreeBitmap free and dealloc the C bitmap
func FreeBitmap(bitmap C.MMBitmapRef) {
	// C.destroyMMBitmap(bitmap)
	C.bitmap_dealloc(bitmap)
}

// ToBitmap trans C.MMBitmapRef to Bitmap
func ToBitmap(bit C.MMBitmapRef) Bitmap {
	bitmap := Bitmap{
		ImgBuf:        (*uint8)(bit.imageBuffer),
		Width:         int(bit.width),
		Height:        int(bit.height),
		Bytewidth:     int(bit.bytewidth),
		BitsPixel:     uint8(bit.bitsPerPixel),
		BytesPerPixel: uint8(bit.bytesPerPixel),
	}

	return bitmap
}

// ToImage convert C.MMBitmapRef to standard image.Image
func ToImage(bit C.MMBitmapRef) image.Image {
	return ToRGBA(bit)
}

// ToRGBA convert C.MMBitmapRef to standard image.RGBA
func ToRGBA(bit C.MMBitmapRef) *image.RGBA {
	bmp1 := ToBitmap(bit)
	return ToRGBAGo(bmp1)
}

/*
.___  ___.   ______    __    __       _______. _______
|   \/   |  /  __  \  |  |  |  |     /       ||   ____|
|  \  /  | |  |  |  | |  |  |  |    |   (----`|  |__
|  |\/|  | |  |  |  | |  |  |  |     \   \    |   __|
|  |  |  | |  `--'  | |  `--'  | .----)   |   |  |____
|__|  |__|  \______/   \______/  |_______/    |_______|

*/

// CheckMouse check the mouse button
func CheckMouse(btn string) C.MMMouseButton {
	// button = args[0].(C.MMMouseButton)
	if btn == "left" {
		return C.LEFT_BUTTON
	}

	if btn == "center" {
		return C.CENTER_BUTTON
	}

	if btn == "right" {
		return C.RIGHT_BUTTON
	}

	return C.LEFT_BUTTON
}

// MoveMouse move the mouse
func MoveMouse(x, y int) {
	Move(x, y)
}

// Move move the mouse to (x, y)
func Move(x, y int) {
	cx := C.int32_t(x)
	cy := C.int32_t(y)
	C.move_mouse(cx, cy)

	MilliSleep(MouseSleep)
}

// DragMouse drag the mouse to (x, y)
func DragMouse(x, y int, args ...string) {
	Drag(x, y, args...)
}

// Drag drag the mouse to (x, y)
func Drag(x, y int, args ...string) {
	var button C.MMMouseButton = C.LEFT_BUTTON
	cx := C.int32_t(x)
	cy := C.int32_t(y)

	if len(args) > 0 {
		button = CheckMouse(args[0])
	}

	C.drag_mouse(cx, cy, button)
	MilliSleep(MouseSleep)
}

// DragSmooth drag the mouse smooth
func DragSmooth(x, y int, args ...interface{}) {
	MouseToggle("down")
	MoveSmooth(x, y, args...)
	MouseToggle("up")
}

// MoveMouseSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
func MoveMouseSmooth(x, y int, args ...interface{}) bool {
	return MoveSmooth(x, y, args...)
}

// MoveSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
//
// robotgo.MoveSmooth(x, y int, low, high float64, mouseDelay int)
func MoveSmooth(x, y int, args ...interface{}) bool {
	cx := C.int32_t(x)
	cy := C.int32_t(y)

	var (
		mouseDelay = 10
		low        C.double
		high       C.double
	)

	if len(args) > 2 {
		mouseDelay = args[2].(int)
	}

	if len(args) > 1 {
		low = C.double(args[0].(float64))
		high = C.double(args[1].(float64))
	} else {
		low = 1.0
		high = 3.0
	}

	cbool := C.move_mouse_smooth(cx, cy, low, high, C.int(mouseDelay))
	MilliSleep(MouseSleep)

	return bool(cbool)
}

// MoveArgs move mouse relative args
func MoveArgs(x, y int) (int, int) {
	mx, my := GetMousePos()
	mx = mx + x
	my = my + y

	return mx, my
}

// MoveRelative move mouse with relative
func MoveRelative(x, y int) {
	Move(MoveArgs(x, y))
}

// MoveSmoothRelative move mouse smooth with relative
func MoveSmoothRelative(x, y int, args ...interface{}) {
	mx, my := MoveArgs(x, y)
	MoveSmooth(mx, my, args...)
}

// GetMousePos get mouse's portion
func GetMousePos() (int, int) {
	pos := C.get_mouse_pos()

	x := int(pos.x)
	y := int(pos.y)

	return x, y
}

// MouseClick click the mouse
//
// robotgo.MouseClick(button string, double bool)
func MouseClick(args ...interface{}) {
	Click(args...)
}

// Click click the mouse
//
// robotgo.Click(button string, double bool)
func Click(args ...interface{}) {
	var (
		button C.MMMouseButton = C.LEFT_BUTTON
		double C.bool
	)

	if len(args) > 0 {
		button = CheckMouse(args[0].(string))
	}

	if len(args) > 1 {
		double = C.bool(args[1].(bool))
	}

	C.mouse_click(button, double)
	MilliSleep(MouseSleep)
}

// MoveClick move and click the mouse
//
// robotgo.MoveClick(x, y int, button string, double bool)
func MoveClick(x, y int, args ...interface{}) {
	MoveMouse(x, y)
	MouseClick(args...)
}

// MovesClick move smooth and click the mouse
func MovesClick(x, y int, args ...interface{}) {
	MoveSmooth(x, y)
	MouseClick(args...)
}

// MouseToggle toggle the mouse
func MouseToggle(togKey string, args ...interface{}) int {
	var button C.MMMouseButton = C.LEFT_BUTTON

	if len(args) > 0 {
		button = CheckMouse(args[0].(string))
	}

	down := C.CString(togKey)
	i := C.mouse_toggle(down, button)

	C.free(unsafe.Pointer(down))
	MilliSleep(MouseSleep)
	return int(i)
}

// ScrollMouse scroll the mouse
func ScrollMouse(x int, direction string) {
	cx := C.size_t(x)
	cy := C.CString(direction)
	C.scroll_mouse(cx, cy)

	C.free(unsafe.Pointer(cy))
	MilliSleep(MouseSleep)
}

// Scroll scroll the mouse to (x, y)
//
// robotgo.Scroll(x, y, msDelay int)
func Scroll(x, y int, args ...int) {
	var msDelay = 10
	if len(args) > 0 {
		msDelay = args[0]
	}

	cx := C.int(x)
	cy := C.int(y)
	cz := C.int(msDelay)

	C.scroll(cx, cy, cz)
	MilliSleep(MouseSleep)
}

// ScrollRelative scroll mouse with relative
func ScrollRelative(x, y int, args ...int) {
	mx, my := MoveArgs(x, y)
	Scroll(mx, my, args...)
}

// SetMouseDelay set mouse delay
func SetMouseDelay(delay int) {
	cdelay := C.size_t(delay)
	C.set_mouse_delay(cdelay)
}

/*
 __  ___  ___________    ____ .______     ______        ___      .______       _______
|  |/  / |   ____\   \  /   / |   _  \   /  __  \      /   \     |   _  \     |       \
|  '  /  |  |__   \   \/   /  |  |_)  | |  |  |  |    /  ^  \    |  |_)  |    |  .--.  |
|    <   |   __|   \_    _/   |   _  <  |  |  |  |   /  /_\  \   |      /     |  |  |  |
|  .  \  |  |____    |  |     |  |_)  | |  `--'  |  /  _____  \  |  |\  \----.|  '--'  |
|__|\__\ |_______|   |__|     |______/   \______/  /__/     \__\ | _| `._____||_______/

*/

// KeyTap tap the keyboard code;
//
// See keys:
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
//
func KeyTap(tapKey string, args ...interface{}) string {
	var (
		akey     string
		keyT     = "null"
		keyArr   []string
		num      int
		keyDelay = 10
	)

	// var ckeyArr []*C.char
	ckeyArr := make([](*C.char), 0)
	// zkey := C.CString(args[0])
	zkey := C.CString(tapKey)
	defer C.free(unsafe.Pointer(zkey))

	if len(args) > 2 && (reflect.TypeOf(args[2]) != reflect.TypeOf(num)) {
		num = len(args)
		for i := 0; i < num; i++ {
			s := args[i].(string)
			ckeyArr = append(ckeyArr, (*C.char)(unsafe.Pointer(C.CString(s))))
		}

		str := C.key_Taps(zkey,
			(**C.char)(unsafe.Pointer(&ckeyArr[0])), C.int(num), 0)
		MilliSleep(KeySleep)
		return C.GoString(str)
	}

	if len(args) > 0 {
		if reflect.TypeOf(args[0]) == reflect.TypeOf(keyArr) {

			keyArr = args[0].([]string)
			num = len(keyArr)
			for i := 0; i < num; i++ {
				ckeyArr = append(ckeyArr, (*C.char)(unsafe.Pointer(C.CString(keyArr[i]))))
			}

			if len(args) > 1 {
				keyDelay = args[1].(int)
			}
		} else {
			akey = args[0].(string)

			if len(args) > 1 {
				if reflect.TypeOf(args[1]) == reflect.TypeOf(akey) {
					keyT = args[1].(string)
					if len(args) > 2 {
						keyDelay = args[2].(int)
					}
				} else {
					keyDelay = args[1].(int)
				}
			}
		}

	} else {
		akey = "null"
		keyArr = []string{"null"}
	}

	if akey == "" && len(keyArr) != 0 {
		str := C.key_Taps(zkey, (**C.char)(unsafe.Pointer(&ckeyArr[0])),
			C.int(num), C.int(keyDelay))

		MilliSleep(KeySleep)
		return C.GoString(str)
	}

	amod := C.CString(akey)
	amodt := C.CString(keyT)
	str := C.key_tap(zkey, amod, amodt, C.int(keyDelay))

	C.free(unsafe.Pointer(amod))
	C.free(unsafe.Pointer(amodt))

	MilliSleep(KeySleep)
	return C.GoString(str)
}

// KeyToggle toggle the keyboard
//
// See keys:
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
//
func KeyToggle(key string, args ...string) string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	ckeyArr := make([](*C.char), 0)
	if len(args) > 3 {
		num := len(args)
		for i := 0; i < num; i++ {
			ckeyArr = append(ckeyArr, (*C.char)(unsafe.Pointer(C.CString(args[i]))))
		}

		str := C.key_Toggles(ckey, (**C.char)(unsafe.Pointer(&ckeyArr[0])), C.int(num))
		MilliSleep(KeySleep)
		return C.GoString(str)
	}

	var (
		down, mKey, mKeyT = "null", "null", "null"
		// keyDelay = 10
	)

	if len(args) > 0 {
		down = args[0]

		if len(args) > 1 {
			mKey = args[1]
			if len(args) > 2 {
				mKeyT = args[2]
			}
		}
	}

	cdown := C.CString(down)
	cmKey := C.CString(mKey)
	cmKeyT := C.CString(mKeyT)

	str := C.key_toggle(ckey, cdown, cmKey, cmKeyT)
	// str := C.key_Toggle(ckey, cdown, cmKey, cmKeyT, C.int(keyDelay))

	C.free(unsafe.Pointer(cdown))
	C.free(unsafe.Pointer(cmKey))
	C.free(unsafe.Pointer(cmKeyT))

	MilliSleep(KeySleep)
	return C.GoString(str)
}

// KeyPress press key string
func KeyPress(key string) {
	KeyDown(key)
	Sleep(1 + rand.Intn(3))
	KeyUp(key)
}

// KeyDown press down a key
func KeyDown(key string) {
	KeyToggle(key, "down")
}

// KeyUp press up a key
func KeyUp(key string) {
	KeyToggle(key, "up")
}

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// CharCodeAt char code at utf-8
func CharCodeAt(s string, n int) rune {
	i := 0
	for _, r := range s {
		if i == n {
			return r
		}
		i++
	}

	return 0
}

// UnicodeType tap uint32 unicode
func UnicodeType(str uint32) {
	cstr := C.uint(str)
	C.unicodeType(cstr)
}

// ToUC trans string to unicode []string
func ToUC(text string) []string {
	var uc []string

	for _, r := range text {
		textQ := strconv.QuoteToASCII(string(r))
		textUnQ := textQ[1 : len(textQ)-1]

		st := strings.Replace(textUnQ, "\\u", "U", -1)
		if st == "\\\\" {
			st = "\\"
		}
		if st == `\"` {
			st = `"`
		}
		uc = append(uc, st)
	}

	return uc
}

func inputUTF(str string) {
	cstr := C.CString(str)
	C.input_utf(cstr)

	C.free(unsafe.Pointer(cstr))
}

// TypeStr send a string, support UTF-8
//
// robotgo.TypeStr(string: The string to send, float64: microsleep time, x11)
func TypeStr(str string, args ...float64) {
	var tm, tm1 = 0.0, 7.0

	if len(args) > 0 {
		tm = args[0]
	}
	if len(args) > 1 {
		tm1 = args[1]
	}

	if runtime.GOOS == "linux" {
		strUc := ToUC(str)
		for i := 0; i < len(strUc); i++ {
			ru := []rune(strUc[i])
			if len(ru) <= 1 {
				ustr := uint32(CharCodeAt(strUc[i], 0))
				UnicodeType(ustr)
			} else {
				inputUTF(strUc[i])
				MicroSleep(tm1)
			}

			MicroSleep(tm)
		}

		return
	}

	for i := 0; i < len([]rune(str)); i++ {
		ustr := uint32(CharCodeAt(str, i))
		UnicodeType(ustr)

		// if len(args) > 0 {
		MicroSleep(tm)
		// }
	}
	MilliSleep(KeySleep)
}

// PasteStr paste a string, support UTF-8
func PasteStr(str string) string {
	err := clipboard.WriteAll(str)
	if err != nil {
		return fmt.Sprint(err)
	}

	if runtime.GOOS == "darwin" {
		return KeyTap("v", "command")
	}

	return KeyTap("v", "control")
}

// Deprecated: TypeString send a string, support unicode
// TypeStr(string: The string to send), Wno-deprecated
func TypeString(str string, delay ...int) {
	tt.Drop("TypeString", "TypeStr")
	var cdelay C.size_t
	cstr := C.CString(str)
	if len(delay) > 0 {
		cdelay = C.size_t(delay[0])
	}

	C.type_string_delayed(cstr, cdelay)

	C.free(unsafe.Pointer(cstr))
}

// TypeStrDelay type string delayed
func TypeStrDelay(str string, delay int) {
	TypeStr(str)
	Sleep(delay)
}

// Deprecated: TypeStringDelayed type string delayed, Wno-deprecated
func TypeStringDelayed(str string, delay int) {
	tt.Drop("TypeStringDelayed", "TypeStrDelay")
	TypeStrDelay(str, delay)
}

// SetKeyDelay set keyboard delay
func SetKeyDelay(delay int) {
	C.set_keyboard_delay(C.size_t(delay))
}

// Deprecated: SetKeyboardDelay set keyboard delay, Wno-deprecated,
// this function will be removed in version v1.0.0
func SetKeyboardDelay(delay int) {
	tt.Drop("SetKeyboardDelay", "SetKeyDelay")
	SetKeyDelay(delay)
}

// SetDelay set the key and mouse delay
func SetDelay(d ...int) {
	v := 10
	if len(d) > 0 {
		v = d[0]
	}

	SetMouseDelay(v)
	SetKeyDelay(v)
}

/*
____    __    ____  __  .__   __.  _______   ______   ____    __    ____
\   \  /  \  /   / |  | |  \ |  | |       \ /  __  \  \   \  /  \  /   /
 \   \/    \/   /  |  | |   \|  | |  .--.  |  |  |  |  \   \/    \/   /
  \            /   |  | |  . `  | |  |  |  |  |  |  |   \            /
   \    /\    /    |  | |  |\   | |  '--'  |  `--'  |    \    /\    /
    \__/  \__/     |__| |__| \__| |_______/ \______/      \__/  \__/

*/

// ShowAlert show a alert window
func ShowAlert(title, msg string, args ...string) bool {
	var (
		// title         string
		// msg           string
		defaultBtn = "Ok"
		cancelBtn  = "Cancel"
	)

	if len(args) > 0 {
		// title = args[0]
		// msg = args[1]
		defaultBtn = args[0]
	}

	if len(args) > 1 {
		cancelBtn = args[1]
	}

	cTitle := C.CString(title)
	cMsg := C.CString(msg)
	defaultButton := C.CString(defaultBtn)
	cancelButton := C.CString(cancelBtn)

	cbool := C.show_alert(cTitle, cMsg, defaultButton, cancelButton)
	ibool := int(cbool)

	C.free(unsafe.Pointer(cTitle))
	C.free(unsafe.Pointer(cMsg))
	C.free(unsafe.Pointer(defaultButton))
	C.free(unsafe.Pointer(cancelButton))

	return ibool == 0
}

// IsValid valid the window
func IsValid() bool {
	abool := C.is_valid()
	gbool := bool(abool)
	// fmt.Println("bool---------", gbool)
	return gbool
}

// SetActive set the window active
func SetActive(win C.MData) {
	C.set_active(win)
}

// GetActive get the active window
func GetActive() C.MData {
	mdata := C.get_active()
	// fmt.Println("active----", mdata)
	return mdata
}

// MinWindow set the window min
func MinWindow(pid int32, args ...interface{}) {
	var (
		state = true
		hwnd  int
	)

	if len(args) > 0 {
		state = args[0].(bool)
	}
	if len(args) > 1 {
		hwnd = args[1].(int)
	}

	C.min_window(C.uintptr(pid), C.bool(state), C.uintptr(hwnd))
}

// MaxWindow set the window max
func MaxWindow(pid int32, args ...interface{}) {
	var (
		state = true
		hwnd  int
	)

	if len(args) > 0 {
		state = args[0].(bool)
	}
	if len(args) > 1 {
		hwnd = args[1].(int)
	}

	C.max_window(C.uintptr(pid), C.bool(state), C.uintptr(hwnd))
}

// CloseWindow close the window
func CloseWindow(args ...int32) {
	if len(args) <= 0 {
		C.close_main_window()
		return
	}

	var hwnd, isHwnd int32
	if len(args) > 0 {
		hwnd = args[0]
	}
	if len(args) > 1 {
		isHwnd = args[1]
	}

	C.close_window(C.uintptr(hwnd), C.uintptr(isHwnd))
}

// SetHandle set the window handle
func SetHandle(hwnd int) {
	chwnd := C.uintptr(hwnd)
	C.set_handle(chwnd)
}

// SetHandlePid set the window handle by pid
func SetHandlePid(pid int32, args ...int32) {
	var isHwnd int32
	if len(args) > 0 {
		isHwnd = args[0]
	}

	C.set_handle_pid_mData(C.uintptr(pid), C.uintptr(isHwnd))
}

// GetHandPid get handle mdata by pid
func GetHandPid(pid int32, args ...int32) C.MData {
	var isHwnd int32
	if len(args) > 0 {
		isHwnd = args[0]
	}

	return C.set_handle_pid(C.uintptr(pid), C.uintptr(isHwnd))
}

// GetHandle get the window handle
func GetHandle() int {
	hwnd := C.get_handle()
	ghwnd := int(hwnd)
	// fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

// Deprecated: GetBHandle get the window handle, Wno-deprecated
func GetBHandle() int {
	tt.Drop("GetBHandle", "GetHandle")
	hwnd := C.bget_handle()
	ghwnd := int(hwnd)
	//fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

func cgetTitle(hwnd, isHwnd int32) string {
	title := C.get_title(C.uintptr(hwnd), C.uintptr(isHwnd))
	gtitle := C.GoString(title)

	return gtitle
}

// GetTitle get the window title
func GetTitle(args ...int32) string {
	if len(args) <= 0 {
		title := C.get_main_title()
		gtitle := C.GoString(title)
		return gtitle
	}

	if len(args) > 1 {
		return internalGetTitle(args[0], args[1])
	}

	return internalGetTitle(args[0])
}

// GetPID get the process id
func GetPID() int32 {
	pid := C.get_PID()
	return int32(pid)
}

// internalGetBounds get the window bounds
func internalGetBounds(pid int32, hwnd int) (int, int, int, int) {
	bounds := C.get_bounds(C.uintptr(pid), C.uintptr(hwnd))
	return int(bounds.X), int(bounds.Y), int(bounds.W), int(bounds.H)
}

// Is64Bit determine whether the sys is 64bit
func Is64Bit() bool {
	b := C.Is64Bit()
	return bool(b)
}

func internalActive(pid int32, hwnd int) {
	C.active_PID(C.uintptr(pid), C.uintptr(hwnd))
}

// ActivePID active the window by PID,
// If args[0] > 0 on the Windows platform via a window handle to active
// func ActivePID(pid int32, args ...int) {
// 	var hwnd int
// 	if len(args) > 0 {
// 		hwnd = args[0]
// 	}

// 	C.active_PID(C.uintptr(pid), C.uintptr(hwnd))
// }

// ActiveName active window by name
func ActiveName(name string) error {
	pids, err := FindIds(name)
	if err == nil && len(pids) > 0 {
		return ActivePID(pids[0])
	}

	return err
}
