package bandit

import (
	"math"
	"sync"
)

// UCB represents the upper confidence bound algorithm
type UCB struct {
	sync.RWMutex
	Counts  []int
	Rewards []float64
}

// Init will initialise the counts and rewards with the provided number of arms
func (b *UCB) Init(nArms int) error {
	b.Lock()
	defer b.Unlock()

	if nArms < 1 {
		return ErrInvalidArms
	}
	b.Counts = make([]int, nArms)
	b.Rewards = make([]float64, nArms)
	return nil
}

// SelectArm chooses an arm that exploits if the value is more than the epsilon
// threshold, and explore if the value is less than epsilon
func (b *UCB) SelectArm(probability float64) int {
	b.RLock()
	defer b.RUnlock()

	nArms := len(b.Counts)

	// Select unplayed arms
	for i := 0; i < nArms; i++ {
		if b.Counts[i] == 0 {
			return i
		}
	}

	totalCounts := sum(b.Counts...)
	ucbValues := make([]float64, nArms)

	for i := 0; i < nArms; i++ {
		count := b.Counts[i]
		reward := b.Rewards[i]
		bonus := math.Sqrt((2.0 * math.Log(float64(totalCounts))) / float64(count))
		ucbValues[i] = bonus + reward
	}

	return max(ucbValues...)
}

// Update will update an arm with some reward value,
// e.g. click = 1, no click = 0
func (b *UCB) Update(chosenArm int, reward float64) error {
	b.Lock()
	defer b.Unlock()

	if chosenArm < 0 || chosenArm >= len(b.Rewards) {
		return ErrArmsIndexOutOfRange
	}
	if reward < 0 {
		return ErrInvalidReward
	}

	b.Counts[chosenArm]++
	n := float64(b.Counts[chosenArm])

	oldRewards := b.Rewards[chosenArm]
	b.Rewards[chosenArm] = (oldRewards*(n-1) + reward) / n

	return nil
}

// GetCounts returns the counts
func (b *UCB) GetCounts() []int {
	b.RLock()
	defer b.RUnlock()

	sCopy := make([]int, len(b.Counts))
	copy(sCopy, b.Counts)
	return sCopy
}

// GetRewards returns the rewards
func (b *UCB) GetRewards() []float64 {
	b.RLock()
	defer b.RUnlock()

	sCopy := make([]float64, len(b.Rewards))
	copy(sCopy, b.Rewards)
	return sCopy
}

// NewUCB returns a pointer to the UCB struct
func NewUCB(counts []int, rewards []float64) (*UCB, error) {
	if len(counts) != len(rewards) {
		return nil, ErrInvalidLength
	}

	return &UCB{
		Rewards: rewards,
		Counts:  counts,
	}, nil
}
