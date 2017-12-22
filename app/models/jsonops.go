package models

type JSONUploadOp struct {
	AuthToken string
	Upload PlogData
}

type JSONUploadResult struct {
	Result bool
	IdPlog int
	ErrorMessage string
}
