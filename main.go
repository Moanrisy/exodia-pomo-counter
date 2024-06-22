package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	// Set to allow transparent window
	FlagWindowTransparent = 0x00000010
	// Set to disable window decoration (frame and buttons)
	FlagWindowUndecorated = 0x00000008
	// Set to support mouse passthrough, only supported when FLAG_WINDOW_UNDECORATED
	FlagWindowMousePassthrough = 0x00004000
)

var bkgColor = rl.NewColor(0, 0, 0, 0)

func checkErr(err error) {
	if err != nil {
		// log.Fatal(err)

		// mylog := log.New(os.Stderr, "app: ", log.LstdFlags|log.Lshortfile)
		// fmt.Println(mylog)

		pc, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
		log.Printf("Hey you have an cathed error")
	}
}

func main() {
	rl.SetConfigFlags(FlagWindowTransparent)
	rl.SetConfigFlags(FlagWindowUndecorated)
	// rl.SetConfigFlags(FlagWindowMousePassthrough)

	rl.InitWindow(117, 130, "exodia-pomo-counter")
	defer rl.CloseWindow()

	// rl.SetExitKey(0)

	// rl.SetTargetFPS(60)
	rl.SetTargetFPS(1)

	var headImg rl.Texture2D = rl.LoadTexture("./.config/exodia-pomo-counter-res-assets/lower_res/head.png")
	var leftHandImg rl.Texture2D = rl.LoadTexture("./.config/exodia-pomo-counter-res-assets/lower_res/left_hand.png")
	var rightHandImg rl.Texture2D = rl.LoadTexture("./.config/exodia-pomo-counter-res-assets/lower_res/right_hand.png")
	var leftFootImg rl.Texture2D = rl.LoadTexture("./.config/exodia-pomo-counter-res-assets/lower_res/left_foot.png")
	var rightFootImg rl.Texture2D = rl.LoadTexture("./.config/exodia-pomo-counter-res-assets/lower_res/right_foot.png")

	isFirstRun := true
	pomodoroCounter := 0
	mod4PomoCounter := 0

	rl.InitAudioDevice()
	// music := rl.LoadMusicStream("./.config/exodia-pomo-counter-res-assets/exodia_obliterate.mp3")
	music := rl.LoadSound("./.config/exodia-pomo-counter-res-assets/exodia_obliterate.mp3")
	// music := rl.LoadSound("./.config/exodia-pomo-counter-res-assets/harrys-avery-farm.mp3")
	// rl.StopMusicStream(music)
	// rl.PlayMusicStream(music)
	// rl.ResumeMusicStream(music)
	// musicPaused := true

	homeDir, _ := os.UserHomeDir()
	filePath := homeDir + "/.variables/test"
	// Check if the path file exist
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("File doesn't exist, creating directory and file: %v", err)

		// Create the directory
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		checkErr(err)

		// Create the file
		file, err = os.Create(filePath)
		checkErr(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	initialStat, err := os.Stat(filePath)
	checkErr(err)

	playAudioOnce := false
	isPaused := false

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		if !isPaused {
			// rl.ClearBackground(rl.RayWhite)
			rl.ClearBackground(bkgColor)

			rl.DrawTextureEx(headImg, rl.NewVector2(float32(40), float32(8)), 0, 1, rl.White)
			switch pomodoroCounter {
			case 0:
				if playAudioOnce {
					rl.PlaySound(music)
					playAudioOnce = false
				}
			case 1:
				rl.DrawTextureEx(leftHandImg, rl.NewVector2(float32(0), float32(15)), 0, 1, rl.White)

			case 2:
				rl.DrawTextureEx(leftHandImg, rl.NewVector2(float32(0), float32(15)), 0, 1, rl.White)
				rl.DrawTextureEx(rightHandImg, rl.NewVector2(float32(80), float32(15)), 0, 1, rl.White)
			case 3:
				rl.DrawTextureEx(leftHandImg, rl.NewVector2(float32(0), float32(15)), 0, 1, rl.White)
				rl.DrawTextureEx(rightHandImg, rl.NewVector2(float32(80), float32(15)), 0, 1, rl.White)
				rl.DrawTextureEx(leftFootImg, rl.NewVector2(float32(20), float32(70)), 0, 1, rl.White)
				playAudioOnce = true
			}
			// rl.DrawTextureEx(rightFootImg, rl.NewVector2(float32(60), float32(70)), 0, 1, rl.White)

			go func() {
				for {
					stat, err := os.Stat(filePath)
					checkErr(err)
					if isFirstRun || stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
						if isFirstRun {
							isFirstRun = false
						}

						buf := make([]byte, stat.Size())
						_, err := file.Read(buf)

						bufString := string(buf)
						// print only the last line
						// because at first run it always return all the text content
						lines := strings.Split(bufString, "\n")
						// lastLine := ""
						if len(lines) > 2 {
							// lastLine = lines[len(lines)-2]
							mod4PomoCounter = (len(lines)) % 4
						} else if err != io.EOF {
							// lastLine = bufString
							mod4PomoCounter = (mod4PomoCounter + 1) % 4
							pomodoroCounter = mod4PomoCounter
						}

						initialStat, err = os.Stat(filePath)
						checkErr(err)
					}
					time.Sleep(1 * time.Second)
				}
			}()

			// isPaused = true
		}
		time.Sleep(1 * time.Second)
		rl.EndDrawing()
	}

	rl.UnloadTexture(headImg)
	rl.UnloadTexture(leftHandImg)
	rl.UnloadTexture(rightHandImg)
	rl.UnloadTexture(leftFootImg)
	rl.UnloadTexture(rightFootImg)

	// rl.UnloadMusicStream(music)
	rl.UnloadSound(music)
	rl.CloseAudioDevice()
}
