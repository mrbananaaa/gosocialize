package logger

import "go.uber.org/zap"

func String(key, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

func Bytes(key string, value []byte) zap.Field {
	return zap.Binary(key, value)
}
