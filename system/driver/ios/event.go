//go:build ios

package ios

import "C"

import (
	"image"

	"cogentcore.org/core/events"
	"cogentcore.org/core/events/key"
)

var TouchIDs [11]uintptr

func sendTouch(cTouch, cTouchType uintptr, x, y float32) {
	id := -1
	for i, val := range TouchIDs {
		if val == cTouch {
			id = i
			break
		}
	}
	if id == -1 {
		for i, val := range TouchIDs {
			if val == 0 {
				TouchIDs[i] = cTouch
				id = i
				break
			}
		}
		if id == -1 {
			panic("out of touchIDs")
		}
	}
	t := events.TouchStart
	switch cTouchType {
	case 0:
		t = events.TouchStart
	case 1:
		t = events.TouchMove
	case 2:
		t = events.TouchEnd

		for idx := range TouchIDs {
			TouchIDs[idx] = 0
		}
	}

	TheApp.Event.Touch(t, events.Sequence(id), image.Pt(int(x), int(y)))
}

func keyboardTyped(str *C.char) {
	for _, r := range C.GoString(str) {
		code := GetCodeFromRune(r)
		TheApp.Event.KeyChord(r, code, 0)
	}
}

func keyboardDelete() {
	TheApp.Event.KeyChord(0, key.CodeBackspace, 0)
}

func scaled(scaleFactor, posX, posY C.float) {
	where := image.Pt(int(posX), int(posY))
	sf := float32(scaleFactor)

	sf = 1 + (sf-1)/10
	TheApp.Event.Magnify(sf, where)
}

var CodeFromRune = map[rune]key.Codes{
	'0':  key.Code0,
	'1':  key.Code1,
	'2':  key.Code2,
	'3':  key.Code3,
	'4':  key.Code4,
	'5':  key.Code5,
	'6':  key.Code6,
	'7':  key.Code7,
	'8':  key.Code8,
	'9':  key.Code9,
	'a':  key.CodeA,
	'b':  key.CodeB,
	'c':  key.CodeC,
	'd':  key.CodeD,
	'e':  key.CodeE,
	'f':  key.CodeF,
	'g':  key.CodeG,
	'h':  key.CodeH,
	'i':  key.CodeI,
	'j':  key.CodeJ,
	'k':  key.CodeK,
	'l':  key.CodeL,
	'm':  key.CodeM,
	'n':  key.CodeN,
	'o':  key.CodeO,
	'p':  key.CodeP,
	'q':  key.CodeQ,
	'r':  key.CodeR,
	's':  key.CodeS,
	't':  key.CodeT,
	'u':  key.CodeU,
	'v':  key.CodeV,
	'w':  key.CodeW,
	'x':  key.CodeX,
	'y':  key.CodeY,
	'z':  key.CodeZ,
	'A':  key.CodeA,
	'B':  key.CodeB,
	'C':  key.CodeC,
	'D':  key.CodeD,
	'E':  key.CodeE,
	'F':  key.CodeF,
	'G':  key.CodeG,
	'H':  key.CodeH,
	'I':  key.CodeI,
	'J':  key.CodeJ,
	'K':  key.CodeK,
	'L':  key.CodeL,
	'M':  key.CodeM,
	'N':  key.CodeN,
	'O':  key.CodeO,
	'P':  key.CodeP,
	'Q':  key.CodeQ,
	'R':  key.CodeR,
	'S':  key.CodeS,
	'T':  key.CodeT,
	'U':  key.CodeU,
	'V':  key.CodeV,
	'W':  key.CodeW,
	'X':  key.CodeX,
	'Y':  key.CodeY,
	'Z':  key.CodeZ,
	',':  key.CodeComma,
	'.':  key.CodeFullStop,
	' ':  key.CodeSpacebar,
	'\n': key.CodeReturnEnter,
	'`':  key.CodeGraveAccent,
	'-':  key.CodeHyphenMinus,
	'=':  key.CodeEqualSign,
	'[':  key.CodeLeftSquareBracket,
	']':  key.CodeRightSquareBracket,
	'\\': key.CodeBackslash,
	';':  key.CodeSemicolon,
	'\'': key.CodeApostrophe,
	'/':  key.CodeSlash,
}

func GetCodeFromRune(r rune) key.Codes {
	if code, ok := CodeFromRune[r]; ok {
		return code
	}
	return key.CodeUnknown
}
