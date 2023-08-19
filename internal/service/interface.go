package service

type Starter interface {
	Start() error
}

type Stopper interface {
	Stop() error
}

type Checker interface {
	Check() error
}

type Service interface {
	Starter
	Stopper
	Checker
	GetPriority() int32
	GetID() int32
	GetPersist() int32
}
