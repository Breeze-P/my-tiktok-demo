namespace go video

include "../base/common.thrift"
include "../base/video.thrift"

struct FeedRequest {
    1: i64 latest_time, // Optional, limit the time stamp of the latest submission of the returned video, accurate to the second
                       // and leave it blank to indicate the current time
    2: i64 current_user_id,
}

struct FeedResponse {
    1: common.BaseResponse base_resp,
    2: list<video.Video> video_list,
    3: i64 next_time // publish the earliest time as the latest time in the next request
}

struct PublishActionRequest {
    1: i64 current_user_id,
    2: string play_url,
    3: string cover_url,
    4: string title
    5: i64 publish_time
}

struct PublishActionResponse {
    1: common.BaseResponse base_resp,
}

struct PublishListRequest {
    1: i64 user_id,
    2: i64 current_user_id,
}

struct PublishListResponse {
    1: common.BaseResponse base_resp,
    2: list<video.Video> video_list
}


struct FavoriteActionRequest {
    1: i64 current_user_id,
    2: i64 video_id, // video id
    3: i32 action_type // 1-like, 2-unlike
}

struct FavoriteActionResponse {
    1: common.BaseResponse base_resp,
    2: bool success // success or not 
}

struct FavoriteListRequest {
    1: i64 user_id,
    2: i64 current_user_id,
}

struct FavoriteListResponse {
    1: common.BaseResponse base_resp,
    2: list<video.Video> video_list // list of videos liked by users
}

struct CommentActionRequest {
    1: i64 current_user_id, // user id
    2: i64 video_id,
    3: i32 action_type, // 1- Post a comment, 2- Delete a comment
    4: string comment_text, // Comment content filled in by users，when action type=1
    5: i64 comment_id // The id of the comment to delete，when action type=2
}

struct CommentActionResponse {
    1: common.BaseResponse base_resp,
    2: video.Comment comment // return the comment content, no need to re-pull the entire list
}

struct CommentListRequest {
    1: i64 current_user_id,
    2: i64 video_id
}

struct CommentListResponse {
    1: common.BaseResponse base_resp,
    2: list<video.Comment> comment_list // return comment list
}

struct GetWorkCountRequest {
    1: i64 current_user_id,
}

struct GetWorkCountResponse {
    1: common.BaseResponse base_resp,
    2: i64 count
}

struct GetFavoriteCountByUserIDRequest {
    1: i64 current_user_id,
}

struct GetFavoriteCountByUserIDResponse {
    1: common.BaseResponse base_resp,
    2: i64 count
}

struct QueryTotalFavoritedByAuthorIDRequest {
    1: i64 current_user_id,
}

struct QueryTotalFavoritedByAuthorIDResponse {
    1: common.BaseResponse base_resp,
    2: i64 count
}

service VideoService {
    FeedResponse GetFeed(1: FeedRequest req)

    PublishActionResponse PublishAction(1: PublishActionRequest request),
    PublishListResponse GetPublishList(1: PublishListRequest request)

    FavoriteActionResponse FavoriteAction(1: FavoriteActionRequest req)
    FavoriteListResponse GetFavoriteList(1: FavoriteListRequest req)
    GetFavoriteCountByUserIDResponse GetFavoriteCountByUserID(1: GetFavoriteCountByUserIDRequest req)
    QueryTotalFavoritedByAuthorIDResponse QueryTotalFavoritedByAuthorID(1: QueryTotalFavoritedByAuthorIDRequest req)

    GetWorkCountResponse GetWorkCount(1: GetWorkCountRequest req)

    // comment
    CommentActionResponse CommentAction(1: CommentActionRequest req),
    CommentListResponse GetCommentList(1: CommentListRequest req)
}