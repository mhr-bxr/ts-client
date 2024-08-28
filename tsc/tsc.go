package tsc

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const (
	MinimumTip          = 500000000 // 0.5 TON
	BloXrouteTipAddress = "UQAw0AJjHbMYQobYXHBoW29ShKx1V2UjaiKanhDYBNJYDPUh"
)

type TransferParams struct {
	Amount  tlb.Coins
	Bounce  bool
	Comment string
}

// NewBundle returns a hex encoding of the requested transfers (and a tip).
func NewBundle(ctx context.Context, w *wallet.Wallet, transfers map[*address.Address]TransferParams, tip tlb.Coins) ([]byte, error) {
	if w == nil {
		return nil, fmt.Errorf("nil wallet")
	}
	if len(transfers) < 1 {
		return nil, nil
	}

	var result []*wallet.Message
	// build transfers from the data passed
	for k, v := range transfers {
		transfer, err := w.BuildTransfer(k, v.Amount, v.Bounce, v.Comment)
		if err != nil {
			return nil, fmt.Errorf("failed to build transfer: %w", err)
		}
		result = append(result, transfer)
	}
	// add the tip transfer
	stip := tip.Nano().Int64()
	if stip < MinimumTip {
		tip = tlb.FromNanoTON(big.NewInt(MinimumTip))
	}
	addr, err := address.ParseAddr(BloXrouteTipAddress)
	if err != nil {
		return nil, fmt.Errorf("internal error: failed to parse BloXroute tip address : %w", err)
	}
	tipTransfer, err := w.BuildTransfer(addr, tip, false, "BloXroute tip")
	if err != nil {
		return nil, fmt.Errorf("error: failed to build BloXroute tip transfer : %w", err)
	}
	result = append(result, tipTransfer)
	ext, err := w.PrepareExternalMessageForMany(ctx, false, result)
	if err != nil {
		return nil, fmt.Errorf("error: failed to prepare external message: %w", err)
	}
	msgCell, err := tlb.ToCell(ext)
	if err != nil {
		return nil, fmt.Errorf("error: failed to convert message to cell: %w", err)
	}
	src := msgCell.ToBOC()
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst, nil
}
