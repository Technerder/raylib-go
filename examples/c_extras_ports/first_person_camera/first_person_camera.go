/*
	The entirety of this file is a direct port of https://github.com/raylib-extras/extras-c/tree/main/cameras/rlFPCamera
*/

package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"math"
)

type FirstPersonCameraControls int32

const (
	MoveFront FirstPersonCameraControls = iota
	MoveBack
	MoveRight
	MoveLeft
	MoveUp
	MoveDown
	TurnLeft
	TurnRight
	TurnUp
	TurnDown
	Sprint
	LastControl
)

type FirstPersonCamera struct {
	ControlKeys              []int32
	MoveSpeed                rl.Vector3
	TurnSpeed                rl.Vector2
	UseMouse                 bool
	MouseSensitivity         float32
	MinimumViewY             float32
	MaximumViewY             float32
	ViewBobbleFreq           float32
	ViewBobbleMagnitude      float32
	ViewBobbleWaverMagnitude float32
	CameraPosition           rl.Vector3
	PlayerEyesPosition       float32
	FOV                      rl.Vector2
	TargetDistance           float32
	ViewAngles               rl.Vector2
	CurrentBobble            float32
	Focused                  bool
	AllowFlight              bool
	ViewCamera               rl.Camera3D
	Forward                  rl.Vector3
	Right                    rl.Vector3
	NearPlane                float32
	FarPlane                 float32
}

func (camera *FirstPersonCamera) Init(fovY float32, position rl.Vector3) {
	camera.ControlKeys = make([]int32, LastControl)
	camera.ControlKeys[0] = rl.KeyW
	camera.ControlKeys[1] = rl.KeyS
	camera.ControlKeys[2] = rl.KeyD
	camera.ControlKeys[3] = rl.KeyA
	camera.ControlKeys[4] = rl.KeyE
	camera.ControlKeys[5] = rl.KeyQ
	camera.ControlKeys[6] = rl.KeyLeft
	camera.ControlKeys[7] = rl.KeyRight
	camera.ControlKeys[8] = rl.KeyUp
	camera.ControlKeys[9] = rl.KeyDown
	camera.ControlKeys[10] = rl.KeyLeftShift

	camera.MoveSpeed = rl.NewVector3(1, 1, 1)
	camera.TurnSpeed = rl.NewVector2(90, 90)

	camera.UseMouse = true
	camera.MouseSensitivity = 600

	camera.MinimumViewY = -89.0
	camera.MaximumViewY = 89.0

	camera.ViewBobbleFreq = 0.0
	camera.ViewBobbleMagnitude = 0.02
	camera.ViewBobbleWaverMagnitude = 0.002
	camera.CurrentBobble = 0

	camera.Focused = rl.IsWindowFocused()

	camera.TargetDistance = 1
	camera.PlayerEyesPosition = 0.5
	camera.ViewAngles = rl.NewVector2(0, 0)

	camera.CameraPosition = position
	camera.FOV.Y = fovY

	camera.ViewCamera.Position = position
	camera.ViewCamera.Position.Y += camera.PlayerEyesPosition
	camera.ViewCamera.Target = rl.Vector3Add(camera.ViewCamera.Position, rl.NewVector3(0, 0, camera.TargetDistance))
	camera.ViewCamera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.ViewCamera.Fovy = fovY
	camera.ViewCamera.Projection = rl.CameraPerspective

	camera.AllowFlight = false
	camera.NearPlane = 0.01
	camera.FarPlane = 1000.0

	camera.ResizeView()
	camera.SetUseMouse(camera.UseMouse)
}

func (camera *FirstPersonCamera) SetUseMouse(useMouse bool) {
	camera.UseMouse = useMouse
	if useMouse && rl.IsWindowFocused() {
		rl.DisableCursor()
	} else {
		rl.EnableCursor()
	}
}

