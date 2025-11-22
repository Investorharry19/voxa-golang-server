package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	OwnerUsername string             `json:"ownerUsername"`
	MessageText   string             `json:"messageText,omitempty"`
	AudioUrl      string             `json:"audioUrl,omitempty"`
	IsOpened      bool               `json:"isOpened"`
	PublicId      string             `json:"publicId,omitempty"`
	IsStarred     bool               `json:"isStarred"`
	CreatedAt     time.Time          `json:"createdAt"`
}

type TextMessageRequestDTO struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	OwnerUsername string             `json:"ownerUsername"`
	MessageText   string             `json:"messageText,omitempty"`
	Type          string             `json:"type"`
	IsOpened      bool               `json:"isOpened"`
	CreatedAt     time.Time          `json:"createdAt"`
	IsStarred     bool               `json:"isStarred"`
}
type AudioMessageRequestDTO struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	OwnerUsername string             `json:"ownerUsername"`
	Type          string             `json:"type"`
	IsOpened      bool               `json:"isOpened"`
	CreatedAt     time.Time          `json:"createdAt"`
	IsStarred     bool               `json:"isStarred"`
	AudioUrl      string             `json:"audioUrl,omitempty"`
	PublicId      string             `json:"publicId,omitempty"`
}

type MessageMarkAsRead struct {
	State *bool `json:"isStarred"`
}

/*

import mongoose, { Schema, mongo } from "mongoose";

const schema = new Schema(
  {
    ownerUsername: {
      type: String,
      required: true,
      lowercase: true,
    },
    type: {
      type: String,
      required: true,
    },
    messageText: {
      type: String,
    },
    audioUrl: {
      type: String,
    },
    isOpened: {
      type: Boolean,
      default: false,
    },
    publicId: {
      type: String,
      default: "",
    },
    isStarred: {
      type: Boolean,
      default: false,
    },
  },
  { timestamps: true }
);




*/
