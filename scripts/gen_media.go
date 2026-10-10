package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
)

// Terminal font definition (simple 7x13 bitmap font for clean terminal look)
// We'll define a crisp font renderer so every text looks sharp and terminal-native.

var font5x7 = map[rune][]string{
	' ': {"     ", "     ", "     ", "     ", "     ", "     ", "     "},
	'A': {" ### ", "#   #", "#   #", "#####", "#   #", "#   #", "#   #"},
	'B': {"#### ", "#   #", "#### ", "#   #", "#   #", "#   #", "#### "},
	'C': {" ####", "#    ", "#    ", "#    ", "#    ", "#    ", " ####"},
	'D': {"#### ", "#   #", "#   #", "#   #", "#   #", "#   #", "#### "},
	'E': {"#####", "#    ", "#### ", "#    ", "#    ", "#    ", "#####"},
	'F': {"#####", "#    ", "#### ", "#    ", "#    ", "#    ", "#    "},
	'G': {" ####", "#    ", "#    ", "# ###", "#   #", "#   #", " ### "},
	'H': {"#   #", "#   #", "#####", "#   #", "#   #", "#   #", "#   #"},
	'I': {"#####", "  #  ", "  #  ", "  #  ", "  #  ", "  #  ", "#####"},
	'J': {"  ###", "    #", "    #", "    #", "    #", "#   #", " ### "},
	'K': {"#   #", "#  # ", "###  ", "#  # ", "#   #", "#   #", "#   #"},
	'L': {"#    ", "#    ", "#    ", "#    ", "#    ", "#    ", "#####"},
	'M': {"#   #", "## ##", "# # #", "#   #", "#   #", "#   #", "#   #"},
	'N': {"#   #", "##  #", "# # #", "#  ##", "#   #", "#   #", "#   #"},
	'O': {" ### ", "#   #", "#   #", "#   #", "#   #", "#   #", " ### "},
	'P': {"#### ", "#   #", "#### ", "#    ", "#    ", "#    ", "#    "},
	'Q': {" ### ", "#   #", "#   #", "#   #", "# # #", "#  ##", " ####"},
	'R': {"#### ", "#   #", "#### ", "#  # ", "#   #", "#   #", "#   #"},
	'S': {" ####", "#    ", " ### ", "    #", "    #", "    #", "#### "},
	'T': {"#####", "  #  ", "  #  ", "  #  ", "  #  ", "  #  ", "  #  "},
	'U': {"#   #", "#   #", "#   #", "#   #", "#   #", "#   #", " ### "},
	'V': {"#   #", "#   #", "#   #", " # # ", " # # ", "  #  ", "  #  "},
	'W': {"#   #", "#   #", "# # #", "# # #", "# # #", "## ##", "#   #"},
	'X': {"#   #", " # # ", "  #  ", "  #  ", " # # ", "#   #", "#   #"},
	'Y': {"#   #", " # # ", "  #  ", "  #  ", "  #  ", "  #  ", "  #  "},
	'Z': {"#####", "    #", "   # ", "  #  ", " #   ", "#    ", "#####"},
	'a': {"     ", " ### ", "    #", " ####", "#   #", "#   #", " ####"},
	'b': {"#    ", "#    ", "#### ", "#   #", "#   #", "#   #", "#### "},
	'c': {"     ", " ####", "#    ", "#    ", "#    ", "#    ", " ####"},
	'd': {"    #", "    #", " ####", "#   #", "#   #", "#   #", " ####"},
	'e': {"     ", " ### ", "#   #", "#####", "#    ", "#    ", " ####"},
	'f': {"  ## ", " #  #", " #   ", "#### ", " #   ", " #   ", " #   "},
	'g': {"     ", " ####", "#   #", " ####", "    #", "#   #", " ### "},
	'h': {"#    ", "#    ", "#### ", "#   #", "#   #", "#   #", "#   #"},
	'i': {"  #  ", "     ", " ##  ", "  #  ", "  #  ", "  #  ", " ### "},
	'j': {"    #", "     ", "   ##", "    #", "    #", "#   #", " ### "},
	'k': {"#    ", "#    ", "#  # ", "###  ", "#  # ", "#  # ", "#   #"},
	'l': {" ##  ", "  #  ", "  #  ", "  #  ", "  #  ", "  #  ", " ### "},
	'm': {"     ", "## ##", "# # #", "#   #", "#   #", "#   #", "#   #"},
	'n': {"     ", "#### ", "#   #", "#   #", "#   #", "#   #", "#   #"},
	'o': {"     ", " ### ", "#   #", "#   #", "#   #", "#   #", " ### "},
	'p': {"     ", "#### ", "#   #", "#### ", "#    ", "#    ", "#    "},
	'q': {"     ", " ####", "#   #", " ####", "    #", "    #", "    #"},
	'r': {"     ", " # ##", " ##  ", " #   ", " #   ", " #   ", " #   "},
	's': {"     ", " ####", "#    ", " ### ", "    #", "    #", "#### "},
	't': {"  #  ", " ### ", "  #  ", "  #  ", "  #  ", "  #  ", "   ##"},
	'u': {"     ", "#   #", "#   #", "#   #", "#   #", "#   #", " ####"},
	'v': {"     ", "#   #", "#   #", " # # ", " # # ", "  #  ", "  #  "},
	'w': {"     ", "#   #", "# # #", "# # #", "# # #", "## ##", "#   #"},
	'x': {"     ", "#   #", " # # ", "  #  ", " # # ", "#   #", "#   #"},
	'y': {"     ", "#   #", "#   #", " ####", "    #", " #  #", "  ## "},
	'z': {"     ", "#####", "   # ", "  #  ", " #   ", "#    ", "#####"},
	'0': {" ### ", "#  ##", "# # #", "# # #", "##  #", "#   #", " ### "},
	'1': {"  #  ", " ##  ", "  #  ", "  #  ", "  #  ", "  #  ", " ### "},
	'2': {" ### ", "#   #", "    #", "  ## ", " #   ", "#    ", "#####"},
	'3': {"#### ", "    #", "  ## ", "    #", "    #", "    #", "#### "},
	'4': {"#  # ", "#  # ", "#  # ", "#####", "   # ", "   # ", "   # "},
	'5': {"#####", "#    ", "#### ", "    #", "    #", "    #", "#### "},
	'6': {" ### ", "#    ", "#### ", "#   #", "#   #", "#   #", " ### "},
	'7': {"#####", "    #", "   # ", "  #  ", "  #  ", "  #  ", "  #  "},
	'8': {" ### ", "#   #", " ### ", "#   #", "#   #", "#   #", " ### "},
	'9': {" ### ", "#   #", " ####", "    #", "    #", "    #", " ### "},
	':': {"     ", "  #  ", "     ", "     ", "  #  ", "     ", "     "},
	'.': {"     ", "     ", "     ", "     ", "     ", "  #  ", "  #  "},
	',': {"     ", "     ", "     ", "     ", "     ", "  #  ", " #   "},
	'-': {"     ", "     ", "     ", "#####", "     ", "     ", "     "},
	'_': {"     ", "     ", "     ", "     ", "     ", "     ", "#####"},
	'/': {"    #", "   # ", "   # ", "  #  ", " #   ", " #   ", "#    "},
	'\\': {"#    ", " #   ", " #   ", "  #  ", "   # ", "   # ", "    #"},
	'[': {" ### ", " #   ", " #   ", " #   ", " #   ", " #   ", " ### "},
	']': {" ### ", "   # ", "   # ", "   # ", "   # ", "   # ", " ### "},
	'(': {"  ## ", " #   ", " #   ", " #   ", " #   ", " #   ", "  ## "},
	')': {" ##  ", "   # ", "   # ", "   # ", "   # ", "   # ", " ##  "},
	'|': {"  #  ", "  #  ", "  #  ", "  #  ", "  #  ", "  #  ", "  #  "},
	'%': {"##  #", "## # ", "  #  ", " #   ", " # ##", "#  ##", "#   #"},
	'#': {" # # ", "#####", " # # ", " # # ", "#####", " # # ", " # # "},
	'*': {"  #  ", "# # #", " ### ", "  #  ", " ### ", "# # #", "  #  "},
	'+': {"     ", "  #  ", "  #  ", "#####", "  #  ", "  #  ", "     "},
	'=': {"     ", "#####", "     ", "#####", "     ", "     ", "     "},
	'<': {"   # ", "  #  ", " #   ", "#    ", " #   ", "  #  ", "   # "},
	'>': {" #   ", "  #  ", "   # ", "    #", "   # ", "  #  ", " #   "},
	'!': {"  #  ", "  #  ", "  #  ", "  #  ", "     ", "  #  ", "  #  "},
	'?': {" ### ", "#   #", "    #", "  ## ", "  #  ", "     ", "  #  "},
	'~': {"  ## ", " #  #", "#    ", "     ", "     ", "     ", "     "},
	'@': {" ### ", "#   #", " ## #", "# # #", "# ###", "#    ", " ### "},
}

