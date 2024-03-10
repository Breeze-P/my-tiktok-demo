namespace go user

include "../base/common.thrift"
include "../base/user.thrift"

struct RegisterRequest {
    1: string username, // registered user name
    2: string password
}

struct RegisterResponse {
    1: common.BaseResponse base_resp,
    2: i64 user_id, // user id
}

struct LoginRequest {
    1: string username, // Login Username
    2: string password // login password
}

struct LoginResponse {
    1: common.BaseResponse base_resp,
    2: i64 user_id
}

struct GetUserInfoRequest {
    1: i64 user_id,
    2: i64 current_user_id
}

struct GetUserInfoResponse {
    1: common.BaseResponse base_resp,
    2: user.User user // User Info
}

struct RelationActionRequest {
    1: i64 current_user_id,
    2: i64 to_user_id,
    3: i32 action_type // 1-Follow, 2-Unfollow
}

struct RelationActionResponse {
    1: common.BaseResponse base_resp,
}

struct RelationFollowListRequest {
    1: i64 user_id,
    2: i64 current_user_id
}

struct RelationFollowListResponse {
    1: common.BaseResponse base_resp,
    2: list<user.User> user_list // User information list
}

struct RelationFollowerListRequest {
    1: i64 user_id,
    2: i64 current_user_id
}

struct RelationFollowerListResponse {
    1: common.BaseResponse base_resp,
    2: list<user.User> user_list
}

struct RelationFriendListRequest {
    1: i64 user_id,
    2: i64 current_user_id
}

struct RelationFriendListResponse {
    1: common.BaseResponse base_resp,
    2: list<user.FriendUser> user_list
}

struct MessageChatRequest {
    1: i64 current_user_id,
    2: i64 to_user_id, // recipient's user id
    3: i64 pre_msg_time
}

struct MessageChatResponse {
    1: common.BaseResponse base_resp,
    2: list<user.Message> message_list // message list
}

struct MessageActionRequest {
    1: i64 current_user_id, // user id of the req
    2: i64 to_user_id, // user id of the recipient
    3: i32 action_type, // 1-Send a message
    4: string content // Message content
}

struct MessageActionResponse {
    1: common.BaseResponse base_resp
}

struct CheckUserExitsByIdRequset {
    1: i64 user_id,
}

struct CheckUserExitsByIdResponse {
    1: common.BaseResponse base_resp,
    2: bool exits
}

service UserService {
    // auth
    LoginResponse Login(1: LoginRequest req)
    RegisterResponse Register(1: RegisterRequest req)
    GetUserInfoResponse GetUserInfo(1: GetUserInfoRequest req)
    CheckUserExitsByIdResponse CheckUserExitsById(1: CheckUserExitsByIdRequset req)

    // relation
    RelationActionResponse RelationAction(1: RelationActionRequest req)
    RelationFollowListResponse GetFollowList(1: RelationFollowListRequest req),
    RelationFollowerListResponse GetFollowerList(1: RelationFollowerListRequest req),
    RelationFriendListResponse GetFriendList(1: RelationFriendListRequest req)

    // message
    MessageChatResponse Chat(1: MessageChatRequest req),
    MessageActionResponse MessageAction(1: MessageActionRequest req)
}