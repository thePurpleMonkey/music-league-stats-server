package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thePurpleMonkey/music-league-stats-server/models"
)

func main() {
	err := models.ConnectDatabase()
	checkErr(err)

	router := gin.Default()

	group := router.Group("/v1")
	{
		group.GET("leagues", getLeagues)
		// group.GET("leagues/:league_id", getLeagueById)
		group.GET("leagues/:league_id/rounds", getRounds)
		group.GET("leagues/:league_id/members", getMembers)
		group.GET("leagues/:league_id/members/:member_id/votes_received", getVotesReceived)
		group.GET("leagues/:league_id/members/:member_id/votes_given", getVotesGiven)
		group.GET("leagues/:league_id/members/:member_id/round_standings", getRoundStandings)
		group.GET("leagues/:league_id/members/:member_id/favorite_songs", getFavoriteSongs)

		group.GET("submissions/:round_id", getSubmissions)
		group.GET("voters/:round_id", getVotesByVoter)
		group.GET("members/:member_id", getMember)
		group.GET("rounds/:round_id", getRound)
		group.GET("rounds/:round_id/rankings", getRoundRankings)
	}

	router.Run("localhost:4040")
}

func getLeagues(c *gin.Context) {
	leagues, err := models.GetLeagues()
	checkErr(err)

	if leagues == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, leagues)
	}
}

func getRounds(c *gin.Context) {
	leagueId := c.Param("league_id")
	rounds, err := models.GetRounds(leagueId)
	checkErr(err)

	if rounds == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, rounds)
	}
}

func getMembers(c *gin.Context) {
	leagueId := c.Param("league_id")
	members, err := models.GetMembers(leagueId)
	checkErr(err)

	if members == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, members)
	}
}

func getRound(c *gin.Context) {
	roundId := c.Param("round_id")
	round, err := models.GetRoundById(roundId)
	checkErr(err)

	if round.Id == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, round)
	}
}

func getRoundRankings(c *gin.Context) {
	roundId := c.Param("round_id")
	rankings, err := models.GetRoundRankings(roundId)
	checkErr(err)

	if rankings == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, rankings)
	}
}

func getMember(c *gin.Context) {
	memberId := c.Param("member_id")
	member, err := models.GetMemberById(memberId)
	checkErr(err)

	if member.Id == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, member)
	}
}

func getVotesReceived(c *gin.Context) {
	leagueId := c.Param("league_id")
	memberId := c.Param("member_id")
	votes, err := models.GetVotesReceived(leagueId, memberId)
	checkErr(err)

	if votes == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, votes)
	}
}

func getVotesGiven(c *gin.Context) {
	leagueId := c.Param("league_id")
	memberId := c.Param("member_id")
	votes, err := models.GetVotesGiven(leagueId, memberId)
	checkErr(err)

	if votes == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, votes)
	}
}

func getRoundStandings(c *gin.Context) {
	leagueId := c.Param("league_id")
	memberId := c.Param("member_id")
	votes, err := models.GetRoundStandings(leagueId, memberId)
	checkErr(err)

	if votes == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, votes)
	}
}

func getFavoriteSongs(c *gin.Context) {
	leagueId := c.Param("league_id")
	memberId := c.Param("member_id")
	votes, err := models.GetFavoriteSongs(leagueId, memberId)
	checkErr(err)

	if votes == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, votes)
	}
}

func getSubmissions(c *gin.Context) {
	roundId := c.Param("round_id")
	round, err := models.GetSubmissions(roundId)
	checkErr(err)

	if round == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, round)
	}
}

func getVotesByVoter(c *gin.Context) {
	roundId := c.Param("round_id")
	round, err := models.GetVotesByVoter(roundId)
	checkErr(err)

	if round == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, round)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
