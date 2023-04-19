// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package load

import (
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"github.com/changaolee/skeleton/pkg/log"
)

// NotificationCommand 定义一个新的通知类型.
type NotificationCommand string

// 定义 Redis 订阅相关键和事件.
const (
	RedisPubSubChannel                      = "skt.notifications"
	NoticePolicyChanged NotificationCommand = "PolicyChanged"
	NoticeSecretChanged NotificationCommand = "SecretChanged"
)

// Notification 是一个用于发布和订阅消息编码类型.
type Notification struct {
	Command       NotificationCommand `json:"command"`
	Payload       string              `json:"payload"`
	Signature     string              `json:"signature"`
	SignatureAlgo crypto.Hash         `json:"algorithm"`
}

// Sign 使用 SHA256 算法进行签名.
func (n *Notification) Sign() {
	n.SignatureAlgo = crypto.SHA256
	hash := sha256.Sum256([]byte(string(n.Command) + n.Payload))
	n.Signature = hex.EncodeToString(hash[:])
}

func handleRedisEvent(v interface{}, handled func(NotificationCommand), reloaded func()) {
	message, ok := v.(*redis.Message)
	if !ok {
		return
	}

	notif := Notification{}
	if err := json.Unmarshal([]byte(message.Payload), &notif); err != nil {
		log.Errorf("Unmarshalling message body failed, malformed: ", err)

		return
	}
	log.Infow("Receive redis message", "command", notif.Command, "payload", message.Payload)

	switch notif.Command {
	case NoticePolicyChanged, NoticeSecretChanged:
		log.Infow("Reloading secrets and policies")
		reloadQueue <- reloaded
	default:
		log.Warnf("Unknown notification command: %q", notif.Command)

		return
	}

	if handled != nil {
		handled(notif.Command)
	}
}
