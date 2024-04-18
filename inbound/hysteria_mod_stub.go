//go:build !with_quic

package inbound

func NewHysteriaM(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.HysteriaInboundOptions) (adapter.Inbound, error) {
	return nil, C.ErrQUICNotIncluded
}

func NewHysteria2M(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.Hysteria2InboundOptions) (adapter.Inbound, error) {
	return nil, C.ErrQUICNotIncluded
}
