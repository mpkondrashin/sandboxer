package shm

import "testing"

type Data struct {
	X int
	Y string
	Z bool
}

func TestSHM(t *testing.T) {
	w, err := NewSHMWriter(1024)
	if err != nil {
		t.Fatal(err)
	}
	r, err := NewSHMReader(1024)
	if err != nil {
		t.Fatal(err)
	}
	dataSource := Data{
		X: 1,
		Y: "2",
		Z: true,
	}
	if err := w.Write(dataSource); err != nil {
		t.Fatal(err)
	}
	var dataTarget Data
	if err := r.Read(&dataTarget); err != nil {
		t.Fatal(err)
	}
	if dataSource != dataTarget {
		t.Fatalf("Expecter %v, but got %v", dataSource, dataTarget)
	}
}