func (camera *FirstPersonCamera) ResizeView() {
	width := float32(rl.GetScreenWidth())
	height := float32(rl.GetScreenHeight())
	camera.FOV.Y = camera.ViewCamera.Fovy
	if height != 0 {
		camera.FOV.X = camera.FOV.Y * (width / height)
	}
}

func (camera *FirstPersonCamera) GetPosition() rl.Vector3 {
	return camera.CameraPosition
}

func (camera *FirstPersonCamera) SetPosition(pos rl.Vector3) {
	camera.CameraPosition = pos
	forward := rl.Vector3Subtract(camera.ViewCamera.Target, camera.ViewCamera.Position)
	camera.ViewCamera.Position = camera.CameraPosition
	camera.ViewCamera.Target = rl.Vector3Add(camera.CameraPosition, forward)
}

func (camera *FirstPersonCamera) GetViewRay() rl.Ray {
	return rl.NewRay(camera.CameraPosition, camera.Forward)
}

func GetSpeedForAxis(camera *FirstPersonCamera, axis FirstPersonCameraControls, speed float32) float32 {
	key := camera.ControlKeys[axis]
	if key == -1 {
		return 0
	}
	var factor float32
	factor = 1.0
	if rl.IsKeyDown(camera.ControlKeys[Sprint]) {
		factor = 2
	}
	if rl.IsKeyDown(camera.ControlKeys[axis]) {
		return speed * rl.GetFrameTime() * factor
	}
	return 0.0
}

func (camera *FirstPersonCamera) Update() {
	if rl.IsWindowFocused() != camera.Focused && camera.UseMouse {
		camera.Focused = rl.IsWindowFocused()
		if camera.Focused {
			rl.DisableCursor()
		} else {
			rl.EnableCursor()
		}
	}
	mousePositionDelta := rl.GetMouseDelta()
	direction := make([]float32, MoveDown+1)
	direction[0] = GetSpeedForAxis(camera, MoveFront, camera.MoveSpeed.Z)
	direction[1] = GetSpeedForAxis(camera, MoveBack, camera.MoveSpeed.Z)
	direction[2] = GetSpeedForAxis(camera, MoveRight, camera.MoveSpeed.X)
	direction[3] = GetSpeedForAxis(camera, MoveLeft, camera.MoveSpeed.X)
	direction[4] = GetSpeedForAxis(camera, MoveUp, camera.MoveSpeed.Y)
	direction[5] = GetSpeedForAxis(camera, MoveDown, camera.MoveSpeed.Y)

	turnRotation := GetSpeedForAxis(camera, TurnRight, camera.TurnSpeed.X) - GetSpeedForAxis(camera, TurnLeft, camera.TurnSpeed.X)
	tiltRotation := GetSpeedForAxis(camera, TurnUp, camera.TurnSpeed.Y) - GetSpeedForAxis(camera, TurnDown, camera.TurnSpeed.Y)

	if turnRotation != 0 {
		camera.ViewAngles.X -= turnRotation * rl.Deg2rad
	} else {
		camera.ViewAngles.X += mousePositionDelta.X / -camera.MouseSensitivity
	}

	if tiltRotation != 0 {
		camera.ViewAngles.Y += tiltRotation * rl.Deg2rad
	} else if camera.UseMouse && camera.Focused {
		camera.ViewAngles.Y += mousePositionDelta.Y / -camera.MouseSensitivity
	}

	if camera.ViewAngles.Y < camera.MinimumViewY*rl.Deg2rad {
		camera.ViewAngles.Y = camera.MinimumViewY * rl.Deg2rad
	} else if camera.ViewAngles.Y > camera.MaximumViewY*rl.Deg2rad {
		camera.ViewAngles.Y = camera.MaximumViewY * rl.Deg2rad
	}

	target := rl.Vector3Transform(rl.NewVector3(0, 0, 1), rl.MatrixRotateXYZ(rl.NewVector3(camera.ViewAngles.Y, -camera.ViewAngles.X, 0)))

	if camera.AllowFlight {
		camera.Forward = target
	} else {
		camera.Forward = rl.Vector3Transform(rl.NewVector3(0, 0, 1), rl.MatrixRotateXYZ(rl.NewVector3(0, -camera.ViewAngles.X, 0)))
	}

	camera.Right = rl.NewVector3(camera.Forward.Z*-1.0, 0, camera.Forward.X)

	camera.CameraPosition = rl.Vector3Add(camera.CameraPosition, rl.Vector3Scale(camera.Forward, direction[MoveFront]-direction[MoveBack]))
	camera.CameraPosition = rl.Vector3Add(camera.CameraPosition, rl.Vector3Scale(camera.Right, direction[MoveRight]-direction[MoveLeft]))

	camera.CameraPosition.Y += direction[MoveUp] - direction[MoveDown]
	camera.ViewCamera.Position = camera.CameraPosition

	eyeOffset := camera.PlayerEyesPosition

	if camera.ViewBobbleFreq > 0 {
		swingDelta := float32(math.Max(math.Abs(float64(direction[MoveFront]-direction[MoveBack])), math.Abs(float64(direction[MoveRight]-direction[MoveLeft]))))
		camera.CurrentBobble += swingDelta * camera.ViewBobbleFreq
		viewBobbleDampen := float32(8.0)
		eyeOffset -= float32(math.Sin(float64(camera.CurrentBobble/viewBobbleDampen)) * float64(camera.ViewBobbleMagnitude))
		camera.ViewCamera.Up.X = float32(math.Sin(float64(camera.CurrentBobble/(viewBobbleDampen*2)))) * camera.ViewBobbleWaverMagnitude
		camera.ViewCamera.Up.Z = -float32(math.Sin(float64(camera.CurrentBobble/(viewBobbleDampen*2)))) * camera.ViewBobbleWaverMagnitude
	} else {
		camera.CurrentBobble = 0
		camera.ViewCamera.Up.X = 0
		camera.ViewCamera.Up.Z = 0
	}

	camera.ViewCamera.Position.Y += eyeOffset

	camera.ViewCamera.Target.X = camera.ViewCamera.Position.X + target.X
	camera.ViewCamera.Target.Y = camera.ViewCamera.Position.Y + target.Y
	camera.ViewCamera.Target.Z = camera.ViewCamera.Position.Z + target.Z
}

