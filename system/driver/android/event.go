//go:build android

package android

import "C"

import (
	"image"
	"log"

	"cogentcore.org/core/events"
	"cogentcore.org/core/events/key"
)

func keyboardTyped(str *C.char) {
	for _, r := range C.GoString(str) {
		code := ConvAndroidKeyCode(r)
		TheApp.Event.KeyChord(r, code, 0)
	}
}

func keyboardDelete() {
	TheApp.Event.KeyChord(0, key.CodeBackspace, 0)
}

func scaled(scaleFactor, posX, posY C.float) {
	where := image.Pt(int(posX), int(posY))
	TheApp.Event.Magnify(float32(scaleFactor), where)
}

func ProcessEvents(env *C.JNIEnv, q *C.AInputQueue) {
	var e *C.AInputEvent
	for C.AInputQueue_getEvent(q, &e) >= 0 {
		if C.AInputQueue_preDispatchEvent(q, e) != 0 {
			continue
		}
		ProcessEvent(env, e)
		C.AInputQueue_finishEvent(q, e, 0)
	}
}

func ProcessEvent(env *C.JNIEnv, e *C.AInputEvent) {
	switch C.AInputEvent_getType(e) {
	case C.AINPUT_EVENT_TYPE_KEY:
		ProcessKey(env, e)
	case C.AINPUT_EVENT_TYPE_MOTION:

		upDownIndex := C.size_t(C.AMotionEvent_getAction(e)&C.AMOTION_EVENT_ACTION_POINTER_INDEX_MASK) >> C.AMOTION_EVENT_ACTION_POINTER_INDEX_SHIFT
		upDownType := events.TouchMove
		switch C.AMotionEvent_getAction(e) & C.AMOTION_EVENT_ACTION_MASK {
		case C.AMOTION_EVENT_ACTION_DOWN, C.AMOTION_EVENT_ACTION_POINTER_DOWN:
			upDownType = events.TouchStart
		case C.AMOTION_EVENT_ACTION_UP, C.AMOTION_EVENT_ACTION_POINTER_UP:
			upDownType = events.TouchEnd
		}

		for i, n := C.size_t(0), C.AMotionEvent_getPointerCount(e); i < n; i++ {
			t := events.TouchMove
			if i == upDownIndex {
				t = upDownType
			}
			seq := events.Sequence(C.AMotionEvent_getPointerId(e, i))
			x := int(C.AMotionEvent_getX(e, i))
			y := int(C.AMotionEvent_getY(e, i))
			TheApp.Event.Touch(t, seq, image.Pt(x, y))
		}
	default:
		log.Printf("unknown input event, type=%d", C.AInputEvent_getType(e))
	}
}

func ProcessKey(env *C.JNIEnv, e *C.AInputEvent) {
	deviceID := C.AInputEvent_getDeviceId(e)
	if deviceID == 0 {
		return
	}

	r := rune(C.getKeyRune(env, e))
	code := ConvAndroidKeyCode(int32(C.AKeyEvent_getKeyCode(e)))

	if r >= '0' && r <= '9' {
		return
	}
	typ := events.KeyDown
	if C.AKeyEvent_getAction(e) == C.AKEY_STATE_UP {
		typ = events.KeyUp
	}

	TheApp.Event.Key(typ, r, code, 0)
}

