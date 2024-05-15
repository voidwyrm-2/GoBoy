package cartridges

import (
	"fmt"
	"image/color"
	"math/rand"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
	npg "github.com/voidwyrm-2/npglib"
)

func generateIDCode(prefix string, length int) string {
	var out string
	if prefix != "" {
		out += prefix + "-"
	}
	for range length {
		out += fmt.Sprint(rand.Intn(10))
	}

	return out
}

type Basic struct {
	plX       int
	plY       int
	boardSize [2]int
	plColor   color.RGBA
}

func (crt *Basic) Init(ticks *map[string]int, alarms *map[string]int) ([2]int, int32, string) {
	// init cardtridge variables
	crt.plX = 0
	crt.plY = 0
	crt.boardSize = [2]int{10, 10}
	crt.plColor = rl.Red

	// init external variables
	tempTicks := *ticks
	tempTicks["Basic"] = 0
	*ticks = tempTicks

	return crt.boardSize, 60, ""
}

func (crt *Basic) Update(board *[][]color.RGBA, ticks *map[string]int, alarms *map[string]int) (int, string) {
	// make indexable copy of the board
	mutBoard := *board

	if rl.IsKeyPressed(rl.KeyUp) {
		if crt.plY > 0 {
			crt.plY--
		}
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		if crt.plY < crt.boardSize[1]-1 {
			crt.plY++
		}
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		if crt.plX > 0 {
			crt.plX--
		}
	}

	if rl.IsKeyPressed(rl.KeyRight) {
		if crt.plX < crt.boardSize[0]-1 {
			crt.plX++
		}
	}

	// call self.draw()
	drawnBoard, errStr := crt.Draw(mutBoard)
	if errStr != "" {
		return -1, errStr
	}

	// apply changes to board
	*board = drawnBoard
	return 0, ""
}

func (crt Basic) Draw(mutableBoard [][]color.RGBA) ([][]color.RGBA, string) {
	mBoard := mutableBoard
	mBoard[crt.plY][crt.plX] = crt.plColor
	return mBoard, ""
}

func (crt *Basic) Exit() string {
	return ""
}

type SpaceInvaders struct {
	plX             int
	plY             int
	bullets         map[string][2]int
	enemies         map[string][3]int
	startingEnemies int
	boardSize       [2]int
	plColor         color.RGBA
	bulletColor     color.RGBA
	enemyColor      color.RGBA
	bulletSpeed     int
	enemySpeed      int
	exitKey         int
	hasLost         bool
	hasWon          bool
	lossSprite      npg.Sprite
	winSprite       npg.Sprite
}

func (crt *SpaceInvaders) Init(ticks *map[string]int, alarms *map[string]int) ([2]int, int32, string) {
	// init cardtridge variables
	crt.boardSize = [2]int{10, 10}
	crt.plX = 0
	crt.plY = crt.boardSize[1] - 1
	crt.startingEnemies = 5
	crt.plColor = color.RGBA{150, 150, 255, 255}
	crt.bulletColor = color.RGBA{155, 0, 0, 255}
	crt.enemyColor = color.RGBA{10, 255, 20, 255}
	crt.bulletSpeed = 6
	crt.enemySpeed = 8
	crt.exitKey = rl.KeyEscape
	crt.hasLost = false
	crt.lossSprite = npg.Sprite{Size: [2]int{10, 10}}
	crt.lossSprite.GenerateFromString("          \n"+
		//"          \n"+
		"  #       \n"+
		"  #       \n"+
		"  #       \n"+
		"  #       \n"+
		"  #       \n"+
		"  ####    ", 0)
	crt.winSprite = npg.Sprite{Size: [2]int{10, 10}}
	crt.winSprite.GenerateFromString("          \n"+
		"#       # \n"+
		"#       # \n"+
		" #  #  #  \n"+
		" # # # #  \n"+
		"  #   #   \n"+
		"          ", 0)

	crt.bullets = make(map[string][2]int)
	crt.enemies = make(map[string][3]int)

	eY := 0
	eX := 0
	for range crt.startingEnemies {
		if eX > crt.boardSize[0]-1 {
			eY++
			if eY%2 == 0 || eY == 0 {
				eX = 1
			} else {
				eX = 0
			}
		}
		if eY > crt.boardSize[1]-1 {
			break
		}
		crt.enemies[generateIDCode("e", 4)] = [3]int{eX, eY, 0}
		eX += 2
	}

	// create alarms
	tempAlarms := *alarms
	tempAlarms["BulletCooldown"] = 0
	tempAlarms["BulletClock"] = crt.bulletSpeed
	tempAlarms["EnemyClock"] = crt.enemySpeed
	*alarms = tempAlarms

	return crt.boardSize, 15, ""
}

