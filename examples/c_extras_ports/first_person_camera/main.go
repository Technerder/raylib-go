/*
	The entirety of this file is a direct port of https://github.com/raylib-extras/extras-c/blob/main/cameras/rlFPCamera/samples/example.c
*/

package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialization
	//--------------------------------------------------------------------------------------
	screenWidth := int32(1900)
	screenHeight := int32(900)

	rl.SetConfigFlags(rl.FlagFullscreenMode | rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "raylib-c_extras_ports [camera] example - First person camera")
	rl.SetTargetFPS(144)

	//--------------------------------------------------------------------------------------
	img := rl.GenImageChecked(256, 256, 32, 32, rl.DarkGray, rl.White)
	tx := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	rl.SetTextureFilter(tx, rl.FilterAnisotropic16x)
	rl.SetTextureWrap(tx, rl.WrapClamp)

	// setup initial camera data
	camera := &FirstPersonCamera{}
	camera.Init(45, rl.NewVector3(1, 0, 0))
	camera.MoveSpeed.Z = 10
	camera.MoveSpeed.X = 5

	camera.FarPlane = 5000

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		if rl.IsKeyPressed(rl.KeyF1) {
			camera.AllowFlight = !camera.AllowFlight
		}

		camera.Update()
		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)

		camera.BeginMode3D()

		// grid of cube trees on a plane to make a "world"
		rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(50, 50), rl.Beige) // simple world plane

		spacing := 4
		count := 5

		for x := -count * spacing; x <= count*spacing; x += spacing {
			for z := -count * spacing; z <= count*spacing; z += spacing {
				rl.DrawCubeTexture(tx, rl.NewVector3(float32(x), 1.5, float32(z)), 1, 1, 1, rl.Green)
				rl.DrawCubeTexture(tx, rl.NewVector3(float32(x), 0.5, float32(z)), 0.25, 1, 0.25, rl.Brown)
			}
		}

		camera.EndMode3D()

		if camera.AllowFlight {
			rl.DrawText("(F1) Flight", 2, 20, 20, rl.Black)
		} else {
			rl.DrawText("(F1) Running", 2, 20, 20, rl.Black)
		}
		// instructions
		rl.DrawFPS(0, 0)
		rl.EndDrawing()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.CloseWindow() // Close window and OpenGL context
	//--------------------------------------------------------------------------------------
}
