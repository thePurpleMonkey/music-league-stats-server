package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./music_league.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

type League struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Round struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	TotalVotes int    `json:"total_votes"`
}

type Member struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Track struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	Submitter Member `json:"submitter"`
}

type Vote struct {
	Voter     Member `json:"voter"`
	Votes     int    `json:"votes"`
	Comment   string `json:"comment"`
	Track     Track  `json:"track"`
	Round     Round  `json:"round"`
	Placement int    `json:"placement"`
}

type Submission struct {
	Track     Track  `json:"track"`
	Submitter Member `json:"submitter"`
	Votes     []Vote `json:"votes"`
	Comment   string `json:"comment"`
}

type VotesGiven struct {
	Voter Member `json:"voter"`
	Votes []Vote `json:"votes"`
}

type Placement struct {
	Member    Member `json:"member"`
	Votes     int    `json:"votes"`
	Placement int    `json:"placement"`
	Round     Round  `json:"round"`
}

func GetLeagues() ([]League, error) {
	rows, err := DB.Query("SELECT id, name FROM leagues")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	leagues := make([]League, 0)

	for rows.Next() {
		league := League{}
		err = rows.Scan(&league.Id, &league.Name)

		if err != nil {
			return nil, err
		}

		leagues = append(leagues, league)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return leagues, err
}

func GetLeagueById(id string) (League, error) {
	stmt, err := DB.Prepare("SELECT id, name FROM leagues WHERE id = ?")
	if err != nil {
		return League{}, err
	}

	league := League{}

	sqlErr := stmt.QueryRow(id).Scan(&league.Id, &league.Name)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return League{}, nil
		}
		return League{}, sqlErr
	}

	return league, nil
}