func (crt *SpaceInvaders) Update(board *[][]color.RGBA, ticks *map[string]int, alarms *map[string]int) (int, string) {
	// make indexable copies of global variables
	mutBoard := *board
	mutAlarms := *alarms

	// add extra variables
	var slatedForRemoval []string

	if crt.hasWon {
		if rl.IsKeyPressed(int32(crt.exitKey)) {
			return 1, ""
		}

		if rl.IsKeyPressed(rl.KeyR) {
			crt.Init(ticks, alarms)
			crt.hasWon = false
		}
	} else if crt.hasLost {
		if rl.IsKeyPressed(int32(crt.exitKey)) {
			return 1, ""
		}

		if rl.IsKeyPressed(rl.KeyR) {
			crt.Init(ticks, alarms)
			crt.hasLost = false
		}
	} else {

		if rl.IsKeyPressed(int32(crt.exitKey)) {
			return 1, ""
		}

		// bullet update
		if mutAlarms["BulletClock"] == 0 {
			for i := range crt.bullets {
				if crt.bullets[i][1] > 0 {
					temp := crt.bullets[i] // make addressable copy
					temp[1]--
					crt.bullets[i] = temp // update original with copy
				}
			}

			mutAlarms["BulletClock"] = crt.bulletSpeed
		}

		if mutAlarms["EnemyClock"] == 0 {
			for i := range crt.enemies {
				if crt.enemies[i][1] < crt.boardSize[1]-1 {
					if crt.enemies[i][2] == 0 {
						if crt.enemies[i][0] < crt.boardSize[0]-1 {
							temp := crt.enemies[i] // make addressable copy
							temp[0]++
							crt.enemies[i] = temp // update original with copy
						} else {
							temp := crt.enemies[i] // make addressable copy
							temp[1]++              // move down
							//crt.enemies[i][0] = crt.boardSize[0] - 1  // reset postion(not needed)
							temp[2] = 1           // flip movement(goes left now)
							crt.enemies[i] = temp // update original with copy
						}
					} else if crt.enemies[i][2] == 1 {
						if crt.enemies[i][0] > 0 {
							temp := crt.enemies[i] // make addressable copy
							temp[0]--
							crt.enemies[i] = temp // update original with copy
						} else {
							temp := crt.enemies[i] // make addressable copy
							temp[1]++              // move down
							//crt.enemies[i][0] = 0 // reset postion(not needed)
							temp[2] = 0           // flip movement(goes right now)
							crt.enemies[i] = temp // update original with copy
						}
					} else {
						return -1, fmt.Sprintf("enemy id'%s' should have a direction of '0' or '1', but had '%d' instead", i, crt.enemies[i][2])
					}
				} else {
					//slatedForRemoval = append(slatedForRemoval, i)
					crt.hasLost = true
				}
			}

			mutAlarms["EnemyClock"] = crt.enemySpeed
		}

		// shoot
		if rl.IsKeyPressed(rl.KeyUp) {
			if mutAlarms["BulletCooldown"] == 0 {
				crt.bullets[generateIDCode("b", 4)] = [2]int{crt.plX, crt.plY - 1}
				mutAlarms["BulletCooldown"] = 30
			}
		}

		// flag
		for i, b := range crt.bullets {
			if b[1] == 0 {
				slatedForRemoval = append(slatedForRemoval, i)
			}
		}

		// bullet -> enemy collision
		for i1, e := range crt.enemies {
			for i2, b := range crt.bullets {
				if b[0] == e[0] && b[1] == e[1] {
					if !slices.Contains(slatedForRemoval, i2) {
						slatedForRemoval = append(slatedForRemoval, i2)
					}
					slatedForRemoval = append(slatedForRemoval, i1)
				}
			}
		}

		// BEGIN remove objects

		for _, id := range slatedForRemoval {
			_, b_ok := crt.bullets[id]
			if b_ok {
				delete(crt.bullets, id)
				continue
			}

			_, e_ok := crt.enemies[id]
			if e_ok {
				delete(crt.enemies, id)
				continue
			}
		}
		// END remove objects

		if len(crt.enemies) == 0 {
			crt.hasWon = true
		}

		// BEGIN movement controls
		if rl.IsKeyPressed(rl.KeyLeft) {
			if crt.plX > 0 {
				crt.plX--
			}
		}

		if rl.IsKeyPressed(rl.KeyRight) {
			if crt.plX < crt.boardSize[0]-1 {
				crt.plX++
			}
		}
		// END movement controls
	}

	// call self.draw()
	drawnBoard, errStr := crt.Draw(mutBoard)
	if errStr != "" {
		return -1, errStr
	}

	// apply changes to board
	*board = drawnBoard
	return 0, ""
}

func (crt SpaceInvaders) Draw(mutableBoard [][]color.RGBA) ([][]color.RGBA, string) {
	mBoard := mutableBoard

	if crt.hasWon {
		crt.winSprite.DrawSpriteOnBoard(0, 0, &mBoard)
	} else if crt.hasLost {
		crt.lossSprite.DrawSpriteOnBoard(0, 0, &mBoard)
	} else {
		// draw player
		mBoard[crt.plY][crt.plX] = crt.plColor

		// draw bullets
		for _, b := range crt.bullets {
			mBoard[b[1]][b[0]] = crt.bulletColor
		}

		// draw enemies
		for _, e := range crt.enemies {
			mBoard[e[1]][e[0]] = crt.enemyColor
		}
	}

	return mBoard, ""
}

func (crt *SpaceInvaders) Exit() string {
	return ""
}
