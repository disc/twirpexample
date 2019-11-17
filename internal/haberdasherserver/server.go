package haberdasherserver

import (
	"context"
	"math/rand"

	pb "github.com/disc/twirpexample/rpc/haberdasher"
	"github.com/twitchtv/twirp"
)

// Server implements the Haberdasher internal
type Server struct{}

func (s *Server) MakeHat(ctx context.Context, size *pb.Size) (hat *pb.Hat, err error) {
	if size.Inches <= 0 {
		twerr := twirp.InvalidArgumentError("inches", "I can't make a hat that small!")
		twerr = twerr.WithMeta("retryable", "false")
		twerr = twerr.WithMeta("retry_after", "15s")
		return nil, twerr
	}

	return &pb.Hat{
		Inches: size.Inches,
		Color:  []pb.Color{pb.Color_RED, pb.Color_GREEN, pb.Color_BLUE}[rand.Intn(2)],
		Name:   []string{"bowler", "baseball cap", "top hat", "derby"}[rand.Intn(3)],
	}, nil
}
