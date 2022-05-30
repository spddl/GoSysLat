module GoSysLat

go 1.17

replace github.com/spddl/RTSSClient => ./RTSSClient

replace github.com/spddl/USBController => ./USBController

replace github.com/spddl/TargetWindow_D3D9 => ./TargetWindow_D3D9

replace github.com/spddl/TargetWindow_OpenGL => ./TargetWindow_OpenGL

require (
	github.com/MakeNowJust/hotkey v0.0.0-20200628032113-41fa0caa507a
	github.com/VividCortex/ewma v1.2.0
	github.com/eclesh/welford v0.0.0-20150116075914-eec62615b1f0
	github.com/go-echarts/go-echarts/v2 v2.2.4
	github.com/lxn/walk v0.0.0-20210112085537-c389da54e794
	github.com/spddl/RTSSClient v0.0.0-00010101000000-000000000000
	github.com/spddl/TargetWindow_D3D9 v0.0.0-00010101000000-000000000000
	github.com/spddl/TargetWindow_OpenGL v0.0.0-00010101000000-000000000000
	github.com/spddl/USBController v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a
)

require (
	github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6 // indirect
	github.com/go-gl/glfw v0.0.0-20220320163800-277f93cfa958 // indirect
	github.com/gonutz/d3d9 v1.2.1 // indirect
	github.com/gonutz/w32/v2 v2.4.0 // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
)
