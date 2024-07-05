package postgres

import (
	"context"
	"fmt"
	"time"

	pb "sustainabilityService/genproto/SustainabilityService"

	"github.com/jmoiron/sqlx"
)

type SustainabilityRepo struct {
	db *sqlx.DB
}

func NewSustainabilityRepo(db *sqlx.DB) *SustainabilityRepo {
	return &SustainabilityRepo{db: db}
}

// 1
func (s *SustainabilityRepo) LogImpact(ctx context.Context, in *pb.LogImpactRequest) (*pb.LogImpactResponse, error) {
	query := `
		INSERT INTO impact_logs (
			user_id, 
			community_id,
			category, 
			amount, 
			unit)
		VALUES 
			($1, $2, $3, $4, $5)
		RETURNING 
			id, 
			user_id, 
			community_id,
			category, 
			amount, 
			unit, 
			logged_at
	`

	var (
		logResp   pb.LogImpactResponse
		logged_at time.Time
	)

	row := s.db.QueryRowxContext(ctx, query, in.UserId, in.CommunityId, in.Category, in.Amount, in.Unit)
	if err := row.Scan(
		&logResp.Id,
		&logResp.UserId,
		&logResp.CommunityId,
		&logResp.Category,
		&logResp.Amount,
		&logResp.Unit,
		&logged_at,
	); err != nil {
		return nil, err
	}

	logResp.LoggedAt = logged_at.Format("2006-01-02 15:04:05")

	return &logResp, nil
}

