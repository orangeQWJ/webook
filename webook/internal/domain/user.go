package domain

// User 领域对象, 是DDD中的(聚合根) entity
// BO (business object)
type User struct {
	Email    string
	Password string
}