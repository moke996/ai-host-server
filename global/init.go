package global

func Init() {
	// 连接MongoDB
	InitMongo()
	// 连接Redis
	InitRedis()

}

func Stop() {
	StopRedis()
}
