package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=500"`
}

type Ranking struct {
	RankingValue int
	RankingName  string
}

type Movies struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ImdbID      string        `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title       string        `bson:"title" json:"title" validate:"required,min=2,max=500"`
	PosterPath  string        `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YouTubeID   string        `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre       []Genre       `bson:"genre" json:"genre" validate:"dive,required"`
	AdminReview string        `bson:"admin_review" json:"admin_review"`
	Ranking     Ranking       `bson:"ranking" json:"ranking" validate:"required"`
}
