package service

import (
	"github.com/sirupsen/logrus"

	"github.com/jatgam/wishlist-api/models"
	"github.com/jatgam/wishlist-api/types"
)

func GetWantedItems(logger *logrus.Entry) (*[]types.Items, error) {
	wantedItems, err := models.GetWantedItems()
	if err != nil {
		logger.Errorf("GetWantedItems: Failed DB Query: %s", err.Error())
		return nil, types.ErrGetWantedItemsDB
	}
	if wantedItems != nil {
		logger.Infof("Got %v wanted items.", len(*wantedItems))
	} else {
		logger.Info("Got 0 wanted items.")
	}

	return itemDBModelToResponse(wantedItems), nil
}

func itemDBModelToResponse(items *[]models.ItemModel) *[]types.Items {
	var resp []types.Items
	for _, item := range *items {
		resp = append(resp, types.Items{Name: item.Name, Rank: item.Rank, Url: item.URL})
	}
	return &resp
}
