package rl

/*
#include "raylib.h"
#include "rlgl.h"
#include "raymath.h"
#include <stdlib.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

// RLGL Constants
var (
	RL_MODELVIEW          = int32(0x1700)
	RL_PROJECTION         = int32(0x1701)
	RL_CULL_DISTANCE_NEAR = 0.01
)

// DrawRenderBatchActive calls rlgl#rlDrawRenderBatchActive
func DrawRenderBatchActive() {
	C.rlDrawRenderBatchActive()
}

// MatrixMode calls rlgl#rlMatrixMode
func MatrixMode(mode int32) {
	cflag := (C.int)(mode)
	C.rlMatrixMode(cflag)
}

// PushMatrix calls rlgl#rlPushMatrix
func PushMatrix() {
	C.rlPushMatrix()
}

// LoadIdentity calls rlgl#rlLoadIdentity
func LoadIdentity() {
	C.rlLoadIdentity()
}

// EnableDepthTest calls rlgl#rlEnableDepthTest
func EnableDepthTest() {
	C.rlEnableDepthTest()
}

// Frustum calls rlgl#rlFrustum
func Frustum(left float32, right float32, bottom float32, top float32, znear float32, zfar float32) {
	Left := (C.double)(left)
	Right := (C.double)(right)
	Bottom := (C.double)(bottom)
	Top := (C.double)(top)
	ZNear := (C.double)(znear)
	ZFar := (C.double)(zfar)
	C.rlFrustum(Left, Right, Bottom, Top, ZNear, ZFar)
}

// Ortho calls rlgl#rlOrtho
func Ortho(left float32, right float32, bottom float32, top float32, znear float32, zfar float32) {
	Left := (C.double)(left)
	Right := (C.double)(right)
	Bottom := (C.double)(bottom)
	Top := (C.double)(top)
	ZNear := (C.double)(znear)
	ZFar := (C.double)(zfar)
	C.rlOrtho(Left, Right, Bottom, Top, ZNear, ZFar)
}

// MultMatrixf calls rlgl#rlMultMatrixf
func MultMatrixf(matf []float32) {
	ccount := (*C.float)(unsafe.Pointer(&matf[0]))
	C.rlMultMatrixf(ccount)
}

type Float16 struct {
	V []float32
}

// MatrixToFloatV returns Float16 object
func MatrixToFloatV(mat Matrix) Float16 {
	result := Float16{
		V: make([]float32, 16),
	}
	result.V[0] = mat.M0
	result.V[1] = mat.M1
	result.V[2] = mat.M2
	result.V[3] = mat.M3
	result.V[4] = mat.M4
	result.V[5] = mat.M5
	result.V[6] = mat.M6
	result.V[7] = mat.M7
	result.V[8] = mat.M8
	result.V[9] = mat.M9
	result.V[10] = mat.M10
	result.V[11] = mat.M11
	result.V[12] = mat.M12
	result.V[13] = mat.M13
	result.V[14] = mat.M14
	result.V[15] = mat.M15
	return result
}

// cptr returns C pointer
func (s *Shader) cptr() *C.Shader {
	return (*C.Shader)(unsafe.Pointer(s))
}

// LoadShader - Load a custom shader and bind default locations
func LoadShader(vsFileName string, fsFileName string) Shader {
	cvsFileName := C.CString(vsFileName)
	defer C.free(unsafe.Pointer(cvsFileName))

	cfsFileName := C.CString(fsFileName)
	defer C.free(unsafe.Pointer(cfsFileName))

	if vsFileName == "" {
		cvsFileName = nil
	}

	if fsFileName == "" {
		cfsFileName = nil
	}

	ret := C.LoadShader(cvsFileName, cfsFileName)
	v := newShaderFromPointer(unsafe.Pointer(&ret))

	return v
}

// LoadShaderFromMemory - Load shader from code strings and bind default locations
func LoadShaderFromMemory(vsCode string, fsCode string) Shader {
	cvsCode := C.CString(vsCode)
	defer C.free(unsafe.Pointer(cvsCode))

	cfsCode := C.CString(fsCode)
	defer C.free(unsafe.Pointer(cfsCode))

	if vsCode == "" {
		cvsCode = nil
	}

	if fsCode == "" {
		cfsCode = nil
	}

	ret := C.LoadShaderFromMemory(cvsCode, cfsCode)
	v := newShaderFromPointer(unsafe.Pointer(&ret))

	return v
}

// UnloadShader - Unload a custom shader from memory
func UnloadShader(shader Shader) {
	cshader := shader.cptr()
	C.UnloadShader(*cshader)
}

// GetShaderLocation - Get shader uniform location
func GetShaderLocation(shader Shader, uniformName string) int32 {
	cshader := shader.cptr()
	cuniformName := C.CString(uniformName)
	defer C.free(unsafe.Pointer(cuniformName))

	ret := C.GetShaderLocation(*cshader, cuniformName)
	v := (int32)(ret)
	return v
}

// GetShaderLocationAttrib - Get shader attribute location
func GetShaderLocationAttrib(shader Shader, attribName string) int32 {
	cshader := shader.cptr()
	cuniformName := C.CString(attribName)
	defer C.free(unsafe.Pointer(cuniformName))

	ret := C.GetShaderLocationAttrib(*cshader, cuniformName)
	v := (int32)(ret)
	return v
}

// SetShaderValue - Set shader uniform value (float)
func SetShaderValue(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	cvalue := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&value)).Data)
	cuniformType := (C.int)(uniformType)
	C.SetShaderValue(*cshader, clocIndex, cvalue, cuniformType)
}

// SetShaderValueV - Set shader uniform value (float)
func SetShaderValueV(shader Shader, locIndex int32, value []float32, uniformType ShaderUniformDataType, count int32) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	cvalue := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&value)).Data)
	cuniformType := (C.int)(uniformType)
	ccount := (C.int)(count)
	C.SetShaderValueV(*cshader, clocIndex, cvalue, cuniformType, ccount)
}

// SetShaderValueMatrix - Set shader uniform value (matrix 4x4)
func SetShaderValueMatrix(shader Shader, locIndex int32, mat Matrix) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	cmat := mat.cptr()
	C.SetShaderValueMatrix(*cshader, clocIndex, *cmat)
}

// SetShaderValueTexture - Set shader uniform value for texture (sampler2d)
func SetShaderValueTexture(shader Shader, locIndex int32, texture Texture2D) {
	cshader := shader.cptr()
	clocIndex := (C.int)(locIndex)
	ctexture := texture.cptr()
	C.SetShaderValueTexture(*cshader, clocIndex, *ctexture)
}

// SetMatrixProjection - Set a custom projection matrix (replaces internal projection matrix)
func SetMatrixProjection(proj Matrix) {
	cproj := proj.cptr()
	C.rlSetMatrixProjection(*cproj)
}

// SetMatrixModelview - Set a custom modelview matrix (replaces internal modelview matrix)
func SetMatrixModelview(view Matrix) {
	cview := view.cptr()
	C.rlSetMatrixModelview(*cview)
}

// BeginShaderMode - Begin custom shader drawing
func BeginShaderMode(shader Shader) {
	cshader := shader.cptr()
	C.BeginShaderMode(*cshader)
}

// EndShaderMode - End custom shader drawing (use default shader)
func EndShaderMode() {
	C.EndShaderMode()
}

// BeginBlendMode - Begin blending mode (alpha, additive, multiplied)
func BeginBlendMode(mode BlendMode) {
	cmode := (C.int)(mode)
	C.BeginBlendMode(cmode)
}

// EndBlendMode - End blending mode (reset to default: alpha blending)
func EndBlendMode() {
	C.EndBlendMode()
}
