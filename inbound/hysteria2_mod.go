//go:build with_quic

package inbound

import (
	"context"
	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
)

type Hysteria2M struct {
	userPasswordList []string
	*Hysteria2
}

func NewHysteria2M(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.Hysteria2InboundOptions) (*Hysteria2M, error) {
	userPasswordList := make([]string, 0, len(options.Users))
	for _, user := range options.Users {
		userPasswordList = append(userPasswordList, user.Password)
	}
	h, err := NewHysteria2(ctx, router, logger, tag, options)
	if err != nil {
		return nil, err
	}
	return &Hysteria2M{
		userPasswordList: userPasswordList,
		Hysteria2:        h,
	}, nil
}

func (h *Hysteria2M) AddUsers(users []option.Hysteria2User) error {
	indexs := make([]int, 0, len(users)+len(h.userNameList))
	names := make([]string, len(users)+len(h.userNameList))
	pws := make([]string, len(users)+len(h.userPasswordList))
	for i := range users {
		names[i] = users[i].Name
		pws[i] = users[i].Password
	}
	if cap(h.userNameList)-len(h.userNameList) >= len(users) {
		h.userNameList = append(h.userNameList, names...)
		h.userPasswordList = append(h.userPasswordList, pws...)
	} else {
		tmp := make([]string, 0, len(h.userNameList)+len(users)+10)
		tmp = append(tmp, h.userNameList...)
		tmp = append(tmp, pws...)
		h.userNameList = tmp
		tmp = make([]string, 0, len(h.userPasswordList)+len(users)+10)
		tmp = append(tmp, h.userPasswordList...)
		tmp = append(tmp, pws...)
	}
	for i := range h.userNameList {
		indexs = append(indexs, i)
	}
	h.service.UpdateUsers(indexs, h.userPasswordList)
	return nil
}

func (h *Hysteria2M) DelUsers(name []string) error {
	if len(name) == 0 {
		return nil
	}
	is := make([]int, 0, len(name))
	ulen := len(name)
	for i := range h.userNameList {
		for _, u := range name {
			if h.userNameList[i] == u {
				is = append(is, i)
				ulen--
			}
			if ulen == 0 {
				break
			}
		}
	}
	ulen = len(h.userNameList)
	for _, i := range is {
		h.userNameList[i] = h.userNameList[ulen-1]
		h.userNameList[ulen-1] = ""
		h.userNameList = h.userNameList[:ulen-1]

		h.userPasswordList[i] = h.userPasswordList[ulen-1]
		h.userPasswordList[ulen-1] = ""
		h.userPasswordList = h.userPasswordList[:ulen-1]
		ulen--
	}
	indexs := make([]int, 0, len(h.userNameList))
	for i := range h.userNameList {
		indexs = append(indexs, i)
	}
	h.service.UpdateUsers(indexs, h.userPasswordList)
	return nil
}
