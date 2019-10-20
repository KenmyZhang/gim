package cache

import (
	"goim/logic/db"
	"goim/logic/model"
	"goim/public/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	GroupUserKey = "group_user:"
	GroupUserExp = 2 * time.Hour
)

type groupUserCache struct{}

var GroupUserCache = new(groupUserCache)

func (*groupUserCache) Key(appId, groupId int64) string {
	return GroupUserKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

// Set 保存群组所有用户的信息
func (c *groupUserCache) Set(appId, groupId int64, userInfos []model.GroupUser) error {
	err := set(c.Key(appId, groupId), userInfos, GroupUserExp)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// GetAll 获取群组的所有用户，如果缓存里面没有，返回nil
func (c *groupUserCache) Get(appId, groupId int64) ([]model.GroupUser, error) {
	var users []model.GroupUser
	err := get(c.Key(appId, groupId), &users)
	if err != nil && err != redis.Nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	if err == redis.Nil {
		return nil, nil
	}
	return users, nil
}

// Del 删除缓存
func (c *groupUserCache) Del(appId, groupId int64) error {
	_, err := db.RedisCli.Del(c.Key(appId, groupId)).Result()
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}