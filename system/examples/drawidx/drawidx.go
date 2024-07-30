package main

import (
	"embed"
	"fmt"
	"image"
	"time"
	"unsafe"

	"github.com/ddkwork/golibrary/mylog"
	vk "github.com/goki/vulkan"

	"cogentcore.org/core/events"
	"cogentcore.org/core/math32"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/system"
	_ "cogentcore.org/core/system/driver"
	"cogentcore.org/core/vgpu"
)

//go:embed *.spv
var content embed.FS

type CamView struct {
	Model      math32.Matrix4
	View       math32.Matrix4
	Projection math32.Matrix4
}

func main() {
	opts := &system.NewWindowOptions{
		Size:      image.Pt(1024, 768),
		StdPixels: true,
		Title:     "System Test Window",
	}
	w := mylog.Check2(system.TheApp.NewWindow(opts))

	var sf *vgpu.Surface
	var sy *vgpu.System
	var pl *vgpu.Pipeline
	var cam *vgpu.Value
	var camo CamView

	make := func() {
		sf = w.Drawer().Surface().(*vgpu.Surface)
		sy = sf.GPU.NewGraphicsSystem("drawidx", &sf.Device)

		destroy := func() {
			sy.Destroy()
		}
		w.SetDestroyGPUResourcesFunc(destroy)

		pl = sy.NewPipeline("drawidx")

		sy.ConfigRender(&sf.Format, vgpu.Depth32)
		sf.SetRender(&sy.Render)
		sy.SetClearColor(0.2, 0.2, 0.2, 1)
		sy.SetRasterization(vk.PolygonModeFill, vk.CullModeNone, vk.FrontFaceCounterClockwise, 1.0)

		pl.AddShaderEmbed("indexed", vgpu.VertexShader, content, "indexed.spv")
		pl.AddShaderEmbed("vtxcolor", vgpu.FragmentShader, content, "vtxcolor.spv")

		vars := sy.Vars()
		vset := vars.AddVertexSet()
		set := vars.AddSet()

		nPts := 3

		posv := vset.Add("Pos", vgpu.Float32Vector3, nPts, vgpu.Vertex, vgpu.VertexShader)
		clrv := vset.Add("Color", vgpu.Float32Vector3, nPts, vgpu.Vertex, vgpu.VertexShader)

		idxv := vset.Add("Index", vgpu.Uint16, nPts, vgpu.Index, vgpu.VertexShader)

		camv := set.Add("Camera", vgpu.Struct, 1, vgpu.Uniform, vgpu.VertexShader)
		camv.SizeOf = vgpu.Float32Matrix4.Bytes() * 3

		vset.ConfigValues(1)
		set.ConfigValues(1)
		sy.Config()

		triPos, _ := posv.Values.ValueByIndexTry(0)
		triPosA := triPos.Floats32()
		triPosA.Set(0,
			-0.5, 0.5, 0.0,
			0.5, 0.5, 0.0,
			0.0, -0.5, 0.0)
		triPos.SetMod()

		triClr, _ := clrv.Values.ValueByIndexTry(0)
		triClrA := triClr.Floats32()
		triClrA.Set(0,
			1.0, 0.0, 0.0,
			0.0, 1.0, 0.0,
			0.0, 0.0, 1.0)
		triClr.SetMod()

		triIndex, _ := idxv.Values.ValueByIndexTry(0)
		idxs := []uint16{0, 1, 2}
		triIndex.CopyFromBytes(unsafe.Pointer(&idxs[0]))

		cam, _ = camv.Values.ValueByIndexTry(0)
		campos := math32.Vec3(0, 0, 2)
		target := math32.Vec3(0, 0, 0)
		var lookq math32.Quat
		lookq.SetFromRotationMatrix(math32.NewLookAt(campos, target, math32.Vec3(0, 1, 0)))
		scale := math32.Vec3(1, 1, 1)
		var cview math32.Matrix4
		cview.SetTransform(campos, lookq, scale)
		view, _ := cview.Inverse()

		camo.Model.SetIdentity()
		camo.View.CopyFrom(view)
		aspect := float32(sf.Format.Size.X) / float32(sf.Format.Size.Y)
		fmt.Printf("aspect: %g\n", aspect)

		camo.Projection.SetVkPerspective(45, aspect, 0.01, 100)

		cam.CopyFromBytes(unsafe.Pointer(&camo))

		sy.Mem.SyncToGPU()

		vars.BindDynamicValue(0, camv, cam)
	}

	frameCount := 0
	stTime := time.Now()

	renderFrame := func() {
		if sf == nil {
			make()
		}

		camo.Model.SetRotationY(.1 * float32(frameCount))
		cam.CopyFromBytes(unsafe.Pointer(&camo))
		sy.Mem.SyncToGPU()

		idx, ok := sf.AcquireNextImage()
		if !ok {
			return
		}

		descIndex := 0
		cmd := sy.CmdPool.Buff
		sy.ResetBeginRenderPass(cmd, sf.Frames[idx], descIndex)
		pl.BindDrawVertex(cmd, descIndex)
		sy.EndRenderPass(cmd)

		sf.SubmitRender(cmd)

		sf.PresentImage(idx)

		frameCount++
		eTime := time.Now()
		dur := float64(eTime.Sub(stTime)) / float64(time.Second)
		if dur > 10 {
			fps := float64(frameCount) / dur
			fmt.Printf("fps: %.0f\n", fps)
			frameCount = 0
			stTime = eTime
		}
	}

	haveKeyboard := false
	go func() {
		for {
			evi := w.Events().Deque.NextEvent()
			et := evi.Type()
			if et != events.WindowPaint && et != events.MouseMove {
				fmt.Println("got event", evi)
			}
			switch et {
			case events.Window:
				ev := evi.(*events.WindowEvent)
				fmt.Println("got window event", ev)
				switch ev.Action {
				case events.WinShow:
					make()
				case events.WinClose:
					fmt.Println("got events.Close; quitting")
					system.TheApp.Quit()
				}
			case events.WindowPaint:
				if w.IsVisible() {
					renderFrame()
				} else {
					fmt.Println("skipping paint event")
				}
			case events.MouseDown:
				if haveKeyboard {
					system.TheApp.HideVirtualKeyboard()
				} else {
					system.TheApp.ShowVirtualKeyboard(styles.KeyboardMultiLine)
				}
				haveKeyboard = !haveKeyboard
			}
		}
	}()
	system.TheApp.MainLoop()
}
