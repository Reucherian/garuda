package raft

import (
	"context"
	"errors"
)

func (r *Raft) GetUserHelper(username string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	resp, err := r.cli.Get(ctx, username)
	cancel()
	if err != nil {
		return []byte{}, err
	}
	if len(resp.Kvs) == 0 {
		return []byte{}, errors.New("user not found")
	}
	return resp.Kvs[0].Value, nil
}
