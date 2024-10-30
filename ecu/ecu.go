package ecu

// Struktur data untuk ECU
type ECUData struct {
	RPM      int
	Gear     int
	Speed    int
	IsAttack bool // true jika status=1
}