var AndroidKeyCodes = map[int32]key.Codes{
	C.AKEYCODE_HOME:            key.CodeHome,
	C.AKEYCODE_0:               key.Code0,
	C.AKEYCODE_1:               key.Code1,
	C.AKEYCODE_2:               key.Code2,
	C.AKEYCODE_3:               key.Code3,
	C.AKEYCODE_4:               key.Code4,
	C.AKEYCODE_5:               key.Code5,
	C.AKEYCODE_6:               key.Code6,
	C.AKEYCODE_7:               key.Code7,
	C.AKEYCODE_8:               key.Code8,
	C.AKEYCODE_9:               key.Code9,
	C.AKEYCODE_VOLUME_UP:       key.CodeVolumeUp,
	C.AKEYCODE_VOLUME_DOWN:     key.CodeVolumeDown,
	C.AKEYCODE_A:               key.CodeA,
	C.AKEYCODE_B:               key.CodeB,
	C.AKEYCODE_C:               key.CodeC,
	C.AKEYCODE_D:               key.CodeD,
	C.AKEYCODE_E:               key.CodeE,
	C.AKEYCODE_F:               key.CodeF,
	C.AKEYCODE_G:               key.CodeG,
	C.AKEYCODE_H:               key.CodeH,
	C.AKEYCODE_I:               key.CodeI,
	C.AKEYCODE_J:               key.CodeJ,
	C.AKEYCODE_K:               key.CodeK,
	C.AKEYCODE_L:               key.CodeL,
	C.AKEYCODE_M:               key.CodeM,
	C.AKEYCODE_N:               key.CodeN,
	C.AKEYCODE_O:               key.CodeO,
	C.AKEYCODE_P:               key.CodeP,
	C.AKEYCODE_Q:               key.CodeQ,
	C.AKEYCODE_R:               key.CodeR,
	C.AKEYCODE_S:               key.CodeS,
	C.AKEYCODE_T:               key.CodeT,
	C.AKEYCODE_U:               key.CodeU,
	C.AKEYCODE_V:               key.CodeV,
	C.AKEYCODE_W:               key.CodeW,
	C.AKEYCODE_X:               key.CodeX,
	C.AKEYCODE_Y:               key.CodeY,
	C.AKEYCODE_Z:               key.CodeZ,
	C.AKEYCODE_COMMA:           key.CodeComma,
	C.AKEYCODE_PERIOD:          key.CodeFullStop,
	C.AKEYCODE_ALT_LEFT:        key.CodeLeftAlt,
	C.AKEYCODE_ALT_RIGHT:       key.CodeRightAlt,
	C.AKEYCODE_SHIFT_LEFT:      key.CodeLeftShift,
	C.AKEYCODE_SHIFT_RIGHT:     key.CodeRightShift,
	C.AKEYCODE_TAB:             key.CodeTab,
	C.AKEYCODE_SPACE:           key.CodeSpacebar,
	C.AKEYCODE_ENTER:           key.CodeReturnEnter,
	C.AKEYCODE_DEL:             key.CodeBackspace,
	C.AKEYCODE_GRAVE:           key.CodeGraveAccent,
	C.AKEYCODE_MINUS:           key.CodeHyphenMinus,
	C.AKEYCODE_EQUALS:          key.CodeEqualSign,
	C.AKEYCODE_LEFT_BRACKET:    key.CodeLeftSquareBracket,
	C.AKEYCODE_RIGHT_BRACKET:   key.CodeRightSquareBracket,
	C.AKEYCODE_BACKSLASH:       key.CodeBackslash,
	C.AKEYCODE_SEMICOLON:       key.CodeSemicolon,
	C.AKEYCODE_APOSTROPHE:      key.CodeApostrophe,
	C.AKEYCODE_SLASH:           key.CodeSlash,
	C.AKEYCODE_PAGE_UP:         key.CodePageUp,
	C.AKEYCODE_PAGE_DOWN:       key.CodePageDown,
	C.AKEYCODE_ESCAPE:          key.CodeEscape,
	C.AKEYCODE_FORWARD_DEL:     key.CodeDelete,
	C.AKEYCODE_CTRL_LEFT:       key.CodeLeftControl,
	C.AKEYCODE_CTRL_RIGHT:      key.CodeRightControl,
	C.AKEYCODE_CAPS_LOCK:       key.CodeCapsLock,
	C.AKEYCODE_META_LEFT:       key.CodeLeftMeta,
	C.AKEYCODE_META_RIGHT:      key.CodeRightMeta,
	C.AKEYCODE_INSERT:          key.CodeInsert,
	C.AKEYCODE_F1:              key.CodeF1,
	C.AKEYCODE_F2:              key.CodeF2,
	C.AKEYCODE_F3:              key.CodeF3,
	C.AKEYCODE_F4:              key.CodeF4,
	C.AKEYCODE_F5:              key.CodeF5,
	C.AKEYCODE_F6:              key.CodeF6,
	C.AKEYCODE_F7:              key.CodeF7,
	C.AKEYCODE_F8:              key.CodeF8,
	C.AKEYCODE_F9:              key.CodeF9,
	C.AKEYCODE_F10:             key.CodeF10,
	C.AKEYCODE_F11:             key.CodeF11,
	C.AKEYCODE_F12:             key.CodeF12,
	C.AKEYCODE_NUM_LOCK:        key.CodeKeypadNumLock,
	C.AKEYCODE_NUMPAD_0:        key.CodeKeypad0,
	C.AKEYCODE_NUMPAD_1:        key.CodeKeypad1,
	C.AKEYCODE_NUMPAD_2:        key.CodeKeypad2,
	C.AKEYCODE_NUMPAD_3:        key.CodeKeypad3,
	C.AKEYCODE_NUMPAD_4:        key.CodeKeypad4,
	C.AKEYCODE_NUMPAD_5:        key.CodeKeypad5,
	C.AKEYCODE_NUMPAD_6:        key.CodeKeypad6,
	C.AKEYCODE_NUMPAD_7:        key.CodeKeypad7,
	C.AKEYCODE_NUMPAD_8:        key.CodeKeypad8,
	C.AKEYCODE_NUMPAD_9:        key.CodeKeypad9,
	C.AKEYCODE_NUMPAD_DIVIDE:   key.CodeKeypadSlash,
	C.AKEYCODE_NUMPAD_MULTIPLY: key.CodeKeypadAsterisk,
	C.AKEYCODE_NUMPAD_SUBTRACT: key.CodeKeypadHyphenMinus,
	C.AKEYCODE_NUMPAD_ADD:      key.CodeKeypadPlusSign,
	C.AKEYCODE_NUMPAD_DOT:      key.CodeKeypadFullStop,
	C.AKEYCODE_NUMPAD_ENTER:    key.CodeKeypadEnter,
	C.AKEYCODE_NUMPAD_EQUALS:   key.CodeKeypadEqualSign,
	C.AKEYCODE_VOLUME_MUTE:     key.CodeMute,
}

func ConvAndroidKeyCode(aKeyCode int32) key.Codes {
	if code, ok := AndroidKeyCodes[aKeyCode]; ok {
		return code
	}
	return key.CodeUnknown
}
