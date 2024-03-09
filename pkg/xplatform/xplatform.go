/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

xplatform.go

General platform function
*/
package xplatform

import "runtime"

func IsWindows() bool {
	return runtime.GOOS == "windows"
}
