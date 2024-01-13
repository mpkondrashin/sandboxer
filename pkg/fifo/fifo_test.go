package fifo

import "testing"

type Data struct {
	X int
	Y string
	Z bool
}

func TestMainFifo(t *testing.T) {
	dataSource := Data{
		X: 1,
		Y: "2",
		Z: true,
	}
	strSource := "test string"
	go func(t *testing.T) {
		t.Log("NewWriter")
		w, err := NewWriter()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("write data")
		if err := w.Write(&dataSource); err != nil {
			t.Error(err)
			return
		}
		t.Log("write str")
		if err := w.Write(strSource); err != nil {
			t.Error(err)
			return
		}
	}(t)
	t.Log("NewReader")
	r, err := NewReader()
	if err != nil {
		t.Fatal(err)
	}
	var dataTarget Data
	t.Log("read data")
	if err := r.Read(&dataTarget); err != nil {
		t.Fatal(err)
	}
	if dataSource != dataTarget {
		t.Fatalf("Expecter %v, but got %v", dataSource, dataTarget)
	}
	var strTarget string
	t.Log("read str")
	if err := r.Read(&strTarget); err != nil {
		t.Fatal(err)
	}
	if strSource != strTarget {
		t.Fatalf("Expecter %v, but got %v", strSource, strTarget)
	}
}
