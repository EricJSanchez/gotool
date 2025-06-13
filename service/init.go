package service

type (
	List struct {
		Nacos *Nacos
	}
)

var (
	Factory *List
)

func Register() {
	Factory = &List{
		Nacos: new(Nacos),
	}
	return
}
