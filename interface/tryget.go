package _interface


type TryGet struct {
}

type ITryGet interface {
	TrxGet(address string, sourceAddress string, lable string) string
	TryGetCoin(address string ,coins ...string)
	TryGetCoinAll(address string ,coins ...string)
}