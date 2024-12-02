package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"notificationService/domain"

	"github.com/gocql/gocql"
)

type NotificationCassandraStore struct {
	session *gocql.Session
	logger  *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment and creates a keyspace
func New(logger *log.Logger) (*NotificationCassandraStore, error) {
	db := os.Getenv("CASS_DB")

	logger.Println("Database Host: ", db)

	// Connect to default keyspace
	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Create 'student' keyspace
	err = session.Query(
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
					WITH replication = {
						'class' : 'SimpleStrategy',
						'replication_factor' : %d
					}`, "notification", 1)).Exec()
	if err != nil {
		logger.Println(err)
	}
	session.Close()

	// Connect to student keyspace
	cluster.Keyspace = "notification"
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &NotificationCassandraStore{
		session: session,
		logger:  logger,
	}, nil
}

// Disconnect from database
func (sr *NotificationCassandraStore) CloseSession() {
	sr.session.Close()
}

// Create tables
func (sr *NotificationCassandraStore) CreateTables() {
	// Create notifications table
	err := sr.session.Query(`
		CREATE TABLE IF NOT EXISTS notifications (
			user_id UUID,
			id UUID,
			message TEXT,
			status TEXT,
			created_at TIMESTAMP,
			PRIMARY KEY (user_id, created_at, id)
		)
	`).Exec()
	if err != nil {
		sr.logger.Println("Error creating notifications table:", err)

	}
}

// Disconnects the session
func (store *NotificationCassandraStore) Close() {
	store.session.Close()
}

// Inserts a new notification
func (store *NotificationCassandraStore) Insert(notification *domain.Notification) error {

	id := gocql.TimeUUID()
	err := store.session.Query(`
		INSERT INTO notifications (id, user_id, message, status, created_at) 
		VALUES (?, ?, ?, ?, ?)`,
		id, notification.UserID, notification.Message, notification.Status, time.Now(),
	).Exec()
	if err != nil {
		store.logger.Println("Error inserting notification:", err)
		return err
	}
	notification.ID = id
	return nil
}

func (store *NotificationCassandraStore) GetAll() ([]*domain.Notification, error) {
	iter := store.session.Query(`
		SELECT id, user_id, message, status, created_at 
		FROM notifications`).Iter()

	var notifications []*domain.Notification
	for {
		var notification domain.Notification
		if !iter.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.Status, &notification.CreatedAt) {
			break
		}
		// Append a copy of the notification object
		notifications = append(notifications, &notification)
	}
	if err := iter.Close(); err != nil {
		store.logger.Println("Error iterating notifications:", err)
		return nil, err
	}
	return notifications, nil
}

// Gets a notification by ID
func (store *NotificationCassandraStore) Get(id gocql.UUID) (*domain.Notification, error) {
	var notification domain.Notification
	err := store.session.Query(`
		SELECT id, user_id, message, status, created_at 
		FROM notifications WHERE id = ?`, id).
		Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.Status, &notification.CreatedAt)
	if err != nil {
		store.logger.Println("Error fetching notification:", err)
		return nil, err
	}
	return &notification, nil
}

func (store *NotificationCassandraStore) GetByUserId(userId gocql.UUID) ([]*domain.Notification, error) {
	iter := store.session.Query(`
		SELECT id, user_id, message, status, created_at 
		FROM notifications WHERE user_id = ? ALLOW FILTERING`, userId).Iter()

	var notifications []*domain.Notification
	for {
		var notification domain.Notification
		if !iter.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.Status, &notification.CreatedAt) {
			break
		}
		// Append a copy of the notification object
		notifications = append(notifications, &notification)
	}
	if err := iter.Close(); err != nil {
		store.logger.Println("Error iterating notifications:", err)
		return nil, err
	}
	return notifications, nil
}
func (repo *NotificationCassandraStore) SetStatusDiscarded(userID gocql.UUID, createdAt time.Time, id gocql.UUID) error {
	err := repo.session.Query(`
        UPDATE notifications 
        SET status = ? 
        WHERE user_id = ? AND created_at = ? AND id = ?`,
		"discarded", userID, createdAt, id).Exec()
	if err != nil {
		repo.logger.Println("Error updating notification status:", err)
		return err
	}
	return nil
}

// Deletes all notifications (not recommended for production use)
func (store *NotificationCassandraStore) DeleteAll() error {
	err := store.session.Query(`
		TRUNCATE notifications`).Exec()
	if err != nil {
		store.logger.Println("Error deleting all notifications:", err)
		return err
	}
	return nil
}
