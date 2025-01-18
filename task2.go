package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type Item struct {
	point    Point
	distance int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func findShortestPath(matrix [][]int, start, end Point) ([]Point, error) {
	directions := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	n, m := len(matrix), len(matrix[0])
	dist := make([][]int, n)
	prev := make([][]Point, n)

	for i := range dist {
		dist[i] = make([]int, m)
		prev[i] = make([]Point, m)
		for j := range dist[i] {
			dist[i][j] = math.MaxInt
		}
	}

	dist[start.x][start.y] = 0

	pq := &PriorityQueue{}
	heap.Push(pq, &Item{
		point:    start,
		distance: 0,
	})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item)
		x, y := current.point.x, current.point.y

		if current.point == end {
			var path []Point
			for current.point != start {
				path = append([]Point{current.point}, path...)
				current.point = prev[current.point.x][current.point.y]
			}
			path = append([]Point{start}, path...)
			return path, nil
		}

		for _, dir := range directions {
			next := Point{x + dir.x, y + dir.y}
			if next.x >= 0 && next.x < n && next.y >= 0 && next.y < m && matrix[next.x][next.y] != 0 {
				alt := dist[x][y] + matrix[next.x][next.y]
				if alt < dist[next.x][next.y] {
					dist[next.x][next.y] = alt
					prev[next.x][next.y] = current.point
					heap.Push(pq, &Item{
						point:    next,
						distance: alt,
					})
				}
			}
		}
	}

	return nil, errors.New("нет решения")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var n, m int
	var matrix [][]int
	var start, end Point

	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "ошибка ввода")
		return
	}
	_, err := fmt.Sscanf(scanner.Text(), "%d %d", &n, &m)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка парсинга")
		return
	}

	matrix = make([][]int, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "ошибка ввода")
			return
		}
		line := scanner.Text()
		matrixLine := strings.Fields(line)
		if len(matrixLine) != m {
			fmt.Fprintln(os.Stderr, "неверное количество элементов в строке матрицы")
			return
		}
		matrix[i] = make([]int, m)

		for j := 0; j < m; j++ {
			matrix[i][j], err = strconv.Atoi(matrixLine[j])
			if err != nil {
				fmt.Fprintln(os.Stderr, "ошибка парсинга")
				return
			}
		}
	}

	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "ошибка ввода")
		return
	}
	_, err = fmt.Sscanf(scanner.Text(), "%d %d %d %d", &start.x, &start.y, &end.x, &end.y)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка парсинга")
		return
	}

	path, err := findShortestPath(matrix, start, end)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка работы функции поиска пути -", err.Error())
		return
	}

	for _, point := range path {
		fmt.Printf("%d %d\n", point.x, point.y)
	}
	fmt.Println(".")
}