// 2
func (s *SustainabilityRepo) GetUserImpact(ctx context.Context, in *pb.GetUserImpactRequest) (*pb.GetUserImpactResponse, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			community_id,
			category, 
			amount, 
			unit, 
			logged_at
		FROM 
			impact_logs
		WHERE
			user_id = $1
	`
	rows, err := s.db.QueryxContext(ctx, query, in.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*pb.LogImpactResponse
	for rows.Next() {
		var logResp pb.LogImpactResponse
		var logged_at time.Time
		if err := rows.Scan(
			&logResp.Id,
			&logResp.UserId,
			&logResp.CommunityId,
			&logResp.Category,
			&logResp.Amount,
			&logResp.Unit,
			&logged_at,
		); err != nil {
			return nil, err
		}
		logResp.LoggedAt = logged_at.Format("2006-01-02 15:04:05")
		logs = append(logs, &logResp)
	}

	return &pb.GetUserImpactResponse{Impacts: logs}, nil
}

// 3
func (s *SustainabilityRepo) GetCommunityImpact(ctx context.Context, in *pb.GetCommunityImpactRequest) (*pb.GetCommunityImpactResponse, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			community_id,
			category, 
			amount, 
			unit, 
			logged_at
		FROM 
			impact_logs
		WHERE
			community_id = $1
	`
	rows, err := s.db.QueryxContext(ctx, query, in.CommunityId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*pb.LogImpactResponse
	for rows.Next() {
		var logResp pb.LogImpactResponse
		var logged_at time.Time
		if err := rows.Scan(
			&logResp.Id,
			&logResp.UserId,
			&logResp.CommunityId,
			&logResp.Category,
			&logResp.Amount,
			&logResp.Unit,
			&logged_at,
		); err != nil {
			return nil, err
		}
		logResp.LoggedAt = logged_at.Format("2006-01-02 15:04:05")
		logs = append(logs, &logResp)
	}

	return &pb.GetCommunityImpactResponse{Impacts: logs}, nil
}

// 4
func (s *SustainabilityRepo) GetChallenges(ctx context.Context, in *pb.GetChallengesRequest) (*pb.GetChallengesResponse, error) {
	query := `
		SELECT
			id,
			title,
			description,
			goal_amount,
			goal_unit,
			start_date,
			end_date,
		FROM
			sustainability_challenges
	`
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challanges []*pb.PostChallengesResponse
	for rows.Next() {
		var challange pb.PostChallengesResponse
		var start_date time.Time
		var end_date time.Time
		if err := rows.Scan(
			&challange.Id,
			&challange.Title,
			&challange.Description,
			&challange.GoalAmount,
			&challange.GoalUnit,
			&start_date,
			&end_date,
		); err != nil {
			return nil, err
		}
		challange.StartDate = start_date.Format("2006-01-02")
		challange.EndDate = end_date.Format("2006-01-02")
		challanges = append(challanges, &challange)
	}

	return &pb.GetChallengesResponse{Challenges: challanges}, nil
}

// 5
func (s *SustainabilityRepo) JoinChallenge(ctx context.Context, in *pb.JoinChallengeRequest) (*pb.JoinChallengeResponse, error) {
	query := `
		INSERT INTO user_challenges (
			user_id,
			community_id,
			challenge_id,
			progress,
			completed_at)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			user_id,
			community_id,
			challenge_id,
			progress,
			completed_at
	`
	var (
		joinResp     pb.JoinChallengeResponse
		completed_at time.Time
	)
	row := s.db.QueryRowxContext(ctx, query, in.UserId, in.CommunityId, in.ChallengeId, in.Progres, in.CompletedAt)
	if err := row.Scan(
		&joinResp.UserId,
		&joinResp.CommunityId,
		&joinResp.ChallengeId,
		&joinResp.Progres,
		&completed_at,
	); err != nil {
		return nil, err
	}
	joinResp.CompletedAt = completed_at.Format("2006-01-02 15:04:05")

	return &joinResp, nil
}

// 6
func (s *SustainabilityRepo) UpdateChallengeProgress(ctx context.Context, in *pb.UpdateChallengeProgressRequest) (*pb.UpdateChallengeProgressResponse, error) {
	query := `
		UPDATE 
			user_challenges
		SET 
			user_id = $1,
			community_id = $2,
			challenge_id = $3,
			progress = $4,
			completed_at = $5
		RETURNING
			user_id,
			community_id,
			challenge_id,
			progress,
			completed_at
	`
	var (
		updateResp     pb.UpdateChallengeProgressResponse
		completed_at time.Time
	)
	row := s.db.QueryRowxContext(ctx, query, in.UserId, in.CommunityId, in.ChallengeId, in.Progres, in.CompletedAt)
	if err := row.Scan(
		&updateResp.UserId,
		&updateResp.CommunityId,
		&updateResp.ChallengeId,
		&updateResp.Progres,
		&completed_at,
	); err != nil {
		return nil, err
	}
	updateResp.CompletedAt = completed_at.Format("2006-01-02 15:04:05")

	return &updateResp, nil
}

// 7
func (s *SustainabilityRepo) GetUserChallenges(ctx context.Context, in *pb.GetUserChallengesRequest) (*pb.GetUserChallengesResponse, error) {
	query := `
		SELECT
			user_id,
			community_id,
			challenge_id,
			progress,
			completed_at
		FROM
			user_challenges
		WHERE
			user_id = $1
	`
	rows, err := s.db.QueryxContext(ctx, query, in.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challanges []*pb.JoinChallengeResponse
	for rows.Next() {
		var challange pb.JoinChallengeResponse
		var completed_at time.Time
		if err := rows.Scan(
			&challange.UserId,
			&challange.CommunityId,
			&challange.ChallengeId,
			&challange.Progres,
			&completed_at,
		); err != nil {
			return nil, err
		}
		challange.CompletedAt = completed_at.Format("2006-01-02 15:04:05")
		challanges = append(challanges, &challange)
	}

	return &pb.GetUserChallengesResponse{Challenges: challanges}, nil
}

// 8
func (s *SustainabilityRepo) GetUserLeaderboard(ctx context.Context, in *pb.GetUserLeaderboardRequest) (*pb.GetUserLeaderboardResponse, error) {
	query := `
		SELECT user_id, community_id, SUM(progress) AS total_progress
		FROM user_challenges
		GROUP BY user_id, community_id
		ORDER BY total_progress DESC
		LIMIT 10;
	`

	fmt.Println()
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []*pb.LeaderboardUser
	for rows.Next() {
		var lead pb.LeaderboardUser
		if err := rows.Scan(
			&lead.UserId,
			&lead.CommunityId,
			&lead.Progres,
		); err != nil {
			return nil, err
		}
		leads = append(leads, &lead)
	}

	return &pb.GetUserLeaderboardResponse{Leaderboard: leads}, nil
}

// 9
func (s *SustainabilityRepo) GetCommunityLeaderboard(ctx context.Context, in *pb.GetCommunityLeaderboardRequest) (*pb.GetCommunityLeaderboardResponse, error) {
	query := `
		SELECT community_id, SUM(progress) AS total_progress
		FROM user_challenges
		GROUP BY community_id
		ORDER BY total_progress DESC
		LIMIT 10;
	`

	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []*pb.LeaderboardCommunity
	for rows.Next() {
		var lead pb.LeaderboardCommunity
		if err := rows.Scan(
			&lead.CommunityId,
			&lead.CommunityId,
			&lead.Progres,
		); err != nil {
			return nil, err
		}
		leads = append(leads, &lead)
	}

	return &pb.GetCommunityLeaderboardResponse{Leaderboard: leads}, nil
}

// 10
func (s *SustainabilityRepo) PostChallenges(ctx context.Context, in *pb.PostChallengesRequest) (*pb.PostChallengesResponse, error) {
	query := `
		INSERT INTO sustainability_challenges (
			title,
			description,
			goal_amount,
			goal_unit,
			start_date,
			end_date
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			title,
			description,
			goal_amount,
			goal_unit,
			start_date,
			end_date
	`

	row := s.db.QueryRowxContext(ctx, query, in.Title, in.Description, in.GoalAmount, in.GoalUnit, in.StartDate, in.EndDate)
	var (
		challange pb.PostChallengesResponse
		start_date time.Time
		end_date time.Time
	)

	if err := row.Scan(
		&challange.Id,
		&challange.Title,
		&challange.Description,
		&challange.GoalAmount,
		&challange.GoalUnit,
		&start_date,
		&end_date,
	); err != nil {
		return nil, err
	}

	challange.StartDate = start_date.Format("2006-01-02")
	challange.EndDate = end_date.Format("2006-01-02")

	return &challange, nil
}
