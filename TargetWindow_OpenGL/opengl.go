package TargetWindow_OpenGL

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Target struct {
	WhiteBox bool
	vao      uint32
	program  uint32
	window   *glfw.Window
	IsActive bool
}

const title = "TargetWindow OpenGL"

var (
	vertices = []float32{
		-1, 0.8, 0, // Oben Links
		-1, 1, 0,
		-0.8, 0.8, 0,
		-0.8, 1, 0, // Oben Links
		-1, 1, 0,
		-0.8, 0.8, 0,

		1, 0.8, 0, // Oben Rechts
		1, 1, 0,
		0.8, 0.8, 0,
		0.8, 1, 0, // Oben Rechts
		1, 1, 0,
		0.8, 0.8, 0,

		1, -0.8, 0, // Unten Rechts
		1, -1, 0,
		0.8, -0.8, 0,
		0.8, -1, 0, // Unten Rechts
		1, -1, 0,
		0.8, -0.8, 0,

		-1, -0.8, 0, // Unten Links
		-1, -1, 0,
		-0.8, -0.8, 0,
		-0.8, -1, 0, // Unten Links
		-1, -1, 0,
		-0.8, -0.8, 0,
	}

	fragmentShaderSource = `
	#version 410
	out vec4 frag_colour;
	void main() {
		frag_colour = vec4(1, 1, 1, 1);
	}
` + "\x00"
)

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(fullscreen bool, width, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var window *glfw.Window
	var err error
	if fullscreen { // fullscreen
		monitor := glfw.GetPrimaryMonitor()
		mode := monitor.GetVideoMode()
		window, err = glfw.CreateWindow(mode.Width, mode.Height, title, monitor, nil)
	} else {
		window, err = glfw.CreateWindow(width, height, title, nil, nil)
	}
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	// glfw.SwapInterval(1)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true)

		case glfw.KeyF11:
			// https://gist.github.com/pwaller/73593ae93d4f252bfb85
		}
	})
	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	// log.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()

	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func (t *Target) Start(fullscreen bool, fWidth, fHeight int, pushMode bool) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if fWidth == 0 {
		fWidth = 640
	}
	if fHeight == 0 {
		fHeight = 480
	}

	t.window = initGlfw(fullscreen, fWidth, fHeight)
	defer glfw.Terminate()

	t.program = initOpenGL()
	t.vao = makeVao(vertices)

	if !pushMode {
		go func(t *Target) {
			for !t.window.ShouldClose() {
				do(func() {
					t.draw()
				})
			}
		}(t)
	}

	t.IsActive = true
	for fn := range mainfunc {
		fn()
		if t.window.ShouldClose() {
			break
		}
	}

	t.Close()
}

func (t *Target) draw() {
	if !t.IsActive {
		return
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	if t.WhiteBox {
		gl.UseProgram(t.program)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)/3))
	}

	glfw.PollEvents()
	t.window.SwapBuffers()
}

func (t *Target) SetWhite() {
	do(func() {
		t.WhiteBox = true
		t.draw()
	})
}
func (t *Target) SetBlack() {
	do(func() {
		t.WhiteBox = false
		t.draw()
	})
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func (t *Target) Close() {
	t.IsActive = false
	glfw.Terminate()
}

// func getGID() uint64 {
// 	b := make([]byte, 64)
// 	b = b[:runtime.Stack(b, false)]
// 	b = bytes.TrimPrefix(b, []byte("goroutine "))
// 	b = b[:bytes.IndexByte(b, ' ')]
// 	n, _ := strconv.ParseUint(string(b), 10, 64)
// 	return n
// }

// queue of work to run in main thread.
var mainfunc = make(chan func())

// do runs f on the main thread.
func do(f func()) {
	done := make(chan bool, 1)
	mainfunc <- func() {
		f()
		done <- true
	}
	<-done

	// mainfunc <- func() {
	// 	f()
	// }
}
