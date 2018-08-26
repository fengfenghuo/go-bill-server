package redis

type DataRedis interface {
	Set(key interface{}, data interface{}) error
	Get(key interface{}) (interface{}, error)
	HashSet(key interface{}, field interface{}, data interface{}) error
	HashGet(key interface{}, field interface{}) (interface{}, error)
	HashAppend(key interface{}, field interface{}, data interface{}) error
	HashIsExist(key interface{}, field interface{}) bool
	SetAppend(key interface{}, data interface{}) error
	SetDelete(key interface{}, data interface{}) error
	SetGet(key interface{}) (interface{}, error)
}
