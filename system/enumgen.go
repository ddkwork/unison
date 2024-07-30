package system

import (
	"cogentcore.org/core/enums"
)

var _PlatformsValues = []Platforms{0, 1, 2, 3, 4, 5, 6}

const PlatformsN Platforms = 7

var _PlatformsValueMap = map[string]Platforms{`MacOS`: 0, `Linux`: 1, `Windows`: 2, `IOS`: 3, `Android`: 4, `Web`: 5, `Offscreen`: 6}

var _PlatformsDescMap = map[Platforms]string{0: `MacOS is a Mac OS machine (aka Darwin)`, 1: `Linux is a Linux OS machine`, 2: `Windows is a Microsoft Windows machine`, 3: `IOS is an Apple iOS or iPadOS mobile phone or iPad`, 4: `Android is an Android mobile phone or tablet`, 5: `Web is a web browser running the app through WASM`, 6: `Offscreen is an offscreen driver typically used for testing, specified using the &#34;offscreen&#34; build tag`}

var _PlatformsMap = map[Platforms]string{0: `MacOS`, 1: `Linux`, 2: `Windows`, 3: `IOS`, 4: `Android`, 5: `Web`, 6: `Offscreen`}

func (i Platforms) String() string { return enums.String(i, _PlatformsMap) }

func (i *Platforms) SetString(s string) error {
	return enums.SetString(i, s, _PlatformsValueMap, "Platforms")
}

func (i Platforms) Int64() int64 { return int64(i) }

func (i *Platforms) SetInt64(in int64) { *i = Platforms(in) }

func (i Platforms) Desc() string { return enums.Desc(i, _PlatformsDescMap) }

func PlatformsValues() []Platforms { return _PlatformsValues }

func (i Platforms) Values() []enums.Enum { return enums.Values(_PlatformsValues) }

func (i Platforms) MarshalText() ([]byte, error) { return []byte(i.String()), nil }

func (i *Platforms) UnmarshalText(text []byte) error {
	return enums.UnmarshalText(i, text, "Platforms")
}

var _ScreenOrientationValues = []ScreenOrientation{0, 1, 2}

const ScreenOrientationN ScreenOrientation = 3

var _ScreenOrientationValueMap = map[string]ScreenOrientation{`OrientationUnknown`: 0, `Portrait`: 1, `Landscape`: 2}

var _ScreenOrientationDescMap = map[ScreenOrientation]string{0: `OrientationUnknown means device orientation cannot be determined. Equivalent on Android to Configuration.ORIENTATION_UNKNOWN and on iOS to: UIDeviceOrientationUnknown UIDeviceOrientationFaceUp UIDeviceOrientationFaceDown`, 1: `Portrait is a device oriented so it is tall and thin. Equivalent on Android to Configuration.ORIENTATION_PORTRAIT and on iOS to: UIDeviceOrientationPortrait UIDeviceOrientationPortraitUpsideDown`, 2: `Landscape is a device oriented so it is short and wide. Equivalent on Android to Configuration.ORIENTATION_LANDSCAPE and on iOS to: UIDeviceOrientationLandscapeLeft UIDeviceOrientationLandscapeRight`}

var _ScreenOrientationMap = map[ScreenOrientation]string{0: `OrientationUnknown`, 1: `Portrait`, 2: `Landscape`}

func (i ScreenOrientation) String() string { return enums.String(i, _ScreenOrientationMap) }

func (i *ScreenOrientation) SetString(s string) error {
	return enums.SetString(i, s, _ScreenOrientationValueMap, "ScreenOrientation")
}

func (i ScreenOrientation) Int64() int64 { return int64(i) }

func (i *ScreenOrientation) SetInt64(in int64) { *i = ScreenOrientation(in) }

func (i ScreenOrientation) Desc() string { return enums.Desc(i, _ScreenOrientationDescMap) }

func ScreenOrientationValues() []ScreenOrientation { return _ScreenOrientationValues }

func (i ScreenOrientation) Values() []enums.Enum { return enums.Values(_ScreenOrientationValues) }

func (i ScreenOrientation) MarshalText() ([]byte, error) { return []byte(i.String()), nil }

func (i *ScreenOrientation) UnmarshalText(text []byte) error {
	return enums.UnmarshalText(i, text, "ScreenOrientation")
}

var _WindowFlagsValues = []WindowFlags{0, 1, 2, 3, 4, 5}

const WindowFlagsN WindowFlags = 6

var _WindowFlagsValueMap = map[string]WindowFlags{`Dialog`: 0, `Modal`: 1, `Tool`: 2, `Fullscreen`: 3, `Minimized`: 4, `Focused`: 5}

var _WindowFlagsDescMap = map[WindowFlags]string{0: `Dialog indicates that this is a temporary, pop-up window.`, 1: `Modal indicates that this dialog window blocks events going to other windows until it is closed.`, 2: `Tool indicates that this is a floating tool window that has minimized window decoration.`, 3: `Fullscreen indicates a window that occupies the entire screen.`, 4: `Minimized indicates a window reduced to an icon, or otherwise no longer visible or active. Otherwise, the window should be assumed to be visible.`, 5: `Focused indicates that the window has the focus.`}

var _WindowFlagsMap = map[WindowFlags]string{0: `Dialog`, 1: `Modal`, 2: `Tool`, 3: `Fullscreen`, 4: `Minimized`, 5: `Focused`}

func (i WindowFlags) String() string { return enums.BitFlagString(i, _WindowFlagsValues) }

func (i WindowFlags) BitIndexString() string { return enums.String(i, _WindowFlagsMap) }

func (i *WindowFlags) SetString(s string) error { *i = 0; return i.SetStringOr(s) }

func (i *WindowFlags) SetStringOr(s string) error {
	return enums.SetStringOr(i, s, _WindowFlagsValueMap, "WindowFlags")
}

func (i WindowFlags) Int64() int64 { return int64(i) }

func (i *WindowFlags) SetInt64(in int64) { *i = WindowFlags(in) }

func (i WindowFlags) Desc() string { return enums.Desc(i, _WindowFlagsDescMap) }

func WindowFlagsValues() []WindowFlags { return _WindowFlagsValues }

func (i WindowFlags) Values() []enums.Enum { return enums.Values(_WindowFlagsValues) }

func (i WindowFlags) HasFlag(f enums.BitFlag) bool { return enums.HasFlag((*int64)(&i), f) }

func (i *WindowFlags) SetFlag(on bool, f ...enums.BitFlag) { enums.SetFlag((*int64)(i), on, f...) }

func (i WindowFlags) MarshalText() ([]byte, error) { return []byte(i.String()), nil }

func (i *WindowFlags) UnmarshalText(text []byte) error {
	return enums.UnmarshalText(i, text, "WindowFlags")
}
