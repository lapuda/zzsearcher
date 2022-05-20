package model

type Manga struct {
	MangaID    int    `json:"comic_id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	UpdateDate string `json:"update_date" bson:"update_date"`
	Author     string `json:"author" bson:"author"`
	Type       string `json:"tags" bson:"type"`
	Status     string `json:"status" bson:"status"`
	Describe   string `json:"describe" bson:"describe"`
	Cover      string `json:"cover_image" bson:"cover"`
	CreateTime int64  `json:"create_time" bson:"create_time"`
}

type Chapter struct {
	ChapterID   int    `json:"chapter_id" bson:"_id"`
	MangaID     int    `json:"comic_id" bson:"manga_id"`
	ChapterName string `json:"chapter_name" bson:"chapter_name"`
	Images      string `json:"chapter_images" bson:"images"`
	CreateTime  int    `json:"create_time" bson:"create_time"`
	UpdateTime  int64  `json:"update_time" bson:"update_time"`
}
