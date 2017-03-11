package main

type Pipeline map[string] Operation

type Operation interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
	Render(*LEDStripe)
}



