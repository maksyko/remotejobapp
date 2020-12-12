package models

import (
	"github.com/ievgen-ma/remotejobapp/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Post struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Path      string
	Name      string
	Host      string
	Title     string
	Salary    string
	Position  string
	Company   string
	Apply     string
	Processed bool
	Created   time.Time
	Updated   time.Time
}

type PostsHandler struct {
	posts *mgo.Collection
}

func NewPostsHandler() *PostsHandler {
	return &PostsHandler{
		posts: database.CreateConn().Posts,
	}
}

func (h *PostsHandler) FindPosts(limit int) ([]*Post, error) {
	var ps []*Post
	return ps, h.posts.Find(bson.M{
		"processed": false,
	}).Limit(limit).Sort("-created").All(&ps)
}

func (h *PostsHandler) Processed(ps []*Post) error {
	bulk := h.posts.Bulk()
	for _, p := range ps {
		bulk.UpdateAll(bson.M{"_id": p.ID}, bson.M{
			"$set": bson.M{"processed": true, "updated": time.Now()},
		})
	}
	_, err := bulk.Run()
	return err
}

func (h *PostsHandler) GetPostsCount(name, path string) (int, error) {
	return h.posts.Find(bson.M{
		"name": name,
		"path": path,
	}).Count()
}

func (h *PostsHandler) SavePost(post *Post) error {
	post.Created = time.Now()
	post.Processed = false
	return h.posts.Insert(post)
}