func (camera *FirstPersonCamera) SetupCamera(aspect float32) {
	rl.DrawRenderBatchActive()
	rl.MatrixMode(rl.RL_PROJECTION)
	rl.PushMatrix()
	rl.LoadIdentity()
	if camera.ViewCamera.Projection == rl.CameraPerspective {
		top := float32(rl.RL_CULL_DISTANCE_NEAR * math.Tan(float64(camera.ViewCamera.Fovy*0.5*rl.Deg2rad)))
		right := top * aspect
		rl.Frustum(-right, right, -top, top, camera.NearPlane, camera.FarPlane)
	} else if camera.ViewCamera.Projection == rl.CameraOrthographic {
		top := camera.ViewCamera.Fovy / 2.0
		right := top * aspect
		rl.Ortho(-right, right, -top, top, camera.NearPlane, camera.FarPlane)
	}
	rl.MatrixMode(rl.RL_MODELVIEW)
	rl.LoadIdentity()
	matView := rl.MatrixLookAt(camera.ViewCamera.Position, camera.ViewCamera.Target, camera.ViewCamera.Up)
	rl.MultMatrixf(rl.MatrixToFloatV(matView).V)
	rl.EnableDepthTest()
}

func (camera *FirstPersonCamera) BeginMode3D() {
	aspect := rl.GetScreenWidth() / rl.GetScreenHeight()
	camera.SetupCamera(float32(aspect))
}

func (camera *FirstPersonCamera) EndMode3D() {
	rl.EndMode3D()
}
