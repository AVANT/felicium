package models

import (
	"github.com/AVANT/felicium/model"
)

///
//	These are the required methods
///

type Comment struct {
	*model.Model
}

type Comments []*Comment

type CommentSearchReturn struct {
	Partial *Comment
}

func NewComment() *Comment {
	return &Comment{
		Connection.NewModel("comment"),
	}
}

func NewCommentInterface() model.IsModel {
	return NewComment()
}

func CommentsFromModels(m model.Models) *Comments {
	u := Comments{}
	for _, v := range m {
		u = append(u, v.(*Comment))
	}
	return &u
}

func CommentFromModel(m model.Model) *Comment {
	return &Comment{}
}

func CommentFromIsModel(m model.IsModel) *Comment {
	return &Comment{}
}
