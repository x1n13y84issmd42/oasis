package log

type IStyle interface {
	styleDefault(...interface{}) string
	styleURL(...interface{}) string
	styleMethod(...interface{}) string
	styleOp(...interface{}) string
	styleOK(...interface{}) string
	styleFailure(...interface{}) string
	styleSuccess(...interface{}) string
	styleError(...interface{}) string
	styleID(...interface{}) string
	styleValueExpected(...interface{}) string
	styleValueActual(...interface{}) string
}
