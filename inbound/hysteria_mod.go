//go:build with_quic

package inbound

import (
	"context"
	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
)

type HysteriaM struct {
	userPasswordList []string
	*Hysteria
}

func NewHysteriaM(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.HysteriaInboundOptions) (*HysteriaM, error) {
	userPasswordList := make([]string, 0, len(options.Users))
	for i := range options.Users {
		var password string
		if options.Users[i].AuthString != "" {
			password = options.Users[i].AuthString
		} else {
			password = string(options.Users[i].Auth)
		}
		userPasswordList[i] = password
	}
	h, err := NewHysteria(ctx, router, logger, tag, options)
	if err != nil {
		return nil, err
	}
	return &HysteriaM{
		userPasswordList: userPasswordList,
		Hysteria:         h,
	}, nil
}

func (h *HysteriaM) AddUsers(users []option.HysteriaUser) error {
	indexs := make([]int, 0, len(users)+len(h.userNameList))
	names := make([]string, len(users)+len(h.userNameList))
	pws := make([]string, len(users)+len(h.userPasswordList))
	for i := range users {
		var password string
		if users[i].AuthString != "" {
			password = users[i].AuthString
		} else {
			password = string(users[i].Auth)
		}
		names[i] = users[i].Name
		pws[i] = password
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

func (h *HysteriaM) DelUsers(name []string) error {
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
