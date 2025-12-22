package models

type PostView struct {
	Post		
	AuthorName		string
	LikeCount		int
	CommentCount	int
	FormattedDate	string
	Comments		[]Comment
}

type PageData struct {
	Posts	[]PostView
}