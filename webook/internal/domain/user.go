package domain

// User 领域对象, 是DDD中的(聚合根) entity
// BO (business object)
type User struct {
	Id         int64
	Email      string
	Password   string
	Nickname   string
	Birthday   string
	AboutMe    string
	Phone      string
	WechatInfo WechatInfo //将来可能还有钉钉Info
}
