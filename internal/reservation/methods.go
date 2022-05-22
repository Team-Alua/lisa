package reservation

func (rs * System) Initialize() {
	if rs.queue == nil {
		rs.queue = make([]*Slot, 16)
	}
}

func (rs * System) isIdInSlot(id string) bool {
	for _, slotElement := range rs.slots {
		if slotElement == nil {
			continue
		}
		if slotElement.Id == id {
			return true
		}
	}
	return false
}

func (rs * System) insertIntoFreeSlot(r *Slot) {
	index := 0
	for i, slotElement := range rs.slots {
		if slotElement == nil {
			index = i
			break
		}
	}
	rs.slots[index] = r
}

func (rs * System) insertIntoQueue(r *Slot) {

	l := len(rs.queue)
	if rs.queue[l - 1] != nil {
		newQueue := make([]*Slot, l * 2)
		for i, slotElement := range rs.queue {
			newQueue[i] = slotElement
		}
		rs.queue = newQueue
	}

	index := -1
	for i, slotElement := range rs.queue {
		if slotElement == nil {
			index = i
			break
		}
	}
	rs.queue[index] = r
}

func (rs * System) SortQueue() {
	visited := map[string]bool{}

	for _, slotElement := range rs.slots {
		visited[slotElement.Id] = true
	}

	marker := -1

	for i, queueElement := range rs.queue {
		// End of Queue
		if queueElement == nil {
			break
		}

		// In a slot or furthur in queue
		// so need to push everything that isn't
		// closer
		if visited[queueElement.Id] {
			marker = i
			break
		} else {
			visited[queueElement.Id] = true
		}
	}

	// No sorting necessary
	sl := len(rs.slots)
	if marker == -1 || marker >= sl {
		return
	}

	ql := len(rs.queue)
	for i := marker + 1; marker < sl; i++ {
		// Reached physical end of queue
		if i == ql {
			break
		}

		// Reached logical end of queue
		if rs.queue[i] == nil {
			break
		}

		// Ignore already visited elements
		if visited[rs.queue[i].Id] {
			continue
		}

		temp := rs.queue[i]

		// Move all elements back
		for j := i; j > marker; j-- {
			rs.queue[j] = rs.queue[j - 1]
		}
		rs.queue[marker] = temp
		marker++
	} 
}

func (rs * System) ChooseReservationFromQueue() *Slot {
	visited := map[string]bool{}

	for _, slotElement := range rs.slots {
		if slotElement == nil {
			continue
		}
		visited[slotElement.Id] = true
	}
	
	for _, queueElement := range rs.queue {

		if queueElement == nil {
			break
		}

		if visited[queueElement.Id] {
			continue
		}

		return queueElement
	}
	return nil
}

func (rs * System) Add(r *Slot) bool {
	if !rs.isIdInSlot(r.Id) {
		rs.insertIntoFreeSlot(r)
		return true
	}

	rs.insertIntoQueue(r)
	return false
}

func (rs * System) RemoveFromSlot(r *Slot) bool {
	for i, slotElement := range rs.slots {
		if slotElement == r {
			rs.slots[i] = nil
			return true
		}
	}
	return false
}

func (rs * System) RemoveFromQueue(r *Slot) bool {
	index := -1
	for i, queueElement := range rs.queue {
		if queueElement == r {
			index = i
			break
		}
	}


	// Not found
	if index == -1 {
		return false
	}

	ql := len(rs.queue)
	for ; index + 1 < ql; index++ {
		// Logical end of queue
		if rs.queue[index] == nil {
			break
		}
		rs.queue[index] = rs.queue[index + 1]
	}
	rs.queue[ql - 1] = nil

	return true
}

