package models

type Exam struct {
    ID          uint   `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Duration    int    `json:"duration"`
    CreatedAt   time.Time
}