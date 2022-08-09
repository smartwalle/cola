package cola

type Waiter interface {
	Add(delta int)

	Done()

	Wait()
}
