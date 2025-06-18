package utils

import (
	"math"
	"moodmap/models"
)

// Parameters for DBSCAN
const (
	Eps    = 500.0 / 111.0
	MinPts = 2
)

// haversine calculates the great-circle distance between two geo-points.
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371.0 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// regionQuery finds all points within Eps distance from the given point.
func regionQuery(moods []models.Mood, idx int) []int {
	var neighbors []int
	for i := 0; i < len(moods); i++ {
		if i != idx && haversine(moods[idx].Lat, moods[idx].Lng, moods[i].Lat, moods[i].Lng) <= Eps &&
			moods[idx].Emoji == moods[i].Emoji { // Only group same-emotion
			neighbors = append(neighbors, i)
		}
	}
	return neighbors
}

// expandCluster grows the cluster from a given core point.
func expandCluster(moods []models.Mood, labels []int, idx int, clusterID int, neighbors []int) {
	labels[idx] = clusterID

	queue := append([]int{}, neighbors...)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if labels[current] == -1 {
			labels[current] = clusterID
		}
		if labels[current] != 0 {
			continue
		}

		labels[current] = clusterID
		otherNeighbors := regionQuery(moods, current)
		if len(otherNeighbors) >= MinPts {
			queue = append(queue, otherNeighbors...)
		}
	}
}

// DBSCAN clusters the mood posts by location and same emoji.
func DBSCAN(moods []models.Mood) []models.Mood {
	n := len(moods)
	labels := make([]int, n) // 0: unclassified, -1: noise, >0: cluster ID
	clusterID := 1

	for i := 0; i < n; i++ {
		if labels[i] != 0 {
			continue
		}

		neighbors := regionQuery(moods, i)
		if len(neighbors) < MinPts {
			labels[i] = -1 // mark as noise
		} else {
			expandCluster(moods, labels, i, clusterID, neighbors)
			clusterID++
		}
	}

	// Attach cluster IDs to moods
	for i := 0; i < n; i++ {
		moods[i].ClusterID = labels[i]
	}

	return moods
}
