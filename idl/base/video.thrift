namespace go base
include "./user.thrift"

struct Video {
    1: i64 id,
    2: user.User author,
    3: string play_url,
    4: string cover_url,
    5: i64 favorite_count,
    6: i64 comment_count,
    7: bool is_favorite,
    8: string title
}

struct Comment {
    1: i64 id, // video comment id
    2: user.User user, // comment user information
    3: string content, // comment
    4: string create_date // comment publication date, format mm-dd
}