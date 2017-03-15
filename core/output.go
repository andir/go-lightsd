package core

type Output interface {
    Source() string
    Render(stripe LEDStripeReader)
}
