package service

type (
	List struct {
		Nacos *Nacos
	}
)

var (
	Factory *List
)

func Register() *List {
	Factory = &List{
		Nacos: new(Nacos),
	}
	return Factory
}
