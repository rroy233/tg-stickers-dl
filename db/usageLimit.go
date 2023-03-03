package db

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/StickerDownloader/config"
	"strconv"
	"time"
)

// CheckLimit 用户是否已达到今日限额
func CheckLimit(update *tgbotapi.Update) bool {
	UID := int64(0)
	if update.Message != nil {
		UID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		UID = update.CallbackQuery.Message.Chat.ID
	} else {
		return true
	}
	if UID == config.Get().General.AdminUID {
		return false
	}
	limit := rdb.Get(ctx, fmt.Sprintf("%s:UserLimit:%d", ServicePrefix, UID)).Val()
	if limit == "" {
		rdb.Set(ctx, fmt.Sprintf("%s:UserLimit:%d", ServicePrefix, UID), 1, 24*time.Hour)
		return false
	}

	limitTimes, _ := strconv.Atoi(limit)
	if limitTimes > config.Get().General.UserDailyLimit {
		return true
	}
	rdb.Set(ctx, fmt.Sprintf("%s:UserLimit:%d", ServicePrefix, UID), limitTimes+1, redis.KeepTTL)
	return false
}

// 获取该用户已使用的次数
func getUsed(UID int64) int {
	limit := rdb.Get(ctx, fmt.Sprintf("%s:UserLimit:%d", ServicePrefix, UID)).Val()
	if limit == "" {
		return -1
	}
	limitTimes, _ := strconv.Atoi(limit)
	return limitTimes
}

// GetLimit 获取该用户今日剩余可用次数
func GetLimit(UID int64) int {
	if UID == config.Get().General.AdminUID {
		return -1
	}
	limitTimes := getUsed(UID)
	if limitTimes == -1 {
		return config.Get().General.UserDailyLimit
	}
	return config.Get().General.UserDailyLimit - limitTimes
}
