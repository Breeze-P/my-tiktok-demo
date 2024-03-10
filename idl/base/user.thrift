namespace go base

struct User {
    1: i64 id,
    2: string name,
    3: i64 follow_count,
    4: i64 follower_count,
    5: bool is_follow,
    6: string avatar,
    7: string background_image,
    8: string signature,
    9: i64 total_favorited,
    10: i64 work_count,
    11: i64 favorite_count
}

struct FriendUser {
    1: User user,
    2: string message, // latest chat messages with this friend
    3: i64 msgType // message type, 0 - Messages currently requested by the user, 1 - Messages sent by the current requesting user
}

struct Message {
    1: i64 id, // message id
    2: i64 to_user_id, // The id of the recipient of the message
    3: i64 from_user_id, // The id of the sender of the message
    4: string content, // Message content
    5: i64 create_time // message creation time
}