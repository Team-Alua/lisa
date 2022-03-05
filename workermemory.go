package main

type WorkerMemory struct {
	titleIdSlots [16]string
	activeTitleIds map[string]bool
}


// RegisterTitleId takes a titleId and assigns a slot to it.
// If the titleId is empty, it will return -2
// If the titleId is already registered, it will return -1.
// If there are no free slots, it will also return -1.
// Otherwise it return a slot index starting at zero.
func (wm * WorkerMemory) RegisterTitleId(titleId string) int {
	if titleId == "" {
		return -2
	}

	if _, ok := wm.activeTitleIds[titleId]; ok {
		return -1
	}

	for index, key := range wm.titleIdSlots {
		if key == "" {
			wm.titleIdSlots[index] = titleId
			wm.activeTitleIds[titleId] = true
			return index
		}
	}
	return -1
}

// FreeSlot takes a slot index and disassociates itself with the assigned title id.
// Freeing an empty slot is a no op.
func (wm * WorkerMemory) FreeSlot(index int) {
	titleId := wm.titleIdSlots[index]
	if titleId == "" {
		return
	}
	delete(wm.activeTitleIds, titleId)
	wm.titleIdSlots[index] = ""
	return
}
