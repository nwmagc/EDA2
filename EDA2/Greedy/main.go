package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Esta es la estructura para representar una celda en la matriz
type Cell struct {
	x, y int
}

// Valida si una celda es existente
func isValidCell(rows, columns, new_x, new_y int) bool {
	return new_x >= 0 && new_y >= 0 && new_x < rows && new_y < columns
}

// Verifica si la celda es un parqueadero (si tiene un número o no)
func isParkingCell(symbolMatrix [][]string, x, y int) bool {
	return symbolMatrix[x][y] != ".."
}

// Verifica si la celda es un lugar por donde se puede pasar
func isBlankCell(symbolMatrix [][]string, x, y int) bool {
	return symbolMatrix[x][y] == ".."
}

// Verifica si la celda es una entrada o salida
func isAirportCell(symbolMatrix [][]string, x, y int) bool {
	return symbolMatrix[x][y] == "=="
}

// Verifica si la celda está bloqueada
func isBlockedCell(symbolMatrix [][]string, x, y int) bool {
	return symbolMatrix[x][y] == "##"
}

// Encuentra el estacionamiento disponible (celda con el menor peso) en la vecindad de la celda dada
func findBestParking(symbolMatrix [][]string, x, y, rows, columns int) Cell {
	movements := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	bestParking := Cell{-1, -1}
	minWeight := rows * columns * 2

	for _, move := range movements {
		new_x := x + move[0]
		new_y := y + move[1]

		if isValidCell(rows, columns, new_x, new_y) && !isBlockedCell(symbolMatrix, new_x, new_y) && !isParkingCell(symbolMatrix, new_x, new_y) && !isAirportCell(symbolMatrix, new_x, new_y) && !isBlankCell(symbolMatrix, new_x, new_y) {
			if weight := countParkingCells(symbolMatrix, new_x, new_y, rows, columns); weight < minWeight {
				minWeight = weight
				bestParking = Cell{new_x, new_y}
			}
		}
	}

	return bestParking
}


// Cuenta el número de parqueaderos alcanzables desde la celda dada
func countParkingCells(symbolMatrix [][]string, x, y, rows, columns int) int {
	count := 0
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, columns)
	}

	var dfs func(int, int)
	dfs = func(x, y int) {
		if !isValidCell(rows, columns, x, y) || visited[x][y] || isBlockedCell(symbolMatrix, x, y) || isBlankCell(symbolMatrix, x, y) || isAirportCell(symbolMatrix, x, y) {
			return
		}

		visited[x][y] = true
		count++

		dfs(x+1, y)
		dfs(x-1, y)
		dfs(x, y+1)
		dfs(x, y-1)
	}

	dfs(x, y)
	return count
}


// Resuelve el problema de asignación de estacionamientos utilizando un enfoque greedy
func solveProblem(symbolMatrix [][]string, events []int, rows, columns int) ([]string, bool) {
	solutionList := make([]string, 0)
	assignedEvents := make(map[int]bool)

	for _, event := range events {
		if event > 0 {
			bestParking := findBestParking(symbolMatrix, -1, -1, rows, columns)
			if bestParking.x == -1 && bestParking.y == -1 {
				return nil, false
			}

			symbolMatrix[bestParking.x][bestParking.y] = strconv.Itoa(event)
			assignedEvents[event] = true
			solutionList = append(solutionList, symbolMatrix[bestParking.x][bestParking.y]+" "+strconv.Itoa(event))
		} else {
			event = -event
			if _, ok := assignedEvents[event]; !ok {
				return nil, false
			}

			for i := 0; i < rows; i++ {
				for j := 0; j < columns; j++ {
					if symbolMatrix[i][j] == strconv.Itoa(event) {
						symbolMatrix[i][j] = ".."
						break
					}
				}
			}

			delete(assignedEvents, event)
		}
	}

	return solutionList, true
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	caseNumber := 0

	for scanner.Scan() {
		nfc := strings.Split(scanner.Text(), " ")
		number_of_airplanes, _ := strconv.Atoi(nfc[0])

		if number_of_airplanes == 0 {
			break
		} else {
			caseNumber++
			number_of_rows, _ := strconv.Atoi(nfc[1])
			number_of_columns, _ := strconv.Atoi(nfc[2])
			symbolMatrix := make([][]string, number_of_rows)

			for i := 0; i < number_of_rows; i++ {
				scanner.Scan()
				row := strings.Split(scanner.Text(), " ")
				symbolMatrix[i] = row
			}

			scanner.Scan()
			eventStrings := strings.Split(scanner.Text(), " ")
			events := make([]int, len(eventStrings))

			for i, eventString := range eventStrings {
				event, _ := strconv.Atoi(eventString)
				events[i] = event
			}

			result, success := solveProblem(symbolMatrix, events, number_of_rows, number_of_columns)

			if success {
				fmt.Printf("Case %d: Yes\n", caseNumber)
				for i := 0; i < len(result)-1; i++ {
					fmt.Printf("%s ", result[i])
				}
				fmt.Println(result[len(result)-1])
			} else {
				fmt.Printf("Case %d: No\n", caseNumber)
			}
		}
	}
}
