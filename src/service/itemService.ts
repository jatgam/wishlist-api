import {Item} from "../models/Item";
import {GetItemsResponse, GenericResponse} from 'types';

export async function GetWantedItems(): Promise<GetItemsResponse> {
    return Item.findAll({where: {reserved: false}, order: [['rank', 'ASC'], ['id', 'DESC']]}).then(items => {
        return {success: true, items: items};
    }).catch(err => {
        return {success: false, items: []};
    });
}

export async function GetReservedItems(reserverid: number): Promise<GetItemsResponse> {
    return Item.findAll({where: {reserved: true, reserverid: reserverid}, order: [['rank', 'ASC'], ['id', 'DESC']]}).then(items => {
        return {success: true, items: items};
    }).catch(err => {
        return {success: false, items: []};
    });
}


// export async function WhoReserved(reserverid: number) {
//     // get first/last name
// }

export async function GetAllItems(): Promise<GetItemsResponse> {
    return Item.scope("user").findAll({order: [['rank', 'ASC'], ['id', 'DESC']]}).then(items => {
        return {success: true, items: items};
    }).catch(err => {
        return {success: false, items: []};
    });
}

export async function AddItem(name: string, url: string, rank: number): Promise<GenericResponse> {
    return Item.create({name: name, url: url, rank: rank}).then(item => {
        if (item) {
            return {success: true, message: "Item created."};
        } else {
            return {success: false, message: "Failed to create item."};
        }
    }).catch(err => {
        console.error(err);
        return {success: false, message: "Failed to create item."};
    })
}

export async function DeleteItem(itemid: number): Promise<GenericResponse> {
    return Item.findOne({where: {id: itemid}}).then(async item => {
        if (item) {
            return await item.destroy().then(itemDeleted => {
                return {success: true, message: "Item Deleted."};
            }).catch(err => {
                return {success: false, message: "Failed to delete item."};
            });
        } else {
            return {success: false, message: "Failed to delete item."};
        }
    })
}

export async function ReserveItem(itemid:number, reserverid: number): Promise<GenericResponse> {
    return Item.findOne({where: {id: itemid}}).then(item => {
        if (item) {
            item.reserved = true;
            item.reserverid = reserverid;
            return item.save().then(item => {
                if (item) {
                    return {success: true, message: "Item Reserved."};
                } else {
                    return {success: false, message: 'Failed to reserve the item.'};
                }

            }).catch(err => {
                console.error(err);
                return {success: false, message: 'Failed to reserve the item.'};
            })
        } else {
            return {success: false, message: "Couldn't find item to reserve."};
        }
    }).catch(err => {
        console.error(err);
        return {success: false, message: "Failed to reserve the item."};
    });
}

export async function UnReserveItem(itemid: number, requesterId: number, requesterUserLevel: number): Promise<GenericResponse> {
    return Item.scope("reserver").findOne({where: {id: itemid}}).then(item => {
        if (item) {
            if ((item.reserverid === requesterId) || (requesterUserLevel === 9)) {
                item.reserved = false;
                item.reserverid = null;
                return item.save().then(item => {
                    if (item) {
                        return {success: true, message: "Item UnReserved."};
                    } else {
                        return {success: false, message: 'Failed to UnReserve the item.'};
                    }

                }).catch(err => {
                    console.error(err);
                    return {success: false, message: 'Failed to UnReserve the item.'};
                })
            } else{
                return {success: false, message: 'Failed to UnReserve the item. You aren\'t the reserver or an admin'};
            }
            
        } else {
            return {success: false, message: "Couldn't find item to UnReserve."};
        }
    }).catch(err => {
        console.error(err);
        return {success: false, message: "Failed to UnReserve the item."};
    });
}

export async function EditItemRank(itemid:number, rank: number): Promise<GenericResponse> {
    return Item.findOne({where: {id: itemid}}).then(item => {
        if (item) {
            item.rank = rank;
            return item.save().then(item => {
                if (item) {
                    return {success: true, message: "Item Rank Updated."};
                } else {
                    return {success: false, message: 'Failed to update the Item Rank.'};
                }

            }).catch(err => {
                console.error(err);
                return {success: false, message: 'Failed to update the Item Rank.'};
            })
        } else {
            return {success: false, message: "Couldn't find item to update the rank.."};
        }
    }).catch(err => {
        console.error(err);
        return {success: false, message: "Failed to update the Item Rank."};
    });
}
