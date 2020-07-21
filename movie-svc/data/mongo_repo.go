package data

import (
	"context"
	"time"

	"github.com/0x113/x-media/movie-svc/databases"
	"github.com/0x113/x-media/movie-svc/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionName = "movies"
)

// movieRepository manages the movie CRUD
type movieRepository struct{}

// NewMongoMovieRepository returns new instance of the movie repository
func NewMongoMovieRepository() MovieRepository {
	return &movieRepository{}
}

// Save movie to the database
func (r *movieRepository) Save(movie *models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionCopy := databases.Database.Session
	defer sessionCopy.EndSession(ctx)

	collection := sessionCopy.Client().Database(databases.Database.DbName).Collection(collectionName)

	_, err := collection.InsertOne(ctx, movie)
	if err != nil {
		return err
	}

	return nil
}

// Update existing movie
func (r *movieRepository) Update(movie *models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionCopy := databases.Database.Session
	defer sessionCopy.EndSession(ctx)

	collection := sessionCopy.Client().Database(databases.Database.DbName).Collection(collectionName)

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": movie.ID},
		bson.M{"$set": movie},
	)

	if err != nil {
		return err
	}

	return nil
}

// GetByTitle returns movie from the database based on its name if exists
func (r *movieRepository) GetByTitle(title string) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionCopy := databases.Database.Session
	defer sessionCopy.EndSession(ctx)

	collection := sessionCopy.Client().Database(databases.Database.DbName).Collection(collectionName)

	var movie models.Movie
	if err := collection.FindOne(ctx, bson.M{"title": title}).Decode(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

// GetByOriginalTitle returns movie from the database based on its name if exists
func (r *movieRepository) GetByOriginalTitle(title string) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionCopy := databases.Database.Session
	defer sessionCopy.EndSession(ctx)

	collection := sessionCopy.Client().Database(databases.Database.DbName).Collection(collectionName)

	var movie models.Movie
	if err := collection.FindOne(ctx, bson.M{"original_title": title}).Decode(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

// GetAll returns all movies from the database
func (r *movieRepository) GetAll() ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionCopy := databases.Database.Session
	defer sessionCopy.EndSession(ctx)

	collection := sessionCopy.Client().Database(databases.Database.DbName).Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var movies []*models.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

// GetByID returns movie from the database based on its id
func (r *movieRepository) GetByID(id primitive.ObjectID) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionCopy := databases.Database.Session
	defer sessionCopy.EndSession(ctx)

	collection := sessionCopy.Client().Database(databases.Database.DbName).Collection(collectionName)

	var movie models.Movie
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}
