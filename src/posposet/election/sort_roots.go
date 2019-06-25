package election

import (
	"errors"
	"math/big"
	"sort"
)

type sortedRoot struct {
	stakeAmount *big.Int
	root        RootHash
}
type sortedRoots []sortedRoot

func (s sortedRoots) Len() int {
	return len(s)
}
func (s sortedRoots) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// compare by stake amount, root hash
func (s sortedRoots) Less(i, j int) bool {
	if s[i].stakeAmount.Cmp(s[j].stakeAmount) > 0 {
		return true
	}
	return s[i].root.Big().Cmp(s[j].root.Big()) < 0
}

// Chooses the famous witness with the greatest stake amount.
// This root serves as a "checkpoint" within DAG, as it's guaranteed to be final and consistent unless more than 1/3n are Byzantine.
// Other nodes will come to the same SfWitness not later than current highest frame + 2.
func (el *Election) chooseSfWitness() (*ElectionRes, error) {
	famousRoots := make(sortedRoots, 0, len(el.nodes))
	// fill famousRoots
	for slot, vote := range el.decidedRoots {
		if vote.yes {
			famousRoots = append(famousRoots, sortedRoot{
				root:        vote.seenRoot,
				stakeAmount: el.nodes[slot.nodeindex].StakeAmount,
			})
		}
	}
	if len(famousRoots) == 0 {
		return nil, errors.New("All the roots aren't famous, which is possible only if more than 1/3n are Byzantine")
	}

	// sort by stake amount, root hash
	sort.Sort(famousRoots)

	// take root with greatest stake
	return &ElectionRes{
		DecidedFrame:     el.frameToDecide,
		DecidedSfWitness: famousRoots[0].root,
	}, nil;
}