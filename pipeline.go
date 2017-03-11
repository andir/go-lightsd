package main

type Operation interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
	Render(*LEDStripe)
}



