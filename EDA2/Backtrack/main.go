package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

//Esta es la estructura para representar una celda en la matriz 
type Cell struct {
	x, y int
}

//Valida si una celda es existente
func isValidCell(rows, columns, new_x, new_y int) bool {
	return new_x >= 0 && new_y >= 0 && new_x < rows && new_y < columns //condición que verifica si una celda dada por las coordenadas (new_x, new_y) está dentro de los límites de la matriz
}

//Verifica si la celda es un parqueadero (si tiene un número o no)
func isParkingCell(symbolMatrix [][]string, x, y int) bool {
	_, err := strconv.Atoi(symbolMatrix[x][y])
	return err == nil
}

//Verifica si la celda es un lugar por donde se puede pasar
func isBlankCell(symbolMatrix [][]string, x, y int) bool {
	return symbolMatrix[x][y] == ".."
}

//Verifica si la celda es una entrada o salida
func isAirportCell(symbolMatrix [][]string, x, y int) bool {
	return symbolMatrix[x][y] == "=="
}

//Verifica si la celda está bloqueada
func isBlockedCell(symbolMatrix [][]string, x, y int) bool {
	return symbolMatrix[x][y] == "##"
}

// Asigna los pesos a las celdas dentro de la matriz
func assignWeights(symbolMatrix [][]string, weightMatrix [][]int, queue *[]Cell, movements [][]int, rows, columns int) ([][]int, map[Cell]bool) {
	visited := make(map[Cell]bool)
	parkingLots := make(map[Cell]bool)

	for len(*queue) > 0 {
		currentCell := (*queue)[0]
		*queue = (*queue)[1:]

		for _, move := range movements {
			new_x := currentCell.x + move[0]
			new_y := currentCell.y + move[1]

			if isValidCell(rows, columns, new_x, new_y) && !visited[Cell{new_x, new_y}] {
				if isBlockedCell(symbolMatrix, new_x, new_y) {
					weightMatrix[new_x][new_y] = -1
					visited[Cell{new_x, new_y}] = true
				}
				if isBlankCell(symbolMatrix, new_x, new_y) {
					weightMatrix[new_x][new_y] = weightMatrix[currentCell.x][currentCell.y]
					*queue = append([]Cell{{new_x, new_y}}, *queue...)
					visited[Cell{new_x, new_y}] = true
				}
				if isParkingCell(symbolMatrix, new_x, new_y) {
					weightMatrix[new_x][new_y] = weightMatrix[currentCell.x][currentCell.y] + 1
					*queue = append(*queue, Cell{new_x, new_y})
					visited[Cell{new_x, new_y}] = true
					parkingLots[Cell{new_x, new_y}] = true
				}
			}
		}
	}

	return weightMatrix, parkingLots
}

// Crea una matriz de pesos inicializada con ceros
func createWeightMatrix(rows, columns int) [][]int {
	weightMatrix := make([][]int, rows)
	for i := range weightMatrix {
		weightMatrix[i] = make([]int, columns)
	}
	return weightMatrix
}

func preprocess(symbolMatrix [][]string, rows, columns int) ([][]int, map[Cell]bool) {
	queue := []Cell{}
	movements := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	weightMatrix := createWeightMatrix(rows, columns)

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if isAirportCell(symbolMatrix, i, j) {
				weightMatrix[i][j] = 0
				queue = append(queue, Cell{i, j})
			}
		}
	}

	return assignWeights(symbolMatrix, weightMatrix, &queue, movements, rows, columns)
}

// Resuelve el problema de asignación de estacionamientos
func solveProblem(symbolMatrix [][]string, weightMatrix [][]int, events *[]int, rows, columns int, solutionList *[]string, dictionary map[int]Cell, parkingLots map[Cell]bool) bool {
	eventCount := 0

	if len(*events) == 0 {
		return true
	}

	for _, event := range *events {
		if event > 0 {
			eventCount++
		} else {
			break
		}
	}

	if eventCount > len(parkingLots) {
		return false
	} else {
		currentEvent := (*events)[0]
		*events = (*events)[1:]

		if currentEvent > 0 {
			for parkingLot := range parkingLots {
				if weightMatrix[parkingLot.x][parkingLot.y] != -1 {
					if weightMatrix[parkingLot.x][parkingLot.y] > 0 {
						*solutionList = append(*solutionList, symbolMatrix[parkingLot.x][parkingLot.y]+" "+strconv.Itoa(currentEvent))
						dictionary[currentEvent] = Cell{parkingLot.x, parkingLot.y}
						symbolMatrix[parkingLot.x][parkingLot.y] = "##"
						weightMatrix, parkingLots = preprocess(symbolMatrix, rows, columns)
						if solveProblem(symbolMatrix, weightMatrix, events, rows, columns, solutionList, dictionary, parkingLots) {
							return true
						} else {
							*solutionList = (*solutionList)[:len(*solutionList)-1]
							symbolMatrix[dictionary[currentEvent].x][dictionary[currentEvent].y] = strconv.Itoa(currentEvent)
							weightMatrix, parkingLots = preprocess(symbolMatrix, rows, columns)
							delete(dictionary, currentEvent)
						}
					}
				}
			}

			*events = append([]int{currentEvent}, *events...)
			return false
		} else {
			currentEvent2 := -currentEvent
			movements := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

			for _, move := range movements {
				new_x := dictionary[currentEvent2].x + move[0]
				new_y := dictionary[currentEvent2].y + move[1]

				if isValidCell(rows, columns, new_x, new_y) {
					if weightMatrix[new_x][new_y] != -11 && weightMatrix[new_x][new_y] != -1 {
						if weightMatrix[new_x][new_y] >= 0 {
							symbolMatrix[dictionary[currentEvent2].x][dictionary[currentEvent2].y] = strconv.Itoa(currentEvent)
							weightMatrix, parkingLots = preprocess(symbolMatrix, rows, columns)
							if solveProblem(symbolMatrix, weightMatrix, events, rows, columns, solutionList, dictionary, parkingLots) {
								return true
							}
							break
						}
					}
				}
			}

			symbolMatrix[dictionary[currentEvent2].x][dictionary[currentEvent2].y] = "##"
			weightMatrix, parkingLots = preprocess(symbolMatrix, rows, columns)
			*events = append([]int{currentEvent}, *events...)
			return false
		}
	}
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
			solutionList := []string{}
			dictionary := make(map[int]Cell)
			symbolMatrix := [][]string{}

			for i := 0; i < number_of_rows; i++ {
				scanner.Scan()
				row := strings.Split(scanner.Text(), " ")
				symbolMatrix = append(symbolMatrix, row)
			}

			scanner.Scan()
			eventStrings := strings.Split(scanner.Text(), " ")
			events := []int{}
			for _, eventString := range eventStrings {
				event, _ := strconv.Atoi(eventString)
				events = append(events, event)
			}

			weightMatrix, parkingLots := preprocess(symbolMatrix, number_of_rows, number_of_columns)
			result := solveProblem(symbolMatrix, weightMatrix, &events, number_of_rows, number_of_columns, &solutionList, dictionary, parkingLots)
			sortedSolution := make([]string, len(solutionList))
			copy(sortedSolution, solutionList)
			sort.Strings(sortedSolution)

			if result {
				fmt.Printf("Case %d: Yes\n", caseNumber)
				for i := 0; i < len(solutionList)-1; i++ {
					fmt.Printf("%s ", sortedSolution[i])
				}
				fmt.Println(sortedSolution[len(solutionList)-1])
			} else {
				fmt.Printf("Case %d: No\n", caseNumber)
			}
		}
	}
}


