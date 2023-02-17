package manager

import (
	"os"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
)
var (
	courseService   pb.CourseServiceClient
	chapterService   pb.ChapterServiceClient
	collectionService   pb.CollectionServiceClient
	searchService   pb.SearchServiceClient
)

func GetCourseService() pb.CourseServiceClient {
	if courseService == nil {
		rpc, _ := utils.GetRPCService("course_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		courseService = pb.NewCourseServiceClient(rpc)
	}
	return courseService
}

func GetChapterService() pb.ChapterServiceClient {
	if chapterService == nil {
		rpc, _ := utils.GetRPCService("course_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		chapterService = pb.NewChapterServiceClient(rpc)
	}
	return chapterService
}

func GetCollectionService() pb.CollectionServiceClient {
	if collectionService == nil {
		rpc, _ := utils.GetRPCService("course_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		collectionService = pb.NewCollectionServiceClient(rpc)
	}
	return collectionService
}


func GetSearchService() pb.SearchServiceClient {
	if searchService == nil {
		rpc, _ := utils.GetRPCService("course_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		searchService = pb.NewSearchServiceClient(rpc)
	}
	return searchService
}