func drawChar(img *image.RGBA, x, y int, char rune, col color.RGBA, scale int) int {
	matrix, ok := font5x7[char]
	if !ok {
		matrix = font5x7['?']
	}
	charWidth := 5 * scale
	for rowIdx, row := range matrix {
		for colIdx, pixel := range row {
			if pixel == '#' {
				for dx := 0; dx < scale; dx++ {
					for dy := 0; dy < scale; dy++ {
						px := x + colIdx*scale + dx
						py := y + rowIdx*scale + dy
						if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
							img.Set(px, py, col)
						}
					}
				}
			}
		}
	}
	return charWidth + (1 * scale)
}

func drawText(img *image.RGBA, x, y int, text string, col color.RGBA, scale int) {
	curX := x
	for _, ch := range text {
		curX += drawChar(img, curX, y, ch, col, scale)
	}
}

func drawFilledRect(img *image.RGBA, x1, y1, x2, y2 int, col color.RGBA) {
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
				img.Set(x, y, col)
			}
		}
	}
}

func drawRoundedRect(img *image.RGBA, x1, y1, x2, y2, r int, col color.RGBA) {
	drawFilledRect(img, x1+r, y1, x2-r, y2, col)
	drawFilledRect(img, x1, y1+r, x2, y2-r, col)
}