func GetRounds(leagueId string) ([]Round, error) {
	rows, err := DB.Query("SELECT round_id, name, SUM(votes) FROM results JOIN rounds ON results.round_id = rounds.id WHERE league_id = ? GROUP BY round_id ORDER BY sequence", leagueId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rounds := make([]Round, 0)

	for rows.Next() {
		round := Round{}
		err = rows.Scan(&round.Id, &round.Name, &round.TotalVotes)

		if err != nil {
			return nil, err
		}

		rounds = append(rounds, round)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return rounds, err
}

func GetAllMembers() ([]Member, error) {
	rows, err := DB.Query("SELECT id, name, picture FROM members")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	members := make([]Member, 0)

	for rows.Next() {
		member := Member{}
		if err = rows.Scan(&member.Id, &member.Name, &member.Picture); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, err
}

func GetRoundMembers(roundId string) ([]Member, error) {
	rows, err := DB.Query("SELECT id, name, picture FROM members WHERE id IN (SELECT DISTINCT recipient_id FROM results WHERE round_id = ?)", roundId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	members := make([]Member, 0)

	for rows.Next() {
		member := Member{}
		if err = rows.Scan(&member.Id, &member.Name, &member.Picture); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, err
}

func GetMembers(leagueId string) ([]Member, error) {
	rows, err := DB.Query("SELECT id, name, picture FROM members WHERE id IN (SELECT DISTINCT recipient_id FROM results WHERE league_id = ?)", leagueId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	members := make([]Member, 0)

	for rows.Next() {
		member := Member{}
		if err = rows.Scan(&member.Id, &member.Name, &member.Picture); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, err
}

func GetVotesReceived(leagueId string, memberId string) ([]Vote, error) {
	rows, err := DB.Query("SELECT voter_id, SUM(votes) FROM results WHERE league_id = ? AND recipient_id = ? GROUP BY voter_id ORDER BY SUM(votes) DESC", leagueId, memberId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	votes := make([]Vote, 0)

	for rows.Next() {
		vote := Vote{}
		var voterId string
		if err = rows.Scan(&voterId, &vote.Votes); err != nil {
			return nil, err
		}
		if vote.Voter, err = GetMemberById(voterId); err != nil {
			return nil, err
		}

		votes = append(votes, vote)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return votes, err
}

func GetVotesGiven(leagueId string, memberId string) ([]Vote, error) {
	rows, err := DB.Query("SELECT recipient_id, SUM(votes) FROM results WHERE league_id = ? AND voter_id = ? GROUP BY recipient_id ORDER BY SUM(votes) DESC", leagueId, memberId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	votes := make([]Vote, 0)

	for rows.Next() {
		vote := Vote{}
		var voterId string
		if err = rows.Scan(&voterId, &vote.Votes); err != nil {
			return nil, err
		}
		if vote.Voter, err = GetMemberById(voterId); err != nil {
			return nil, err
		}

		votes = append(votes, vote)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return votes, err
}

func GetRoundStandings(leagueId string, memberId string) ([]Vote, error) {
	rows, err := DB.Query("SELECT round_id, SUM(votes) FROM results JOIN rounds ON results.round_id = rounds.id WHERE league_id = ? AND recipient_id = ? GROUP BY round_id ORDER BY sequence", leagueId, memberId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	votes := make([]Vote, 0)

	member, err := GetMemberById(memberId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		vote := Vote{
			Voter: member,
		}
		var roundId string
		if err = rows.Scan(&roundId, &vote.Votes); err != nil {
			return nil, err
		}

		if vote.Round, err = GetRoundById(roundId); err != nil {
			return nil, err
		}

		votes = append(votes, vote)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return votes, err
}

func GetRoundRankings(roundId string) ([]Placement, error) {
	rows, err := DB.Query("SELECT id, SUM(votes) FROM results JOIN members ON results.recipient_id = members.id WHERE round_id = ? GROUP BY recipient_id ORDER BY SUM(votes) DESC", roundId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ranking := make([]Placement, 0)
	var rank = 1
	for rows.Next() {
		placement := Placement{}
		var memberId string
		if err = rows.Scan(&memberId, &placement.Votes); err != nil {
			return nil, err
		}
		if placement.Member, err = GetMemberById(memberId); err != nil {
			return nil, err
		}
		if placement.Round, err = GetRoundById(roundId); err != nil {
			return nil, err
		}
		placement.Placement = rank

		ranking = append(ranking, placement)
		rank += 1
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return ranking, err
}

func GetFavoriteSongs(leagueId string, memberId string) ([]Vote, error) {
	rows, err := DB.Query("SELECT track_id, name, picture, votes, comment, recipient_id FROM results JOIN track_names ON results.track_id = track_names.id WHERE voter_id = ? AND league_id = ? ORDER BY votes DESC", memberId, leagueId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	votes := make([]Vote, 0)

	member, err := GetMemberById(memberId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		vote := Vote{
			Voter: member,
		}
		var submitterId string
		if err = rows.Scan(&vote.Track.Id, &vote.Track.Name, &vote.Track.Picture, &vote.Votes, &vote.Comment, &submitterId); err != nil {
			return nil, err
		}
		if vote.Track.Submitter, err = GetMemberById(submitterId); err != nil {
			return nil, err
		}

		votes = append(votes, vote)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return votes, err
}

func GetSubmissions(roundId string) ([]Submission, error) {
	rows, err := DB.Query("SELECT voter_id, recipient_id, votes, track_id, name, picture, comment FROM results JOIN track_names ON track_id = track_names.id WHERE round_id = ?", roundId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make(map[string]*Submission)

	for rows.Next() {
		submission := Submission{}
		track := Track{}
		vote := Vote{}

		var voterId, submitterId string
		if err = rows.Scan(&voterId, &submitterId, &vote.Votes, &track.Id, &track.Name, &track.Picture, &vote.Comment); err != nil {
			return nil, err
		}
		submission.Track = track
		if submission.Submitter, err = GetMemberById(submitterId); err != nil {
			return nil, err
		}
		if vote.Voter, err = GetMemberById(voterId); err != nil {
			return nil, err
		}
		vote.Track = track

		if sub, exists := submissions[submitterId]; exists {
			sub.Votes = append(sub.Votes, vote)
		} else {
			submission.Votes = []Vote{vote}
			submissions[submitterId] = &submission
		}
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	result := []Submission{}
	for _, sub := range submissions {
		result = append(result, *sub)
	}

	return result, err
}

func GetVotesByVoter(roundId string) ([]VotesGiven, error) {
	rows, err := DB.Query("SELECT voter_id, recipient_id, votes, track_id, name, picture, comment FROM results JOIN track_names ON track_id = track_names.id WHERE round_id = ?", roundId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	voteMap := make(map[string]*VotesGiven)

	for rows.Next() {
		votesGiven := VotesGiven{}
		track := Track{}
		vote := Vote{}

		var voterId, submitterId string
		if err = rows.Scan(&voterId, &submitterId, &vote.Votes, &track.Id, &track.Name, &track.Picture, &vote.Comment); err != nil {
			return nil, err
		}
		vote.Track = track
		if track.Submitter, err = GetMemberById(submitterId); err != nil {
			return nil, err
		}
		if votesGiven.Voter, err = GetMemberById(voterId); err != nil {
			return nil, err
		}
		vote.Voter = votesGiven.Voter
		vote.Track = track

		if votes, exists := voteMap[voterId]; exists {
			votes.Votes = append(votes.Votes, vote)
		} else {
			votesGiven.Votes = []Vote{vote}
			voteMap[voterId] = &votesGiven
		}
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	result := []VotesGiven{}
	for _, votes := range voteMap {
		result = append(result, *votes)
	}

	return result, err

}

func GetMemberById(memberId string) (Member, error) {
	stmt, err := DB.Prepare("SELECT id, name, picture FROM members WHERE id = ?")
	if err != nil {
		return Member{}, err
	}

	member := Member{}

	sqlErr := stmt.QueryRow(memberId).Scan(&member.Id, &member.Name, &member.Picture)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Member{}, nil
		}
		return Member{}, sqlErr
	}

	return member, nil
}

func GetRoundById(roundId string) (Round, error) {
	stmt, err := DB.Prepare("SELECT round_id, name, SUM(votes) FROM results JOIN rounds ON results.round_id = rounds.id WHERE round_id = ? ORDER BY sequence")
	if err != nil {
		return Round{}, err
	}

	round := Round{}
	err = stmt.QueryRow(roundId).Scan(&round.Id, &round.Name, &round.TotalVotes)
	if err != nil {
		if err == sql.ErrNoRows {
			return Round{}, nil
		}
		return Round{}, err
	}

	return round, nil
}

func GetSimilarity(roundId string, memberId string) (map[string]float32, error) {
	similarities := make(map[string]float32)
	votes := make(map[string][]string)

	rows, err := DB.Query("SELECT voter_id, track_id FROM results WHERE round_id = ? AND votes > 0", roundId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var voterId, trackId string
		if err = rows.Scan(&voterId, &trackId); err != nil {
			return nil, err
		}

		votes[voterId] = append(votes[voterId], trackId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	memberVotes := listToSet(votes[memberId])

	for voterId, votes := range votes {
		if voterId == memberId {
			continue // Don't calculate similarity with yourself
		}

		otherVotes := listToSet(votes)
		similarities[voterId] = calculateJaccardSimilarity(memberVotes, otherVotes)
	}

	return similarities, nil
}

func listToSet(list []string) map[string]bool {
	set := make(map[string]bool)

	for _, item := range list {
		set[item] = true
	}

	return set
}

// Calculate Jaccard similarity coefficient
func calculateJaccardSimilarity(votes1, votes2 map[string]bool) float32 {
	intersection := 0
	for item := range votes1 {
		if votes2[item] {
			intersection++
		}
	}

	union := len(votes1) + len(votes2) - intersection

	return float32(intersection) / float32(union)
}
