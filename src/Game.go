var frame string                         // game frames
var dir string                           // current directory
var events chan string                   // keyboard events
var gameOver = false                     // end of game
var windowWidth, windowHeight = 400, 300 // width and height of the window
var frameRate int                        // how many frames to show per second (fps)
var gameDelay int                        // delay time added to each game loop

func init() {
	// events is a channel of string events that come from the front end
	events = make(chan string, 1000)
	// getting the current directory to access resources
	var err error
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	frameRate = 50                                         // 50 fps
	gameDelay = 20                                         // 20 ms delay
	sprites = getImage(dir + "/box.jpg") // spritesheet
	background = getImage(dir + "/circle.jpg")   // background image
}

// main function
func main() {
	// run the web server in a separate goroutine
	go app()
	// create a web view
	err := webview.Open("Game", "/public/html/index.html",
	windowWidth, windowHeight, false)
	if err != nil {
		log.Fatal(err)
	}

	for !gameOver {
		// to slow up or speed up the game
		time.Sleep(time.Millisecond * time.Duration(gameDelay))
		// if any of the keyboard events are captured
		select {
		case ev := <-events:
			// exit the game
			if ev == "38" { // Up arrow key?
				player.Position.Y += 10
			}
			if ev == "40" { // Down arrow key?
				player.Position.Y -= 10
			}
			if ev == "39" { // right arrow key
				player.Position.X += 10
			}
			if ev == "37" { // left arrow key
				player.Position.X -= 10
			}
		}
	}
}
func app() {
	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(dir+"/public"))))
	mux.HandleFunc("/start", start)
	mux.HandleFunc("/frame", getFrame)
	mux.HandleFunc("/key", captureKeys)
	server := &http.Server{
		Addr:    "localhost:9090",
		Handler: mux,
	}
	server.ListenAndServe()
}
// start the game
func start(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(dir + "/public/game.html")
	// start generating frames in a new goroutine
	go generateFrames()
	t.Execute(w, 1000/frameRate)
}
// capture keyboard events
func captureKeys(w http.ResponseWriter, r *http.Request) {
	ev := r.FormValue("event")
	// what to react to when the game is over
	if gameOver {
		if ev == "83" { // s
			gameOver = false
			go generateFrames()
		}
		if ev == "81" { // q
			os.Exit(0)
		}

	} else {
		events <- ev
	}
	w.Header().Set("Cache-Control", "no-cache")
}
// get the game frames
func getFrame(w http.ResponseWriter, r *http.Request) {
	str := "data:image/png;base64," + frame
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(str))
}
