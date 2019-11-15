package sparkle

import (
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"golang.org/x/exp/shiny/driver/gldriver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

type Context struct {
	Buffer *image.RGBA
	GL     gl.Context
	Size   size.Event

	touchX, touchY float32
	screen         screen.Screen
	texture        screen.Texture
	windowBuffer   screen.Buffer
	window         screen.Window
	images         *glutil.Images
	fps            *debug.FPS
	drawFuncs      map[uint]func(*Context)
	drawFuncMut    sync.RWMutex
	drawNum        uint
}

func (c *Context) GetTouchX() float32 {
	return c.touchX
}

func (c *Context) GetTouchY() float32 {
	return c.touchY
}

func (c *Context) AddDrawer(draw func(*Context)) uint {
	uid := c.drawNum
	c.drawNum++

	c.drawFuncMut.Lock()
	c.drawFuncs[uid] = draw
	c.drawFuncMut.Unlock()

	return uid
}

func (c *Context) RemoveDrawer(uid uint) {
	delete(c.drawFuncs, uid)
}

func setupContext() *Context {
	return &Context{
		drawFuncs: make(map[uint]func(*Context)),
	}
}

func LoadShader(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bVal, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(bVal), nil
}

func Main(opts *screen.NewWindowOptions, start, stop, update func(*Context)) (*Context, error) {
	var (
		ctx = setupContext()
		err error
	)

	gldriver.Main(func(s screen.Screen) {
		ctx.screen = s

		ctx.window, err = s.NewWindow(opts)
		if err != nil {
			err = fmt.Errorf("creating new window: %w", err)
			return
		}

		ctx.texture, err = s.NewTexture(image.Point{opts.Width, opts.Height})
		if err != nil {
			return
		}

		// tx, _ := s.NewTexture(image.Point{800, 600})
		ctx.windowBuffer, err = s.NewBuffer(image.Point{opts.Width, opts.Height})
		if err != nil {
			return
		}

		ctx.Buffer = ctx.windowBuffer.RGBA()

		// go eventLoop(ctx, start, stop, update)
		eventLoop(ctx, start, stop, update)
	})

	if err != nil {
		return nil, fmt.Errorf("sparkle main: %w", err)
	}

	return ctx, nil
}

func eventLoop(ctx *Context, start, stop, update func(*Context)) {
	quit := make(chan struct{})
	var err error

	for {
		switch e := ctx.window.NextEvent().(type) {
		case lifecycle.Event:
			switch e.Crosses(lifecycle.StageVisible) {
			case lifecycle.CrossOn:
				ctx.GL, _ = e.DrawContext.(gl.Context)

				sparkleStart(ctx)
				start(ctx)

				go sparkleUpdate(ctx, quit, update)
				go sparkleDraw(ctx, quit)
			case lifecycle.CrossOff:
				close(quit)
				sparkleStop(ctx)
				stop(ctx)
				ctx.GL = nil
				return
			}
		case size.Event:
			ctx.Size = e
			ctx.touchX = float32(ctx.Size.WidthPx / 2)
			ctx.touchY = float32(ctx.Size.HeightPx / 2)

			ctx.windowBuffer.Release()
			ctx.windowBuffer, err = ctx.screen.NewBuffer(image.Point{
				e.WidthPx, e.HeightPx,
			})
			if err != nil {
				panic(err)
			}

			ctx.Buffer = ctx.windowBuffer.RGBA()
		case key.Event:
			if e.Code == key.CodeEscape {
				ctx.window.Send(lifecycle.Event{
					From: lifecycle.StageVisible,
					To:   lifecycle.StageDead,
				})
			}
		}
	}
}

func sparkleStart(ctx *Context) {
	log.Printf("Initializing Window\n")

	ctx.images = glutil.NewImages(ctx.GL)
	ctx.fps = debug.NewFPS(ctx.images)
}

func sparkleStop(ctx *Context) {
	log.Printf("Killing Window\n")

	ctx.fps.Release()
	ctx.images.Release()
}

func sparkleUpdate(ctx *Context, quit chan struct{}, update func(*Context)) {
	for {
		select {
		case <-quit:
			return
		default:
		}

		update(ctx)
	}
}

func newBuf(ctx *Context) {
	bounds := ctx.Buffer.Bounds().Max
	ctx.Buffer = nil
	ctx.windowBuffer.Release()
	ctx.windowBuffer, _ = ctx.screen.NewBuffer(bounds)
	ctx.Buffer = ctx.windowBuffer.RGBA()
}

func sparkleDraw(ctx *Context, quit chan struct{}) {
	for {
		select {
		case <-quit:
			return
		default:
		}

		ctx.GL.ClearColor(0, 0, 0, 1)
		ctx.GL.Clear(gl.COLOR_BUFFER_BIT)

		newBuf(ctx)

		for i := range ctx.drawFuncs {
			ctx.drawFuncs[i](ctx)
		}

		ctx.texture.Upload(
			image.Point{0, 0},
			ctx.windowBuffer,
			ctx.windowBuffer.Bounds(),
		)
		ctx.window.Scale(
			ctx.windowBuffer.Bounds(),
			ctx.texture,
			ctx.texture.Bounds(),
			draw.Over, // Draw over existing items.
			// draw.Src, // Omit previously existing items.
			&screen.DrawOptions{},
		)

		ctx.fps.Draw(ctx.Size)
		ctx.window.Publish()
	}
}
