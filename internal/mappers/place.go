package mappers

import (
	"database/sql"
	"sort"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ToPlaceResponse(req models.Place) dto.PlaceResponse {
	res := dto.PlaceResponse{
		ID:        req.ID,
		Title:     req.Title,
		CreatedAt: req.CreatedAt.Format(defaultDateTimeLayout),
		UpdatedAt: req.UpdatedAt.Format(defaultDateTimeLayout),
	}

	if req.ParentID.Valid {
		parentID := uint64(req.ParentID.Int64)
		res.ParentID = &parentID
	}

	return res
}

func ToAddPlaceRequest(req dto.AddPlaceRequest) models.AddPlaceRequest {
	res := models.AddPlaceRequest{
		Title: req.Title,
	}

	if req.ParentID != nil {
		res.ParentID = sql.NullInt64{Int64: int64(*req.ParentID), Valid: true}
	}

	return res
}

func ToUpdatePlaceRequest(id uint64, req dto.UpdatePlaceRequest) models.UpdatePlaceRequest {
	res := models.UpdatePlaceRequest{
		ID:    id,
		Title: req.Title,
	}

	if req.ParentID != nil {
		res.ParentID = sql.NullInt64{Int64: int64(*req.ParentID), Valid: true}
	}

	return res
}

func ToPlacesResponse(req []models.Place) dto.PlacesResponse {
	res := make([]dto.PlaceResponse, 0, len(req))

	for _, p := range req {
		res = append(res, ToPlaceResponse(p))
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Title < res[j].Title
	})

	return dto.PlacesResponse{
		Places: res,
	}
}

func ToPlaceTreeResponse(req []models.Place) dto.PlaceTreeResponse {
	var res []dto.PlaceTree
	parentMap := make(map[uint64][]models.Place, len(req))

	for _, p := range req {
		if !p.ParentID.Valid {
			res = append(res, dto.PlaceTree{
				Place: ToPlaceResponse(p),
			})
			continue
		}

		parentID := uint64(p.ParentID.Int64)
		parentMap[parentID] = append(parentMap[parentID], p)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Place.Title < res[j].Place.Title
	})

	for i := range res {
		recursiveFillNested(&res[i], parentMap)
	}

	return dto.PlaceTreeResponse{
		Places: res,
	}
}

func recursiveFillNested(item *dto.PlaceTree, parentMap map[uint64][]models.Place) {
	id := item.Place.ID

	if _, ok := parentMap[id]; ok {
		for _, childPlace := range parentMap[id] {
			nestedItem := dto.PlaceTree{
				Place: ToPlaceResponse(childPlace),
			}

			recursiveFillNested(&nestedItem, parentMap)

			item.NestedPlaces = append(item.NestedPlaces, nestedItem)
		}

		sort.Slice(item.NestedPlaces, func(i, j int) bool {
			return item.NestedPlaces[i].Place.Title < item.NestedPlaces[j].Place.Title
		})
	}
}
