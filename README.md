# go-hes
#### Part of [HES-V2](https://github.com/HackerSchool/HES-V2)
HackerSchool HES driver written in Go  
[Original driver with config utility](https://github.com/HackerSchool/HES)

## Install and Update
```
$ got get -u github.com/hackerschool/go-hes
```

## Running go-hes
Linux and Mac users:
```
sudo $GOPATH/bin/go-hes
```
Without sudo, the driver will fail to send keys.

Windows users:
```
Double click executable in %GOPATH%/bin
```
### Config
Linux and Mac users:
```
$GOPATH/bin/go-hes config
```
Windows Users
```
%GOPATH%\bin\go-hes config
```
## New features
- Pressing select 5 consecutive times closes the driver 
- HTML based config utility
 
## Multiple gamepads
To use multiple HES you should connect all of the gamepads desired and run the driver

#### Valid keys
```go
		"a":     keybd.VK_A,
		"b":     keybd.VK_B,
		"c":     keybd.VK_C,
		"d":     keybd.VK_D,
		"e":     keybd.VK_E,
		"f":     keybd.VK_F,
		"g":     keybd.VK_G,
		"h":     keybd.VK_H,
		"i":     keybd.VK_I,
		"j":     keybd.VK_J,
		"k":     keybd.VK_K,
		"l":     keybd.VK_L,
		"m":     keybd.VK_M,
		"n":     keybd.VK_N,
		"o":     keybd.VK_O,
		"p":     keybd.VK_P,
		"q":     keybd.VK_Q,
		"r":     keybd.VK_R,
		"s":     keybd.VK_S,
		"t":     keybd.VK_T,
		"u":     keybd.VK_U,
		"v":     keybd.VK_V,
		"w":     keybd.VK_W,
		"x":     keybd.VK_X,
		"y":     keybd.VK_Y,
		"z":     keybd.VK_Z,
		"0":     keybd.VK_0,
		"1":     keybd.VK_1,
		"2":     keybd.VK_2,
		"3":     keybd.VK_3,
		"4":     keybd.VK_4,
		"5":     keybd.VK_5,
		"6":     keybd.VK_6,
		"7":     keybd.VK_7,
		"8":     keybd.VK_8,
		"9":     keybd.VK_9,
		"up":    keybd.VK_UP,
		"down":  keybd.VK_DOWN,
		"left":  keybd.VK_LEFT,
		"right": keybd.VK_RIGHT,
		"esc":   keybd.VK_ESC,
		"space": keybd.VK_ENTER,
```
