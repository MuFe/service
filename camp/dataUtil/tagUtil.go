package dataUtil

import (
	pb "mufe_service/jsonRpc"
)

type Tag struct {
	Tag string `json:"title" `
	Id  int64  `json:"id" `
}

func ParseTag(list []*pb.TagData) []Tag {
	tags := make([]Tag, 0)
	for _, v := range list {
		tags = append(tags,Tag{
			Id:v.Id,
			Tag:v.Title,
		})
	}
	return tags
}
