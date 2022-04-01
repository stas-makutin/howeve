package core

type Keyable struct {
	key interface{}
}

func (k *Keyable) Key() interface{} {
	return k.key
}

func (k *Keyable) WithKey(key interface{}) {
	k.key = key
}
