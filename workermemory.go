package main

type WorkerMemory struct {
	slots [16]string
	activeSlotName map[string]bool
	clientQueue []*SaveClient
}


func (wm * WorkerMemory) Initialize() {
	wm.activeSlotName = make(map[string]bool)
	wm.clientQueue = make([]*SaveClient, 0)
}

// ReserveSlot takes a slotName and assigns a slot to it.
// If the slotName is already registered, it will return -1.
// If there are no free slots, it will also return -1.
// Otherwise it return a slot index starting at zero.
func (wm * WorkerMemory) ReserveSlot(slotName string) int {
	if slotName == "" {
		panic("Title Id was an empty string.")
	}

	if _, ok := wm.activeSlotName[slotName]; ok {
		return -1
	}

	for index, key := range wm.slots {
		if key == "" {
			wm.slots[index] = slotName
			wm.activeSlotName[slotName] = true
			return index
		}
	}
	return -1
}

// FreeSlot takes a slot index and disassociates itself with the assigned title id.
// Freeing an empty slot is a no op.
func (wm * WorkerMemory) FreeSlot(index int) {
	slotName := wm.slots[index]
	if slotName == "" {
		return
	}
	delete(wm.activeSlotName, slotName)
	wm.slots[index] = ""
	return
}

func (wm * WorkerMemory) AddToQueue(client *SaveClient) {
	wm.clientQueue = append(wm.clientQueue, client)
}

func (wm * WorkerMemory) GetClientFromQueue() *SaveClient {
	if len(wm.clientQueue) > 0 {
		client := wm.clientQueue[0]
		wm.clientQueue = wm.clientQueue[1:]
		return client
	}
	return nil
}
