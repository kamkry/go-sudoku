package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"image"
	"log"
	"strconv"
)

//+gap
const cellHeight = 5
const cellWidth = 3
const emptyCell = "  "

type Cell struct {
	*widgets.Paragraph
	Value    int
	Editable bool
}

func (c Cell) String() string {
	return fmt.Sprint(c.Value)
}

var values *[9][9]int
var cells [9][9]Cell

var info *widgets.Paragraph
var lvlTab *widgets.TabPane
var solveBtn *widgets.Paragraph
var clearBtn *widgets.Paragraph
var generateBtn *widgets.Paragraph
var exitBtn *widgets.Paragraph

func main() {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	setUI()
	setDifficulty()
	handleEvents()
}

func setUI() {
	lvlTab = widgets.NewTabPane(" Easy", "Medium", "Hard")
	lvlTab.Title = "Difficulty"
	lvlTab.SetRect(50, 3, 74, 6)
	info = widgets.NewParagraph()
	info.SetRect(50, 6, 74, 9)
	info.Title = " info "
	solveBtn = widgets.NewParagraph()
	solveBtn.SetRect(50, 9, 74, 12)
	solveBtn.Text = " Solve"
	solveBtn.BorderStyle.Fg = ui.ColorBlue
	solveBtn.TextStyle.Fg = ui.ColorBlue
	clearBtn = widgets.NewParagraph()
	clearBtn.SetRect(50, 12, 74, 15)
	clearBtn.Text = " Clear "
	clearBtn.BorderStyle.Fg = 250
	generateBtn = widgets.NewParagraph()
	generateBtn.SetRect(50, 15, 74, 18)
	generateBtn.Text = " Randomize "
	generateBtn.TextStyle.Fg = ui.ColorYellow
	generateBtn.BorderStyle.Fg = ui.ColorYellow
	exitBtn = widgets.NewParagraph()
	exitBtn.SetRect(50, 18, 74, 21)
	exitBtn.BorderStyle.Fg = ui.ColorRed
	exitBtn.TextStyle.Fg = ui.ColorRed
	exitBtn.Text = " Exit "
	render()
}

func handleEvents() {
	editedRow, editedCol := -1, -1
	for e := range ui.PollEvents() {
		switch e.ID {
		case "<MouseLeft>":
			mouse := e.Payload.(ui.Mouse)
			if editedRow == -1 && editedCol == -1 {
				editedRow, editedCol = getEditedCell(mouse)
			}
			if inArea(mouse, solveBtn.Rectangle) {
				values := cellsToValues(notEditableCells(cells))
				Solve(0, 0, &values)
				updateBoardUI(values, true)
			}
			if inArea(mouse, clearBtn.Rectangle) {
				values := cellsToValues(notEditableCells(cells))
				updateBoardUI(values, false)
			}
			if inArea(mouse, generateBtn.Rectangle) {
				setDifficulty()
			}
			if inArea(mouse, exitBtn.Rectangle) {
				return
			}

		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if editedRow == -1 && editedCol == -1 {
				break
			}
			edited := &cells[editedRow][editedCol]
			num, err := strconv.Atoi(e.ID)
			if err != nil {
				logM(fmt.Sprintf("failed to convert string to integer: %v", err))
			}
			if num == 0 {
				resetColor(*edited)
				edited.Text = emptyCell
			} else {
				edited.BorderStyle.Fg = ui.ColorGreen
				edited.Text = fmt.Sprintf(" %v ", num)
			}
			edited.Value = num
			ui.Render(edited)
			editedRow, editedCol = -1, -1

		case "<Left>":
			lvlTab.FocusLeft()
			ui.Clear()
			setDifficulty()
			render()

		case "<Right>":
			lvlTab.FocusRight()
			ui.Clear()
			setDifficulty()
			render()

		case "q", "<C-z>", "<C-edited>":
			return

		}

		values := cellsToValues(cells)
		if Correct(values) {
			logM("You Won!!!!")
		} else {
			logM("")
		}
	}

}

func render() {
	ui.Render(info, solveBtn, clearBtn, generateBtn, lvlTab, exitBtn)
}

func getEditedCell(mouse ui.Mouse) (int, int) {
	for row := range cells {
		for col := range cells[row] {
			c := cells[row][col]
			if c.Editable && inArea(mouse, c.Rectangle) {
				c.BorderStyle.Fg = ui.ColorBlue
				ui.Render(c)
				return row, col
			}
		}
	}
	return -1, -1
}

func setDifficulty() {
	values = GenerateBoard(lvlTab.ActiveTabIndex + 1)
	cells = createBoardUI(*values)
}

func cellsToValues(cells [9][9]Cell) [9][9]int {
	var values [9][9]int
	for row := range cells {
		for col := range cells[row] {
			values[row][col] = cells[row][col].Value
		}
	}
	return values
}

func notEditableCells(cells [9][9]Cell) [9][9]Cell {
	var ret [9][9]Cell
	for row := range cells {
		for col := range cells[row] {
			if cells[row][col].Editable {
				cells[row][col].Value = 0
			}
			ret[row][col] = cells[row][col]
		}
	}
	return ret
}

func inArea(mouse ui.Mouse, c image.Rectangle) bool {
	return mouse.X > c.Min.X && mouse.X < c.Max.X &&
		mouse.Y > c.Min.Y && mouse.Y < c.Max.Y
}

func createBoardUI(values [9][9]int) [9][9]Cell {
	var cells [9][9]Cell
	for row := range values {
		for col := range values[row] {
			p := widgets.NewParagraph()
			c := &cells[row][col]
			c.Paragraph = p
			c.Value = values[row][col]
			if c.Value == 0 {
				c.Editable = true
				p.Text = "  "
			} else {
				p.Text = fmt.Sprintf(" %v ", c.Value)
			}
			colorByGrid(row, col, c)
			p.SetRect(col*cellHeight, row*cellWidth,
				col*cellHeight+cellHeight, row*cellWidth+cellWidth)
			ui.Render(c)
		}
	}
	return cells
}

func updateBoardUI(values [9][9]int, showEditable bool) {
	for row := range values {
		for col := range values[row] {
			c := &cells[row][col]
			c.Value = values[row][col]
			if c.Value == 0 {
				c.Text = "  "
				if !showEditable {
					colorByGrid(row, col, c)
				}
			} else if c.Editable {
				c.Text = fmt.Sprintf(" %v ", c.Value)
				c.BorderStyle.Fg = ui.ColorGreen
			}
			ui.Render(c)
		}
	}
}

func colorByGrid(row int, col int, c *Cell) {
	if isDiagonal(row, col) {
		c.BorderStyle.Fg = 240
	} else {
		c.BorderStyle.Fg = 255
	}
}

func resetColor(p Cell) {
	row := p.Min.X / cellHeight
	col := p.Min.Y / cellWidth
	if isDiagonal(row, col) {
		p.BorderStyle.Fg = 240
	} else {
		p.BorderStyle.Fg = 255
	}
	ui.Render(p)

}

func isDiagonal(row int, col int) bool {
	return (row < 3 && col < 3) || (row < 3 && col > 5) ||
		(row > 2 && row < 6 && col > 2 && col < 6) ||
		(row > 5 && col < 3) || (row > 5 && col > 5)
}

func logM(mess string) {
	info.Text = fmt.Sprintln(mess)
	ui.Render(info)

}