func drawCircle(img *image.RGBA, cx, cy, r int, col color.RGBA) {
	for y := cy - r; y <= cy + r; y++ {
		for x := cx - r; x <= cx + r; x++ {
			if (x-cx)*(x-cx)+(y-cy)*(y-cy) <= r*r {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, col)
				}
			}
		}
	}
}

// Color palette
var (
	cBg       = color.RGBA{18, 20, 32, 255}    // #121420
	cWindowBg = color.RGBA{26, 27, 42, 255}    // #1a1b2a
	cTitleBar = color.RGBA{36, 40, 59, 255}    // #24283b
	cCardBg   = color.RGBA{31, 35, 53, 255}    // #1f2335
	cText     = color.RGBA{192, 202, 245, 255} // #c0caf5
	cMuted    = color.RGBA{100, 110, 150, 255} // #646e96
	cPrimary  = color.RGBA{122, 162, 247, 255} // #7aa2f7
	cCyan     = color.RGBA{125, 207, 255, 255} // #7dcfff
	cPurple   = color.RGBA{187, 154, 247, 255} // #bb9af7
	cGreen    = color.RGBA{158, 206, 106, 255} // #9ece6a
	cYellow   = color.RGBA{224, 175, 104, 255} // #e0af68
	cRed      = color.RGBA{247, 118, 142, 255} // #f7768e
	cWhite    = color.RGBA{255, 255, 255, 255}
	cHighlight= color.RGBA{65, 72, 104, 255}   // #414868
)

func renderWindowBase(width, height int, title string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Outer background
	drawFilledRect(img, 0, 0, width, height, cBg)

	// Window container
	winX1, winY1 := 20, 20
	winX2, winY2 := width-20, height-20
	drawRoundedRect(img, winX1, winY1, winX2, winY2, 8, cWindowBg)

	// Title bar
	drawFilledRect(img, winX1, winY1, winX2, winY1+36, cTitleBar)

	// Window buttons (red, yellow, green)
	drawCircle(img, winX1+20, winY1+18, 6, cRed)
	drawCircle(img, winX1+38, winY1+18, 6, cYellow)
	drawCircle(img, winX1+56, winY1+18, 6, cGreen)

	// Window title
	drawText(img, winX1+80, winY1+14, title, cMuted, 1)

	return img
}

