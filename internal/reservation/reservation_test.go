package reservation

import (
	"testing"
)

func TestAddSlot(t * testing.T) {
	sys := System{}
	sys.Initialize()
	slot1 := Slot{Id: "a"}
	if !sys.Add(&slot1) {
		t.Error("Expected System.Add to add to unique id to a free slot")
	}

	if sys.slots[0] != &slot1 {
		t.Error("Expected System.Add to add slot at index 0")
	}
}

func TestRemoveSLot(t * testing.T) {
	sys := System{}
	sys.Initialize()
	slot1 := Slot{Id: "a"}
	sys.slots[0] = &slot1

	sys.RemoveFromSlot(&slot1)

	if sys.slots[0] != nil {
		t.Error("Expected slot to be nil'd after it was requested to be removed.")
	}
}

func TestCandidateNoReservation(t * testing.T) {
	sys := System{}
	sys.Initialize()

	slot1 := Slot{Id: "a"}
	sys.slots[0] = &slot1
	sys.queue[0] = &slot1

	if sys.ChooseReservationFromQueue() != nil {
		t.Error("Expected no candidates for slots.")
	}
}

func TestCandidateReservation(t * testing.T) {
	sys := System{}
	sys.Initialize()

	slot1 := Slot{Id: "a"}
	slot2 := Slot{Id: "b"}

	sys.slots[0] = &slot1
	sys.queue[0] = &slot2

	if sys.ChooseReservationFromQueue() != &slot2 {
		t.Error("Expected reservation candidate to be first item in queue.")
	}
}

func TestSorting(t * testing.T) {
	sys := System{}
	sys.Initialize()

	slot1 := Slot{Id: "a"}
	slot2 := Slot{Id: "b"}

	sys.slots[0] = &slot1
	sys.queue[0] = &slot1
	sys.queue[1] = &slot2
	sys.SortQueue()

	if sys.queue[0] == &slot1 || sys.queue[1] == &slot2 || sys.queue[2] != nil {
		t.Error("Did not sort properly.")
	}
}

