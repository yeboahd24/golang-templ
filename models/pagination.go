package models

type PaginationData struct {
    CurrentPage  int
    TotalPages   int
    HasPrevious  bool
    HasNext      bool
    PreviousPage int
    NextPage     int
}
