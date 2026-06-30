package domain

type Task struct {
	MissionName string     `json:"mission_name"`
	Takeoff     float64    `json:"takeoff"`
	RouteType   string     `json:"route_type"`
	Waypoints   []Waypoint `json:"waypoints"`
}

type Waypoint struct {
	Distance float64 `json:"distance"`
	Azimuth  float64 `json:"azimuth"`
	Altitude float64 `json:"altitude"`
	Hold     float64 `json:"hold"`
}
