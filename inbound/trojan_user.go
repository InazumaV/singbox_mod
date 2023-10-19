package inbound

import (
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func (h *Trojan) AddUsers(users []option.TrojanUser) error {
	if cap(h.users)-len(h.users) >= len(users) {
		h.users = append(h.users, users...)
	} else {
		tmp := make([]option.TrojanUser, 0, len(h.users)+len(users)+10)
		tmp = append(tmp, h.users...)
		tmp = append(tmp, users...)
		h.users = tmp
	}
	err := h.service.UpdateUsers(common.MapIndexed(h.users, func(index int, user option.TrojanUser) int {
		return index
	}), common.Map(h.users, func(user option.TrojanUser) string {
		return user.Password
	}))
	if err != nil {
		return err
	}
	return nil
}

func (h *Trojan) DelUsers(names []string) error {
	is := make([]int, 0, len(names))
	ulen := len(names)
	for i := range h.users {
		for _, n := range names {
			if h.users[i].Name == n {
				is = append(is, i)
				ulen--
			}
			if ulen == 0 {
				break
			}
		}
	}
	ulen = len(h.users)
	for _, i := range is {
		h.users[i] = h.users[ulen-1]
		h.users[ulen-1] = option.TrojanUser{}
		h.users = h.users[:ulen-1]
		ulen--
	}
	err := h.service.UpdateUsers(common.MapIndexed(h.users, func(index int, user option.TrojanUser) int {
		return index
	}), common.Map(h.users, func(user option.TrojanUser) string {
		return user.Password
	}))
	return err
}
