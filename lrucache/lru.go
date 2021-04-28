// +build !solution

package lrucache

type cell struct {
	key, value     int
	previous, next *cell
}

type cacheImpl struct {
	capacity int
	cache    map[int]*cell
	head     *cell
	tail     *cell
}

func New(cap int) Cache {
	return &cacheImpl{
		capacity: cap,
		cache:    map[int]*cell{},
	}
}

func (c *cacheImpl) Get(key int) (int, bool) {
	cacheCell, inCache := c.cache[key]

	rez := 0
	if inCache {
		rez = cacheCell.value
		c.updateUsage(cacheCell)
	}

	return rez, inCache
}

func (c *cacheImpl) Set(key, value int) {
	if c.capacity == 0 {
		return
	}

	cacheCell, inCache := c.cache[key]

	if !inCache {
		cacheCell = &cell{
			key:   key,
			value: value,
		}

		if len(c.cache) == c.capacity {
			c.removeLeastUsed()
		}
		c.cache[key] = cacheCell
	} else {
		cacheCell.value = value
	}

	c.updateUsage(cacheCell)
}

func (c *cacheImpl) Range(f func(key int, value int) bool) {
	cur := c.head
	for cur != nil && f(cur.key, cur.value) {
		cur = cur.next
	}
}

func (c *cacheImpl) Clear() {
	c.head, c.tail = nil, nil
	c.cache = make(map[int]*cell)
}

func (c *cacheImpl) removeLeastUsed() {
	forRemove := c.head
	if forRemove == nil {
		return
	}

	if forRemove.next != nil {
		forRemove.next.previous = nil
	} else {
		c.tail = nil
	}

	c.head = forRemove.next

	delete(c.cache, forRemove.key)
}

func (c *cacheImpl) updateUsage(forUpdate *cell) {
	if c.head == forUpdate {
		c.head = forUpdate.next
	}
	if c.tail == forUpdate {
		c.tail = forUpdate.previous
	}

	if forUpdate.previous != nil {
		forUpdate.previous.next = forUpdate.next
	}
	if forUpdate.next != nil {
		forUpdate.next.previous = forUpdate.previous
	}

	forUpdate.previous = c.tail
	if c.tail != nil {
		c.tail.next = forUpdate
	}
	c.tail = forUpdate

	if c.head == nil {
		c.head = forUpdate
	}
}
