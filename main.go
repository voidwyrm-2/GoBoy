package main

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	carts "github.com/voidwyrm-2/GoBoy/cartridges"
)

type Cartridge interface {
	Init(*map[string]int, *map[string]int) ([2]int, int32, string)
	// returned by Init:
	// the x and y of the pixel screen
	// the FPS
	// if the returned string isn't empty, exit with error

	Update(*[][]color.RGBA, *map[string]int, *map[string]int) (int, string)
	// returned by Update:
	// 0 means continue as normal,
	// 1 means stop and call exit, and
	// -1 means exit with an error

	//Draw([][]color.RGBA) ([][]color.RGBA, string)
	// returned by Draw:
	// if the returned string isn't empty, exit with error

	Exit() string
	// returned by Exit:
	// if the returned string isn't empty, exit with error
}

const windowX int32 = 800
const windowY int32 = 600

const pixelScale int = 10

var board [][]color.RGBA

var loadedCart = carts.SpaceInvaders{}

var alarms = make(map[string]int)

var ticks = make(map[string]int)

var loadBarTick int = 0
var loadBarPercent int = 0
var loadBarHeight int = 20

func main() {
	rl.InitWindow(windowX, windowY, "GoBoy")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	var gbThemeMusic = rl.LoadMusicStream("assets/GoBoyTheme.mp3")
	var gbThemeMusicPlayed = false
	var gbThemeMusicStopped = false
	//var gbThemeMusicLength = rl.GetMusicTimeLength(gbThemeMusic)
	var gbThemeMusicTimeRemaining float32

	ticks["Game"] = 0

	boardSize, fps, initErrStr := loadedCart.Init(&ticks, &alarms)
	if initErrStr != "" {
		fmt.Printf("UpdateError: %s\n", initErrStr)
		return
	}

	rl.SetTargetFPS(fps)

	for y := range boardSize[1] {
		var tempY []rl.Color
		board = append(board, tempY)
		for range boardSize[0] {
			tempX := rl.Black
			board[y] = append(board[y], tempX)
		}
	}

	holdingShift := false

	for { //!rl.WindowShouldClose() {

		if loadBarPercent < 400 {

			if !gbThemeMusicPlayed {
				rl.PlayMusicStream(gbThemeMusic)
				gbThemeMusicPlayed = true
			} else {
				if !gbThemeMusicStopped {
					timeRemaining := fmt.Sprintf("%f", gbThemeMusicTimeRemaining)
					if timeRemaining[len(timeRemaining)-3:] == "816" {
						rl.StopMusicStream(gbThemeMusic)
						gbThemeMusicStopped = true
					}

					rl.UpdateMusicStream(gbThemeMusic)
					gbThemeMusicTimeRemaining = rl.GetMusicTimePlayed(gbThemeMusic)
				}

				//fmt.Printf("%f\n", gbThemeMusicLength)
				//fmt.Printf("%f\n", gbThemeMusicTimeRemaining)
			}

			if rl.IsKeyPressed(rl.KeyEscape) {
				fmt.Println("game was manually quit during loading")
				return
			}

			if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
				holdingShift = true
			} else {
				holdingShift = false
			}

			if holdingShift && rl.IsKeyPressed(rl.KeyPeriod) {
				loadBarPercent = 398
			}

			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			if loadBarTick == 0 {
				loadBarPercent++
				//loadBarTick = 0
			} else {
				loadBarTick--
			}

			//rl.DrawText("Creeper, oh man", 316, 200, 20, rl.White)

			goboyX := 280
			rl.DrawText("GO", int32(goboyX), 200, 100, color.RGBA{114, 206, 221, 255})
			rl.DrawText("BOY", int32(goboyX+140), 200, 100, color.RGBA{114, 206, 221, 255})

			rl.DrawRectangle(150, (windowY/2)-int32(loadBarHeight), int32(loadBarPercent), int32(loadBarHeight), rl.White)

			rl.EndDrawing()
		} else {

			for tk := range ticks {
				ticks[tk]++
			}

			for al := range alarms {
				if alarms[al] > 0 {
					alarms[al]--
				}
			}

			for y := range len(board) {
				for x := range len(board[y]) {
					board[y][x] = rl.Black
				}
			}

			code, updateErrStr := loadedCart.Update(&board, &ticks, &alarms)
			if code == 1 {
				break
			} else if code == -1 {
				fmt.Printf("UpdateError: %s\n", updateErrStr)
				return
			} else if code != 0 {
				fmt.Printf("UpdateError: expected '0' but got '%d' instead\n", code)
				return
			}

			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			//rl.DrawText("Creeper, oh man", 316, 200, 20, rl.White)

			/*
				drawErrStr := loadedCart.Draw()
				if drawErrStr != "" {
					fmt.Printf("DrawError: %s\n", drawErrStr)
					return
				}
			*/

			for y := range len(board) {
				for x := range len(board[y]) {
					rl.DrawRectangle(int32(x*pixelScale), int32(y*pixelScale), int32(pixelScale), int32(pixelScale), board[y][x])
				}
			}

			rl.EndDrawing()
		}
	}

	exitErrStr := loadedCart.Exit()

	if exitErrStr != "" {
		fmt.Printf("ExitError: %s\n", exitErrStr)
		return
	}
}
