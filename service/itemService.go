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

func GetAllItems(logger *logrus.Entry) (*[]types.Items, error) {
	allItems, err := models.GetAllItems()
	if err != nil {
		logger.Errorf("GetAllItems: Failed DB Query: %s", err.Error())
		return nil, types.ErrGetAllItemsDB
	}
	if allItems != nil {
		logger.Infof("Got %v items", len(*allItems))
	} else {
		logger.Info("Got 0 items")
	}
	return itemDBModelToResponse(allItems), nil
}

func GetReservedItems(userID int, logger *logrus.Entry) (*[]types.Items, error) {
	reservedItems, err := models.GetReservedItems(userID)
	if err != nil {
		logger.Errorf("GetReservedItems: Failed DB Query: %s", err.Error())
		return nil, types.ErrGetReservedItemsDB
	}
	if reservedItems != nil {
		logger.Infof("Got %v items", len(*reservedItems))
	} else {
		logger.Info("Got 0 items")
	}
	return itemDBModelToResponse(reservedItems), nil
}

func EditItemRank(itemID int, rank int, logger *logrus.Entry) error {
	item, err := models.FindOneItem(map[string]interface{}{"id": itemID})
	if err != nil {
		logger.Errorf("EditItemRank: Failed DB Query to find item: %s", err.Error())
		return types.ErrEditItem
	}
	if item == nil {
		logger.Error("EditItemRank: Item Doesn't Exist")
		return types.ErrEditItem
	}
	updateErr := models.UpdateItemWithMap(item, map[string]interface{}{"rank": rank})
	if updateErr != nil {
		logger.Errorf("EditItemRank: Failed DB Query to update item: %s", err.Error())
		return types.ErrEditItem
	}
	return nil
}

func ReserveItem(userID, itemID int, logger *logrus.Entry) error {
	item, err := models.FindOneItem(map[string]interface{}{"id": itemID, "reserved": false}, models.ItemReserveScope)
	if err != nil {
		logger.Errorf("ReserveItem: Failed DB Query to find item: %s", err.Error())
		return types.ErrEditItem
	}
	if item == nil {
		logger.Error("ReserveItem: Item Doesn't Exist, or already reserved")
		return types.ErrEditItem
	}
	updateErr := models.UpdateItemWithMap(item, map[string]interface{}{"reserverid": userID, "reserved": true})
	if updateErr != nil {
		logger.Errorf("ReserveItem: Failed DB Query to update item: %s", err.Error())
		return types.ErrEditItem
	}
	return nil
}

func UnReserveItem(userID, itemID int, logger *logrus.Entry) error {
	item, err := models.FindOneItem(map[string]interface{}{"id": itemID}, models.ItemReserveScope)
	if err != nil {
		logger.Errorf("UnReserveItem: Failed DB Query to find item: %s", err.Error())
		return types.ErrEditItem
	}
	if item == nil {
		logger.Error("UnReserveItem: Item Doesn't Exist")
		return types.ErrEditItem
	}
	if item.ReserverID == nil || *item.ReserverID != userID {
		logger.Errorf("UnReserveItem: Failed, not original reserver or admin")
		return types.ErrEditItem
	}
	updateErr := models.UpdateItemWithMap(item, map[string]interface{}{"reserverid": nil, "reserved": false})
	if updateErr != nil {
		logger.Errorf("UnReserveItem: Failed DB Query to update item: %s", err.Error())
		return types.ErrEditItem
	}
	return nil
}

func AddItem(name, url string, rank int, logger *logrus.Entry) error {
	if err := models.AddItem(name, url, rank); err != nil {
		logger.Errorf("Failed to Add Item: %v, Error: %v", name, err.Error())
		return types.ErrAddItemErr
	}
	return nil
}

func DeleteItem(itemID int, logger *logrus.Entry) error {
	item, err := models.FindOneItem(map[string]interface{}{"id": itemID})
	if err != nil {
		logger.Errorf("DeleteItem: Failed DB Query to find item: %s", err.Error())
		return types.ErrDeleteItem
	}
	if item == nil {
		logger.Error("DeleteItem: Item Doesn't Exist")
		return types.ErrDeleteItem
	}

	if deleteErr := models.DeleteItem(item); deleteErr != nil {
		logger.Errorf("Failed to Delete Item: %v", item.Name)
		return types.ErrDeleteItem
	}

	return nil
}

func itemDBModelToResponse(items *[]models.ItemModel) *[]types.Items {
	var resp []types.Items
	for _, item := range *items {
		resp = append(resp, types.Items{Name: item.Name, Rank: item.Rank, Url: item.URL})
	}
	return &resp
}
