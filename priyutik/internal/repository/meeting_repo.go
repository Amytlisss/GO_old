package repository

import (
	"database/sql"
	"priyutik/internal/models"
	"time"
)

type MeetingRepo struct {
	db *sql.DB
}

func (r *MeetingRepo) CreateMeeting(userID int, date time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO meetings (user_id, date) VALUES ($1, $2)",
		userID, date,
	)
	return err
}

func (r *MeetingRepo) GetMeetings(userID int) ([]models.Meeting, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, date, cancelled, created_at FROM meetings WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var m models.Meeting
		if err := rows.Scan(&m.ID, &m.UserID, &m.Date, &m.Cancelled, &m.CreatedAt); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}
	return meetings, nil
}

func (r *MeetingRepo) CancelMeeting(id int) error {
	_, err := r.db.Exec(
		"UPDATE meetings SET cancelled = TRUE WHERE id = $1",
		id,
	)
	return err
}

func (r *MeetingRepo) UpdateMeeting(id int, newDate time.Time) error {
	_, err := r.db.Exec(
		"UPDATE meetings SET date = $1 WHERE id = $2",
		newDate, id,
	)
	return err
}

func (r *MeetingRepo) GetMeetingByID(id int) (*models.Meeting, error) {
	var m models.Meeting
	err := r.db.QueryRow(
		"SELECT id, user_id, date, cancelled FROM meetings WHERE id = $1",
		id,
	).Scan(&m.ID, &m.UserID, &m.Date, &m.Cancelled)
	return &m, err
}

func (r *MeetingRepo) GetAllMeetings() ([]models.Meeting, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, date, cancelled, created_at FROM meetings",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var m models.Meeting
		if err := rows.Scan(&m.ID, &m.UserID, &m.Date, &m.Cancelled, &m.CreatedAt); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}
	return meetings, nil
}

func (r *MeetingRepo) GetMeetingsByDate(date time.Time) ([]models.Meeting, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, date, cancelled, created_at FROM meetings WHERE date::date = $1",
		date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var m models.Meeting
		if err := rows.Scan(&m.ID, &m.UserID, &m.Date, &m.Cancelled, &m.CreatedAt); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}
	return meetings, nil
}