func renderTerminalTabs(img *image.RGBA, activeTab int, startX, startY int) {
	tabs := []string{"[ 1 ] Ports", "[ 2 ] Processes", "[ 3 ] System Vitals", "[ 4 ] Dev Cleaner"}
	curX := startX
	for i, t := range tabs {
		tWidth := len(t)*6*2 + 16
		if i == activeTab {
			drawFilledRect(img, curX, startY, curX+tWidth, startY+24, cPrimary)
			drawText(img, curX+8, startY+6, t, cWhite, 2)
		} else {
			drawFilledRect(img, curX, startY, curX+tWidth, startY+24, cCardBg)
			drawText(img, curX+8, startY+6, t, cMuted, 2)
		}
		curX += tWidth + 8
	}
}

func renderFrame(activeTab int, stage int) *image.RGBA {
	width, height := 860, 520
	img := renderWindowBase(width, height, "devpulse -- developer cockpit (zsh)")

	// Logo & Header
	drawFilledRect(img, 40, 75, 190, 105, cCardBg)
	drawText(img, 50, 83, "DEVPULSE", cCyan, 2)
	drawText(img, 205, 87, "v1.0.0 | ALL-IN-ONE TERMINAL COCKPIT", cMuted, 2)

	// Tab bar
	renderTerminalTabs(img, activeTab, 40, 120)

	// Main card area
	drawFilledRect(img, 40, 155, width-40, height-65, cCardBg)

	switch activeTab {
	case 0: // Ports Tab
		filter := ""
		if stage >= 1 {
			filter = "node"
		}
		drawText(img, 55, 170, "Filter: ["+filter+"]  (Press '/' to filter, 'k' to Kill, 'r' to Refresh)", cMuted, 1)

		headers := fmt.Sprintf("%-8s %-6s %-22s %-8s %-14s %-8s", "PORT", "PROTO", "PROCESS", "PID", "BIND IP", "MEMORY")
		drawFilledRect(img, 50, 190, width-50, 212, cTitleBar)
		drawText(img, 55, 196, headers, cPurple, 1)

		ports := []struct {
			Port, Proto, Proc, Pid, IP, Mem string
		}{
			{"3000", "TCP", "node.exe (Next.js)", "34912", "127.0.0.1", "184 MB"},
			{"5000", "TCP", "python.exe (Flask)", "18440", "0.0.0.0", "92 MB"},
			{"5432", "TCP", "postgres.exe", "4102", "127.0.0.1", "240 MB"},
			{"8080", "TCP", "go-backend", "12904", "0.0.0.0", "45 MB"},
			{"27017", "TCP", "mongod", "8890", "127.0.0.1", "310 MB"},
			{"5173", "TCP", "vite.exe", "28376", "127.0.0.1", "120 MB"},
		}

		if stage >= 1 {
			// filtered
			ports = []struct {
				Port, Proto, Proc, Pid, IP, Mem string
			}{
				{"3000", "TCP", "node.exe (Next.js)", "34912", "127.0.0.1", "184 MB"},
			}
		}

		curY := 220
		for idx, p := range ports {
			rowStr := fmt.Sprintf(":%-7s %-6s %-22s %-8s %-14s %-8s", p.Port, p.Proto, p.Proc, p.Pid, p.IP, p.Mem)
			if idx == 0 {
				drawFilledRect(img, 50, curY-4, width-50, curY+16, cHighlight)
				drawText(img, 55, curY, "> "+rowStr, cWhite, 1)
			} else {
				drawText(img, 55, curY, "  "+rowStr, cText, 1)
			}
			curY += 24
		}

		if stage >= 2 {
			// Toast notification: Terminated!
			drawFilledRect(img, width-380, 75, width-50, 105, cPurple)
			drawText(img, width-370, 83, "Terminated PID 34912 (:3000)", cWhite, 2)
		}

	case 1: // Processes Tab
		drawText(img, 55, 170, "Sorted by: [CPU] (press 'c' CPU, 'm' Mem, 'n' Name) | 'k' to Kill", cMuted, 1)
		headers := fmt.Sprintf("%-8s %-24s %-8s %-8s %-12s %-20s", "PID", "PROCESS NAME", "CPU %", "MEM %", "RSS", "COMMAND")
		drawFilledRect(img, 50, 190, width-50, 212, cTitleBar)
		drawText(img, 55, 196, headers, cPurple, 1)

		procs := []struct {
			Pid, Name, Cpu, Mem, Rss, Cmd string
		}{
			{"18440", "python.exe", "18.4%", "3.2%", "92.0 MB", "python server.py --port 5000"},
			{"31916", "node.exe", "12.1%", "6.5%", "184.0 MB", "npm run dev (Next.js SSR)"},
			{"26156", "docker.exe", "8.7%", "8.2%", "420.0 MB", "dockerd daemon"},
			{"12904", "devpulse.exe", "1.2%", "0.4%", "14.2 MB", "devpulse interactive"},
			{"4102", "postgres.exe", "0.9%", "4.1%", "240.0 MB", "postgres -D /data"},
		}

		curY := 220
		for idx, p := range procs {
			rowStr := fmt.Sprintf("%-8s %-24s %-8s %-8s %-12s %-20s", p.Pid, p.Name, p.Cpu, p.Mem, p.Rss, p.Cmd)
			if idx == 0 {
				drawFilledRect(img, 50, curY-4, width-50, curY+16, cHighlight)
				drawText(img, 55, curY, "> "+rowStr, cWhite, 1)
			} else {
				drawText(img, 55, curY, "  "+rowStr, cText, 1)
			}
			curY += 24
		}

	case 2: // System Vitals
		drawText(img, 55, 170, "Host: dev-workstation | OS: Windows 11 / x86_64 | Uptime: 4d 12h", cCyan, 2)

		drawText(img, 55, 205, "CPU LOAD & SPARKLINE", cPrimary, 2)
		drawText(img, 55, 230, "Load: [====================..........] 64.2%", cGreen, 2)
		drawText(img, 55, 255, "Trend: _.-^^---..___..-^^^---..._^^^-", cPurple, 2)

		drawText(img, 55, 290, "RAM & SWAP USAGE", cPrimary, 2)
		drawText(img, 55, 315, "RAM:  [===============...............] 48.1% (15.4 GB / 32.0 GB)", cYellow, 2)
		drawText(img, 55, 340, "Swap: [======........................] 19.4% (3.1 GB / 16.0 GB)", cMuted, 2)

		drawText(img, 55, 375, "NETWORK THROUGHPUT & DISKS", cPrimary, 2)
		drawText(img, 55, 400, "Net: Download: 1.42 MB/s  |  Upload: 284.1 KB/s", cGreen, 2)
		drawText(img, 55, 425, "Disk C: [====================........] 68.4% (680 GB / 1000 GB)", cText, 2)

	case 3: // Dev Cleaner
		drawText(img, 55, 170, "Total Reclaimable: 3.42 GB | Selected: 2.10 GB  ([d] Clean Selected)", cYellow, 2)

		headers := fmt.Sprintf("%-6s %-24s %-12s %-30s", "SEL", "CATEGORY", "SIZE", "DIRECTORY")
		drawFilledRect(img, 50, 195, width-50, 217, cTitleBar)
		drawText(img, 55, 201, headers, cPurple, 1)

		cleanItems := []struct {
			Sel, Cat, Size, Path string
		}{
			{"[X]", "Node Dependencies", "1.42 GB", "frontend/node_modules"},
			{"[X]", "Next.js Cache", "680.0 MB", "frontend/.next"},
			{"[ ]", "Rust Build Artifacts", "890.0 MB", "core/target"},
			{"[ ]", "Python Virtualenv", "430.0 MB", "backend/.venv"},
		}

		curY := 225
		for idx, it := range cleanItems {
			rowStr := fmt.Sprintf("%-6s %-24s %-12s %-30s", it.Sel, it.Cat, it.Size, it.Path)
			if idx == 0 {
				drawFilledRect(img, 50, curY-4, width-50, curY+16, cHighlight)
				drawText(img, 55, curY, "> "+rowStr, cWhite, 1)
			} else {
				drawText(img, 55, curY, "  "+rowStr, cText, 1)
			}
			curY += 24
		}
	}

	// Footer bar
	drawText(img, 55, height-40, "[Tab] Next Tab  |  [/] Filter  |  [k] Kill  |  [r] Refresh  |  [?] Help  |  [q] Quit", cMuted, 1)

	return img
}

