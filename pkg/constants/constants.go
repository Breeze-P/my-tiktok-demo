package constants

const (
	UserTableName      = "users"
	FollowsTableName   = "follows"
	VideosTableName    = "videos"
	MessageTableName   = "messages"
	FavoritesTableName = "likes"
	CommentTableName   = "comments"

	VideoFeedCount       = 30
	FavoriteActionType   = 1
	UnFavoriteActionType = 2

	MinioVideoBucketName = "videobucket"
	MinioImgBucketName   = "imagebucket"

	TestSign       = "测试账号！ offer"
	TestAva        = "avatar/test1.jpg"
	TestBackground = "background/test1.png"
)

// connection information
const (
	MySQLDefaultDSN = "root:111111@tcp(127.0.0.1:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
	RedisAddr       = "localhost:6379"
	RedisPassword   = ""

	MinioEndPoint        = "localhost:9000"
	MinioAccessKeyID     = "mytiktok"
	MinioSecretAccessKey = "mytiktok111"
	MiniouseSSL          = false
)
