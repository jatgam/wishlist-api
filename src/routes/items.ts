import {Router, Request, Response, NextFunction} from "express";
import {check, validationResult} from "express-validator";

import {GetAllItems, GetWantedItems, AddItem, DeleteItem, EditItemRank, ReserveItem, UnReserveItem, GetReservedItems} from "../service/itemService";

export const items = Router();

items.get("", async (req: Request, res: Response, next: NextFunction) => {
    await GetWantedItems().then(items => {
        if (items.success) {
            res.json(items);
        } else {
            res.status(500).json(items);
        }
    }).catch(err => next(err));
});

items.get("/all", async (req: Request, res: Response, next: NextFunction) => {
    if (req.auth.level != 9) return res.sendStatus(401);
    await GetAllItems().then(items => {
        if (items.success) {
            res.json(items);
        } else {
            res.status(500).json(items);
        }
    }).catch(err => next(err));
});

items.get("/reserved", async (req: Request, res: Response, next: NextFunction) => {
    await GetReservedItems(req.auth.id).then(items => {
        if (items.success) {
            res.json(items);
        } else {
            res.status(500).json(items);
        }
    }).catch(err => next(err));
});

items.post("", [
        check("name").trim().escape().isString(),
        check("url").trim().isURL(),
        check("rank").trim().isNumeric().toInt()
    ], async (req: Request, res: Response, next: NextFunction) => {
    if (req.auth.level != 9) return res.sendStatus(401);
    const validateErr = validationResult(req).formatWith(error => ({
        param: error.param,
        msg: error.msg
    }));
    if (!validateErr.isEmpty()) {
        return res.status(422).json({errors: validateErr.array()});
    }
    await AddItem(req.body.name, req.body.url, req.body.rank).then(item => {
        if (item.success) {
            res.json(item);
        } else {
            res.status(400).json(item);
        }
    }).catch(err => next(err));
});

items.delete("/id/:itemid", [
        check("itemid").trim().isNumeric().toInt()
    ], async (req: Request, res: Response, next: NextFunction) => {
    if (req.auth.level != 9) return res.sendStatus(401);
    const validateErr = validationResult(req).formatWith(error => ({
        param: error.param,
        msg: error.msg
    }));
    if (!validateErr.isEmpty()) {
        return res.status(422).json({errors: validateErr.array()});
    }
    await DeleteItem(parseInt(req.params.itemid)).then(item => {
        if (item.success) {
            res.json(item);
        } else {
            res.status(400).json(item);
        }
    }).catch(err => next(err));
});

items.post("/id/:itemid/reserve", [
        check("itemid").trim().isNumeric().toInt()
    ], async (req: Request, res: Response, next: NextFunction) => {
    const validateErr = validationResult(req).formatWith(error => ({
        param: error.param,
        msg: error.msg
    }));
    if (!validateErr.isEmpty()) {
        return res.status(422).json({errors: validateErr.array()});
    }
    await ReserveItem(parseInt(req.params.itemid), req.auth.id).then(item => {
        if (item.success) {
            res.json(item);
        } else {
            res.status(400).json(item);
        }
    });
});

items.post("/id/:itemid/unreserve", [
        check("itemid").trim().isNumeric().toInt()
    ], async (req: Request, res: Response, next: NextFunction) => {
    const validateErr = validationResult(req).formatWith(error => ({
        param: error.param,
        msg: error.msg
    }));
    if (!validateErr.isEmpty()) {
        return res.status(422).json({errors: validateErr.array()});
    }
    await UnReserveItem(parseInt(req.params.itemid), req.auth.id, req.auth.level).then(item => {
        if (item.success) {
            res.json(item);
        } else {
            res.status(400).json(item);
        }
    });
});

items.post("/id/:itemid/rank/:rank", [
        check("itemid").trim().isNumeric().toInt(),
        check("rank").trim().isNumeric().toInt()
    ], async (req: Request, res: Response, next: NextFunction) => {
    const validateErr = validationResult(req).formatWith(error => ({
        param: error.param,
        msg: error.msg
    }));
    if (!validateErr.isEmpty()) {
        return res.status(422).json({errors: validateErr.array()});
    }
    await EditItemRank(parseInt(req.params.itemid), parseInt(req.params.rank)).then(item => {
        if (item.success) {
            res.json(item);
        } else {
            res.status(400).json(item);
        }
    });
});
