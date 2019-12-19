import {Router, Request, Response, NextFunction} from "express";
import {check, validationResult} from "express-validator";

import {User} from "../models/User";
import {Authenticate, Register, PasswordReset, PasswordResetTokenValid, PasswordForgot} from "../service/userService";

export const users = Router();

users.get("/get/:id", async (req, res, next) => {
    try {
        if (req.auth.level != 9) return res.sendStatus(401);
        const user = await User.findByPk(req.params["id"]);
        res.json(user);
    } catch (e) {
        next(e);
    }
});

users.get('/authenticated', async (req: Request, res: Response, next: NextFunction) => {
    try {
        if (req.auth.id) {
            return res.sendStatus(200);
        } else {
            return res.sendStatus(401);
        }
    } catch (e) {
        next(e);
    }
});

users.post("/auth", async (req, res, next) => {
    var username = req.body.username;
    var password = req.body.password;
    await Authenticate(username, password)
        .then(auth => {
            if (auth.success) {
                res.json(auth);
            } else {
                res.status(401).json(auth);
            }
        })
        .catch(err => next(err));
});

users.post("/register", [
        check('username', 'The username must be at least 3 characters long and may only contain letters and numbers').isString().isAlphanumeric().trim().escape().isLength({min: 3}),
        check('password', 'The password must be at least 10 character long and contain at least one of each: lowercase letter, uppercase letter, number, and special character').isLength({min: 10}).matches("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^\w\s]|[_])"),
        check('email', 'The E-Mail is not valid').trim().escape().isEmail(),
        check('firstname', 'The first name must be at least 1 character and letters only.').isLength({min: 1}).isString().isAlpha().trim().escape(),
        check('lastname', 'The last name must be at least 1 character and letters only.').isLength({min: 1}).isString().isAlpha().trim().escape()
    ], async (req: Request, res: Response, next: NextFunction) => {
    const validateErr = validationResult(req).formatWith(error => ({
        param: error.param,
        msg: error.msg
    }));
    if (!validateErr.isEmpty()) {
        return res.status(422).json({errors: validateErr.array()});
    }
    await Register(req.body.username, req.body.password, req.body.email, req.body.firstname, req.body.lastname)
        .then(reg => {
            if (reg.success) {
                res.json(reg);
            } else {
                res.status(400).json(reg);
            }
        }).catch(err => next(err));
});

users.post("/password_forgot", [check('email', 'The E-Mail is not valid').trim().escape().isEmail(),], async (req: Request, res: Response, next: NextFunction) => {
    const validateErr = validationResult(req).formatWith(error => ({
        param: error.param,
        msg: error.msg
    }));
    if (!validateErr.isEmpty()) {
        return res.status(422).json({errors: validateErr.array()});
    }
    if (!req.headers.host) return res.status(400);
    await PasswordForgot(req.body.email, req.headers.host)
        .then(pwForgot => {
            if (pwForgot.success){
                res.json(pwForgot);
            } else {
                res.status(400).json(pwForgot);
            }
        }).catch(err => next(err));
});

users.get("/password_reset/:token", [check('token').trim().escape().isString().isLength({min: 40,max:40})], async (req: Request, res: Response, next: NextFunction) => {
    const validateErr = validationResult(req).formatWith(error => ({
        param: error.param,
        msg: error.msg
    }));
    if (!validateErr.isEmpty()) {
        return res.status(422).json({errors: validateErr.array()});
    }
    await PasswordResetTokenValid(req.params.token)
        .then(tokenValid => {
            if (tokenValid) {
                res.json({success: true, message: 'Token Valid'});
            } else {
                res.status(400).json({success: false, message: 'Password reset token is invalid or expired.'});
            }
        }).catch(err => next(err));
});

users.post("/password_reset/:token", [
        check('token').trim().escape().isString().isLength({min: 40,max:40}),
        check('password', 'The password must be at least 10 character long and contain at least one of each: lowercase letter, uppercase letter, number, and special character').isLength({min: 10}).matches("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^\w\s]|[_])"),
        check('email', 'The E-Mail is not valid').trim().escape().isEmail()
    ], async (req: Request, res: Response, next: NextFunction) => {
    const validateErr = validationResult(req).formatWith(error => ({
        param: error.param,
        msg: error.msg
    }));
    if (!validateErr.isEmpty()) {
        return res.status(422).json({errors: validateErr.array()});
    }
    await PasswordReset(req.body.email, req.body.password, req.params.token)
        .then(passReset => {
            if (passReset.success) {
                res.json(passReset);
            } else{
                res.status(400).json(passReset);
            }
        }).catch(err => next(err));
});

users.get("/reserved/:id", async (req, res, next) => {
    try {
        const items = await User.scope('reserveditems').findByPk(req.params['id']);
        res.json(items);
    } catch (e) {
        next(e);
    }
});
