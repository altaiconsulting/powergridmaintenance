package main

import (
	"fmt"
)

// powerGrid type represents a single power grid with connected power stations
type powerGrid struct {
	// minHeap is a slice representing min heap data structure into which grid power stations are organized by their ids
	minHeap []*powerStation
}

// addPowerStation adds power station with the specified id to power grid
func (grid *powerGrid) addPowerStation(stationID int, stationPool *powerStationPool) {
	station := stationPool.get(stationID)
	station.grid = grid
	station.id = stationID
	station.online = true
	grid.minHeap = append(grid.minHeap, station)
}

// siftDown moves the power station located at the index "index" in min-heap slice down the min heap tree
// by successively exchanging this power station with the power station
// with smaller id among the stations located in child nodes of this power station node
func (grid *powerGrid) siftDown(index int) {
	length := len(grid.minHeap)
	for index < length && index >= 0 {
		leftIndex, rightIndex := index<<1+1, index<<1+2
		if leftIndex >= length {
			break
		}
		smallerChildIndex := leftIndex
		if rightIndex < length && grid.minHeap[rightIndex].id < grid.minHeap[leftIndex].id {
			smallerChildIndex = rightIndex
		}
		if grid.minHeap[smallerChildIndex].id >= grid.minHeap[index].id {
			break
		}
		grid.minHeap[smallerChildIndex], grid.minHeap[index] = grid.minHeap[index], grid.minHeap[smallerChildIndex]
		index = smallerChildIndex
	}
}

// buildHeap urilizes siftDown to implement heapify method restoring min-heap property
func (grid *powerGrid) buildHeap() {
	for i := len(grid.minHeap)>>1 - 1; i >= 0; i-- {
		grid.siftDown(i)
	}
}

// removeMin removes power station with the smallest id from the grid
func (grid *powerGrid) removeMin() {
	lastIndex := len(grid.minHeap) - 1
	grid.minHeap[0] = grid.minHeap[lastIndex]
	grid.minHeap = grid.minHeap[:lastIndex]
	grid.siftDown(0)
}

// getMin retrieves power station with the smallest id in the grid
func (grid *powerGrid) getMin() *powerStation {
	return grid.minHeap[0]
}

// getOperationalStationMinID returns operational power station with the smallest id in the grid
func (grid *powerGrid) getOperationalStationMinID() int {
	for len(grid.minHeap) > 0 {
		stationWithMinID := grid.getMin()
		if stationWithMinID.online {
			return stationWithMinID.id
		}
		grid.removeMin()
	}
	return -1
}

// powerStation type represents a single power station that is part of a power grid
type powerStation struct {
	// grid is a pointer to the power grid this power station belongs to
	grid *powerGrid
	// id is the power station id unique across the entire power grid interconnection
	id int
	// online is true when power station is online and false if it is offline
	online bool
}

// resolveMaintenanceCheck resolves maintenance check request to the power station
func (station *powerStation) resolveMaintenanceCheck() int {
	if station.online {
		return station.id
	}
	return station.grid.getOperationalStationMinID()
}

// moveOffline moves the power station offline
func (station *powerStation) moveOffline() {
	station.online = false
}

// powerStationPool represents the contiguous pool of power stations indexed by their ids
type powerStationPool struct {
	stations *[]powerStation
}

// newStationPool constructs the pool of power stations of the specified size
func newStationPool(size int) *powerStationPool {
	powerStations := make([]powerStation, size+1)
	return &powerStationPool{stations: &powerStations}
}

// get fetches the power station with the specified id from the pool
func (pool *powerStationPool) get(stationID int) *powerStation {
	station := &(*pool.stations)[stationID]
	return station
}

// PowerGridInterconnection represents the power grid interconnection consisting of
// the pool of all interconnection power stations and the collection of disconnected power grids
type PowerGridInterconnection struct {
	stationPool *powerStationPool
	grids       []*powerGrid
}

// NewPowerGridInterconnection constructs new power grid interconnection
// from 2D array "connections" representing the connections between power stations
func NewPowerGridInterconnection(c int, connections [][]int) *PowerGridInterconnection {
	// build adjacency list of the entire interconnection
	graph := make([][]int, c+1)
	for _, edge := range connections {
		node1, node2 := edge[0], edge[1]
		graph[node1] = append(graph[node1], node2)
		graph[node2] = append(graph[node2], node1)
	}
	// create a pool of power stations in interconnection
	stationPool := newStationPool(c)
	// create a slice for power grids in interconnection
	grids := []*powerGrid{}
	// visited represents power station visitation slice in DFS process finding power stations connected into power grids
	visited := make([]bool, c+1)
	// DFS for finding connected components in the graph
	for nd := 1; nd <= c; nd++ {
		if visited[nd] {
			continue
		}
		// create new grid
		grid := &powerGrid{}
		// add grid to interconnection
		grids = append(grids, grid)
		stack := []int{nd}
		for len(stack) > 0 {
			node := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if visited[node] {
				continue
			}
			grid.addPowerStation(node, stationPool)
			visited[node] = true
			for _, neighbor := range graph[node] {
				if visited[neighbor] {
					continue
				}
				stack = append(stack, neighbor)
			}
		}
	}
	for _, grid := range grids {
		grid.buildHeap()
	}
	return &PowerGridInterconnection{
		stationPool: stationPool,
		grids:       grids,
	}
}

