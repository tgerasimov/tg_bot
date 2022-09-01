package datastruct

type Vote struct {
	UserID    int
	OfferText string
	IsVoted   bool
}

type VoteOffer struct {
	OfferText   string
	VotesCount  int
	OfferNumber int
}
