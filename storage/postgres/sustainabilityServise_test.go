package postgres

import (
	"context"
	"fmt"
	"log"
	"testing"

	pb "sustainabilityService/genproto/SustainabilityService"

	"github.com/jmoiron/sqlx"
)

func ConnectDB() *SustainabilityRepo {
	db, err := sqlx.Open("postgres", "host=localhost user=postgres password=root port=5432 dbname=authentification sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return NewSustainabilityRepo(db)
}

func TestLogImpact(t *testing.T) {

	repo := ConnectDB()
	test := pb.LogImpactRequest{
		UserId:   "2a8c33b9-f3c0-4551-88e4-30e3957b331c",
		Category: "water_saved",
		Amount:   100.0,
		Unit:     "liters",
	}

	ctx := context.Background()
	resp, err := repo.LogImpact(ctx, &test)
	if err != nil {
		t.Fatalf("Error getting user: %v", err)
	}

	Wait := pb.LogImpactResponse{
		UserId:   "2a8c33b9-f3c0-4551-88e4-30e3957b331c",
		Category: "water_saved",
		Amount:   100.0,
		Unit:     "liters",
	}

	fmt.Println("response:  ", resp)
	fmt.Println("wait :   ", Wait)

	if Wait.UserId != resp.UserId || Wait.Category != resp.Category || Wait.Amount != resp.Amount || Wait.Unit != resp.Unit {
		t.Errorf("User data does not match. Wait: %+v, got: %+v", Wait, resp)
	}
}

func TestGetUserImpact(t *testing.T) {
	repo := ConnectDB()
	test := pb.GetUserImpactRequest{
		UserId: "2a8c33b9-f3c0-4551-88e4-30e3957b331c",
	}
	ctx := context.Background()
	resp, err := repo.GetUserImpact(ctx, &test)
	if err != nil {
		t.Fatalf("Error getting user: %v", err)
	}
	impacts := resp.Impacts
	wait := pb.LogImpactResponse{
		Id:       "2a8c33b9-f3c0-4551-88e4-30e3957b331a",
		UserId:   "2a8c33b9-f3c0-4551-88e4-30e3957b331c",
		Category: "water_saved",
		Amount:   100.0,
		Unit:     "liters",
	}

	if impacts[0].Id != wait.Id || impacts[0].Category != wait.Category || impacts[0].Amount != wait.Amount {

	}
}

func TestGetChallenges(t *testing.T) {

	repo := ConnectDB()
	test := pb.GetChallengesRequest{}

	ctx := context.Background()
	resp, err := repo.GetChallenges(ctx, &test)
	if err != nil {
		t.Fatalf("Error getting user: %v", err)
	}

	//update qilingan malumotni qo'ying
	Wait := pb.PostChallengesResponse{
		Id: "2a8c33b9-f3c0-4551-88e4-30e3957b331a",
		Title:       "Test Challenge",
        Description: "This is a test challenge",
        GoalAmount:  1000.0,
        GoalUnit:    "liters",
        StartDate:   "2022-01-01",
	}

	res := resp.Challenges
	fmt.Println("response:  ", resp)
	fmt.Println("wait :   ", Wait)

	if Wait.Id != res[0].Id || Wait.Title != res[0].Title || Wait.Description != res[0].Description {
		t.Errorf("User data does not match. Wait: %+v, got: %+v", Wait, res)
	}
}

func TestJoinChallenge(t *testing.T) {

	repo := ConnectDB()
	test := pb.JoinChallengeRequest{
		CommunityId:       "2a8c33b9-f3c0-4551-88e4-30e3957b331a",
		UserId:   "2a8c33b9-f3c0-4551-88e4-30e3957b331c",
		ChallengeId: "738c33b9-f3c0-4551-88e4-30e3957b331c",
		Progres: 34.6,
	}

	ctx := context.Background()
	resp, err := repo.JoinChallenge(ctx, &test)
	if err != nil {
		t.Fatalf("Error deleting user: %v", err)
	}

	Wait := pb.JoinChallengeResponse{
		UserId:   "2a8c33b9-f3c0-4551-88e4-30e3957b331c",
		CommunityId:       "2a8c33b9-f3c0-4551-88e4-30e3957b331a",
		ChallengeId: "738c33b9-f3c0-4551-88e4-30e3957b331c",
		Progres: 34.6,		
	}


	if Wait.UserId != resp.UserId || Wait.CommunityId != resp.CommunityId || Wait.ChallengeId != resp.CommunityId || Wait.Progres != resp.Progres {
		t.Errorf("User data does not match. Wait: %+v, got: %+v", Wait, resp)
	}
}

func TestUpdateChallengeProgress(t *testing.T) {
	repo := ConnectDB()
	test := pb.UpdateChallengeProgressRequest{
		UserId: "2a8c33b9-f3c0-4551-88e4-30e3957b331c",
		ChallengeId: "738c33b9-f3c0-4551-88e4-30e3957b331c",
		Progres: 34.6,		
	}

	ctx := context.Background()
	resp, err := repo.UpdateChallengeProgress(ctx, &test)
	if err != nil {
		t.Fatalf("Error getting user profile: %v", err)

	}

	Wait := pb.UpdateChallengeProgressResponse{
		UserId: "2a8c33b9-f3c0-4551-88e4-30e3957b331c",
		ChallengeId: "738c33b9-f3c0-4551-88e4-30e3957b331c",
		Progres: 34.6,	
	}

	if resp.UserId != Wait.UserId || resp.ChallengeId != Wait.ChallengeId || resp.Progres != Wait.Progres{
		t.Errorf("User profile data does not match. Wait: %+v, got: %+v", Wait, resp)
	}

}

func TestGetUserChallenges(t *testing.T) {
	repo := ConnectDB()
	test := pb.GetUserChallengesRequest{
		UserId: "fe1bedf7-84ba-4aa0-8b3a-f8e1d53e7c13",
	}

	ctx := context.Background()
	resp, err := repo.GetUserChallenges(ctx, &test)
	if err != nil {
		t.Fatalf("Error getting user profile: %v", err)

	}

	res := resp.Challenges

	Wait := pb.JoinChallengeResponse{
		UserId: "2a8c33b9-f3c0-4551-88e4-30e3957b331c",
		ChallengeId: "738c33b9-f3c0-4551-88e4-30e3957b331c",
		Progres: 34.6,	
	}

	if res[0].UserId != Wait.UserId || res[0].ChallengeId != Wait.ChallengeId || res[0].Progres != Wait.Progres {
		t.Errorf("User profile data does not match. Wait: %+v, got: %+v", Wait, resp)
	}

}