func rgbaToPaletted(rgba *image.RGBA) *image.Paletted {
	pal := palette.Plan9
	bounds := rgba.Bounds()
	paletted := image.NewPaletted(bounds, pal)
	draw.Draw(paletted, bounds, rgba, bounds.Min, draw.Src)
	return paletted
}

func main() {
	fmt.Println("Generating visual assets for DevPulse...")

	// 1. Generate high-res Preview Image (PNG)
	previewImg := renderFrame(0, 0)
	fPreview, err := os.Create("assets/preview.png")
	if err == nil {
		png.Encode(fPreview, previewImg)
		fPreview.Close()
		fmt.Println("✓ Generated assets/preview.png")
	}

	// 2. Generate System Vitals Image (PNG)
	sysImg := renderFrame(2, 0)
	fSys, err := os.Create("assets/vitals.png")
	if err == nil {
		png.Encode(fSys, sysImg)
		fSys.Close()
		fmt.Println("✓ Generated assets/vitals.png")
	}

	// 3. Generate Hero Banner (PNG)
	bannerWidth, bannerHeight := 1200, 480
	banner := image.NewRGBA(image.Rect(0, 0, bannerWidth, bannerHeight))
	drawFilledRect(banner, 0, 0, bannerWidth, bannerHeight, cBg)
	// Grid accent lines
	for x := 0; x < bannerWidth; x += 40 {
		drawFilledRect(banner, x, 0, x+1, bannerHeight, color.RGBA{24, 28, 45, 255})
	}
	for y := 0; y < bannerHeight; y += 40 {
		drawFilledRect(banner, 0, y, bannerWidth, y+1, color.RGBA{24, 28, 45, 255})
	}

	// Banner Centerpiece Box
	drawRoundedRect(banner, 60, 40, bannerWidth-60, bannerHeight-40, 12, cWindowBg)

	// Neon Header
	drawText(banner, 100, 90, "DEVPULSE", cCyan, 5)
	drawText(banner, 100, 160, "The Blazingly Fast Developer Cockpit & Process TUI", cText, 2)
	drawText(banner, 100, 195, "Built with Go + Bubble Tea + Lip Gloss | Instant Port Killer | System Vitals", cMuted, 2)

	// Feature Badges
	badges := []struct {
		Text string
		Col  color.RGBA
	}{
		{"1-Key Port Kill (:3000)", cRed},
		{"Process Inspector", cPurple},
		{"CPU & RAM Sparklines", cGreen},
		{"Dev Cache Cleaner", cYellow},
		{"Zero Dependencies", cCyan},
	}
	bx := 100
	for _, b := range badges {
		bw := len(b.Text)*6*2 + 24
		drawFilledRect(banner, bx, 240, bx+bw, 275, cCardBg)
		drawText(banner, bx+12, 250, b.Text, b.Col, 2)
		bx += bw + 15
	}

	// Code snippet preview
	drawFilledRect(banner, 100, 310, bannerWidth-100, 390, cTitleBar)
	drawText(banner, 120, 325, "$ go install github.com/MisterPugs/devpulse/cmd/devpulse@latest", cGreen, 2)
	drawText(banner, 120, 355, "$ devpulse kill 3000   # Terminate process occupying port 3000 in 10ms", cText, 2)

	fBanner, err := os.Create("assets/banner.png")
	if err == nil {
		png.Encode(fBanner, banner)
		fBanner.Close()
		fmt.Println("✓ Generated assets/banner.png")
	}

	// 4. Generate Animated GIF demo.gif
	// Sequence of stages:
	// Frame 0: Ports tab initial (1200ms)
	// Frame 1: Filtering "node" (1000ms)
	// Frame 2: Killed PID toast (1200ms)
	// Frame 3: Processes tab (1400ms)
	// Frame 4: System Vitals tab with sparklines (1600ms)
	// Frame 5: Dev Cleaner tab (1400ms)

	stages := []struct {
		Tab   int
		Stage int
		Delay int // in 100ths of a second (e.g. 100 = 1s)
	}{
		{0, 0, 130},
		{0, 1, 100},
		{0, 2, 130},
		{1, 0, 140},
		{2, 0, 160},
		{3, 0, 140},
	}

	anim := &gif.GIF{
		LoopCount: 0, // infinite loop
	}

	for _, s := range stages {
		fImg := renderFrame(s.Tab, s.Stage)
		palImg := rgbaToPaletted(fImg)
		anim.Image = append(anim.Image, palImg)
		anim.Delay = append(anim.Delay, s.Delay)
	}

	fGif, err := os.Create("assets/demo.gif")
	if err == nil {
		gif.EncodeAll(fGif, anim)
		fGif.Close()
		fmt.Println("✓ Generated assets/demo.gif")
	}

	fmt.Println("Assets generated successfully!")
}
