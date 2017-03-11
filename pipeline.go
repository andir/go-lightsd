package main

type Pipeline []Operation

type Operation interface {
	Name() string

	Lock()
	Unlock()
	RLock()
	RUnlock()

	Render(LEDStripe) LEDStripe
}



