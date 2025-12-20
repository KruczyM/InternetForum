package models

type PostView struct {
	Post		
	AuthorName		string
	LikeCount		int
	CommentCount	int
	FormattedDate	string
}

type PageData struct {
	Posts	[]PostView
}