// getStationByID fetches interconnection power station by its id
func (interconnection *PowerGridInterconnection) getStationByID(stationID int) *powerStation {
	return interconnection.stationPool.get(stationID)
}

// ResolveMaintenanceCheckForStation resolves maintenance check request to the power station with the specified id
func (interconnection *PowerGridInterconnection) ResolveMaintenanceCheckForStation(stationID int) int {
	return interconnection.getStationByID(stationID).resolveMaintenanceCheck()
}

// MoveStationOffline moves offline the power station with the specified id
func (interconnection *PowerGridInterconnection) MoveStationOffline(stationID int) {
	interconnection.getStationByID(stationID).moveOffline()
}

// processQueries constructs power grid interconnections from 2D array "connections"
// and processes queries to this interconnection
func processQueries(c int, connections [][]int, queries [][]int) []int {
	interconnection := NewPowerGridInterconnection(c, connections)
	result := []int{}
	for _, query := range queries {
		op, stationID := query[0], query[1]
		switch op {
		case 1:
			result = append(result, interconnection.ResolveMaintenanceCheckForStation(stationID))
		case 2:
			interconnection.MoveStationOffline(stationID)
		default:
		}
	}
	return result
}

func main() {
	connections := [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}
	queries := [][]int{{1, 3}, {2, 1}, {1, 1}, {2, 2}, {1, 2}}
	fmt.Println(processQueries(5, connections, queries))
	connections = [][]int{}
	queries = [][]int{{1, 1}, {2, 1}, {1, 1}}
	fmt.Println(processQueries(3, connections, queries))
	connections = [][]int{}
	queries = [][]int{{1, 1}, {2, 1}, {2, 1}, {2, 1}, {2, 1}}
	fmt.Println(processQueries(1, connections, queries))
	connections = [][]int{{1, 2}}
	queries = [][]int{{1, 1}, {1, 2}, {1, 2}, {2, 2}, {2, 2}, {1, 1}, {1, 2}, {1, 1}}
	fmt.Println(processQueries(2, connections, queries))
	connections = [][]int{{2, 1}}
	queries = [][]int{{2, 1}, {1, 2}, {2, 1}, {1, 1}, {1, 2}, {1, 1}, {1, 1}, {2, 1}, {2, 2}}
	fmt.Println(processQueries(2, connections, queries))
	connections = [][]int{{2, 4}, {1, 2}, {5, 4}, {4, 6}, {2, 6}, {3, 6}, {4, 1}}
	queries = [][]int{{1, 1}, {1, 1}, {2, 1}, {2, 6}, {1, 6}, {1, 6}, {1, 3}, {1, 4}, {1, 4}, {2, 2}, {1, 2}, {2, 1}, {1, 4}, {2, 1}, {1, 6}, {1, 5}, {1, 2}, {2, 5}, {1, 2}, {2, 4}}
	fmt.Println(processQueries(6, connections, queries))
	connections = [][]int{{17, 9}, {9, 14}, {1, 3}, {10, 12}, {6, 2}, {3, 12}, {3, 15}, {8, 11}, {9, 4}, {13, 1}, {1, 8}, {12, 8}, {17, 7}, {17, 16}, {9, 12}, {13, 3}, {1, 16}, {15, 12}, {7, 14}}
	queries = [][]int{{2, 10}, {1, 12}, {1, 13}, {2, 1}, {2, 7}, {2, 3}, {2, 11}, {1, 9}, {2, 11}, {2, 4}, {2, 5}, {2, 7}, {1, 14}, {2, 1}, {1, 1}, {1, 7}, {2, 4}, {2, 3}, {2, 14}, {1, 1}, {2, 15}, {1, 15}, {1, 6}, {2, 15}, {2, 10}, {1, 1}, {1, 17}, {2, 9}, {1, 12}, {1, 17}, {1, 4}, {1, 5}, {2, 7}, {2, 8}, {1, 14}, {1, 16}, {1, 3}, {2, 17}}
	fmt.Println(processQueries(17, connections, queries))
}
