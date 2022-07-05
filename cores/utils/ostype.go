package utils

import (
	"runtime"
)

// OSType -
type OSType uint8
const (
    OSTypeNone OSType = iota
    OSTypeLinux
    OSTypeMacOS
    OSTypeWindows
    OSTypeAndroid
    OSTypeIOS
    OSTypeWeb
    OSTypeMax
    // OSTypeMobile = OSTypeIOS | OSTypeAndroid
    // OSTypePC = OSTypeLinux | OSTypeMacOS | OSTypeWindows
)

// GOOSType -
func GOOSType() OSType {
    switch (runtime.GOOS) {
    case "linux": return OSTypeLinux
    case "darwin": return OSTypeMacOS
    case "windows": return OSTypeWindows
    case "android": return OSTypeAndroid
    case "ios": return OSTypeIOS
    case "web": return OSTypeWeb
    }
    return OSTypeNone
}

// String -
func (s OSType) String() string {
    switch (s) {
    case OSTypeLinux: return "Linux"
    case OSTypeMacOS: return "MacOS"
    case OSTypeWindows: return "Windows"
    case OSTypeAndroid: return "Android"
    case OSTypeIOS: return "IOS"
    case OSTypeWeb: return "Web"
    }
    return "None"
}

// Match -
func (s OSType) Match(v OSType) bool {
    return s & v == v
}

// Mobile -
func (s OSType) Mobile() bool {
    return s == OSTypeAndroid || s == OSTypeIOS
}

// PC -
func (s OSType) Desktop() bool {
    return s == OSTypeLinux || s == OSTypeMacOS || s == OSTypeWindows
}


var AppIds = map[OSType]string{
    OSTypeWindows: "2071EFE494A482739FC5622E193F4CCB",
    OSTypeAndroid: "AC0A6CAED62777635C12FEEF4664CCBC",
    OSTypeMacOS: "4B63BBA2AA48FEB0E597BE4FEA86CD6D",
    OSTypeIOS: "CDFC62201C14C73EDD34AB0ADC7BD94B",
    OSTypeLinux: "E97D516D8D27160E77D300B7F54FBE5B",
}

// CheckAppID -
func CheckAppID(t OSType, s string) bool {
    if t <= OSTypeNone || t >= OSTypeMax {
        return false
    }
    if i, ok := AppIds[t]; ok && i == s {
        return true
    }
    return false
}

// AppID -
func AppID() string {
    return AppIds[GOOSType()]
}
