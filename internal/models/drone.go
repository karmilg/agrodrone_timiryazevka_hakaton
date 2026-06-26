package models


type DroneStatus string
const (
	StatusAfk DroneStatus = "Ждет"
	StatusFly DroneStatus = "Летит"
	StatusScan DroneStatus = "Сканирует"
	StatusCharge DroneStatus = "Заряжается"
)

type Drone struct {
	ID int
	Battery int
	Name string
	StatusDrone DroneStatus
	X int
	Y int
	BaseX int
	BaseY int
}

func manhatan(x1, x2, y1, y2 int) int {
	return abs(x1 - x2) + abs(y1 - y2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func (d *Drone) Fly(tx, ty int) bool {
	s := manhatan(d.X, tx, d.Y, ty)
	if d.Battery < s {
		return false
	}

	d.X, d.Y = tx, ty
	d.Battery -= s
	d.StatusDrone = StatusFly
	return true
}


func (d *Drone) Charge(s int) {
	d.StatusDrone = StatusCharge
	d.Battery += s
	if d.Battery > 100 {
		d.Battery = 100
	} else if d.Battery == 100 {
		d.StatusDrone = StatusAfk
	}
}


type ScoutDrone struct {
	Drone
	ScanRadius int
}

func (s *ScoutDrone) Scan(c *Cell, ndvi, soil, temp float64) {
	c.NDVI = ndvi
	c.Soil = soil
	c.Temp = temp
	c.Scanned = true
	s.StatusDrone = StatusScan
	s.Battery -= 10
}

type SprayDrone struct {
	Drone
	Cap float64
}

func (sd *SprayDrone) Spray(c *Cell, ls float64) bool {
	if sd.Cap < ls {
		return false
	}
	sd.Cap -= ls 
	sd.Battery -= 10
	c.Soil += 5
	if c.Soil > 100 {
		c.Soil = 100
	}

	return true
}


