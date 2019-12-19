import jwt from "jsonwebtoken";
import md5 from "md5";
import bcrypt from "bcryptjs";
import crypto from "crypto";
import sgMail from '@sendgrid/mail';

import {User} from "../models/User";
import { AuthResponse, JwtPayload, RegisterResponse, ForgotPasswordResponse, ResetPasswordResponse } from "types";


const config = require('config.json')


export async function Authenticate(username: string, password: string): Promise<AuthResponse> {
    return User.scope('auth').findOne({where: {
        username: username.toLowerCase()
    }}).then(user => {
        if (user) {
            if (user.hash.startsWith("old:") && (md5(password.trim()) === user.hash.trim().substr(4))) {
                const token = jwt.sign({ id: user.id, level: user.userlevel, passwordreset: user.passwordreset} as JwtPayload, config.secret, {expiresIn: 24*60*60});
                return { success: true, message: 'Authentication Successful', token: token};
            } else if (bcrypt.compareSync(password, user.hash)) {
                const token = jwt.sign({ id: user.id, level: user.userlevel, passwordreset: user.passwordreset} as JwtPayload, config.secret, {expiresIn: 24*60*60});
                return { success: true, message: 'Authentication Successful', token: token};
            } else {
                console.debug('User failed auth: %s', username);
                return { success: false, message: 'Authentication Failed: User or Password invalid.', token: ''};
            }
            
        } else {
            console.debug('Failed to find user: %s', username);
            return { success: false, message: 'Authentication Failed: User or Password invalid.', token: ''};
        }
    }).catch(err => {
        console.warn(err.message);
        return { success: false, message: 'Authentication Failed: User or Password invalid.', token: ''};
    });
}

export async function Register(username: string, password: string, email: string,
                            firstname: string, lastname: string): Promise<RegisterResponse> {
    const usernameLower = username.toLowerCase();
    const emailLower = email.toLowerCase();
    const failedMessage: RegisterResponse = {success: false, message: "Failed to register user."};
    // Check username and email not used
    const usernameUsed = await usernameTaken(usernameLower).catch(err => {return failedMessage});
    if (usernameUsed) {
        if (typeof usernameUsed === 'boolean') {
            return {success: false, message: "Username Already used!"};
        } else {
            return usernameUsed as unknown as RegisterResponse;
        }
    }
    const emailUsed = await emailTaken(emailLower).catch(err => {return failedMessage});
    if (emailUsed) {
        if (typeof emailUsed === 'boolean') {
            return {success: false, message: "Email Already used!"};
        } else {
        return emailUsed as unknown as RegisterResponse;
        }
    }

    const hash = await hashPassword(password);
    password = "";

    return User.create({username: usernameLower, hash: hash, email: emailLower, firstname: firstname, lastname: lastname, userlevel: 1})
        .then(user=> {
            if (user) {
                return {success: true, message: "User created."}
            } else {
                return {success: false, message: "User NOT created."}
            }
        }).catch(err => {
            console.error(err);
            return {success: false, message: "User NOT created."}
        });
}

export async function PasswordForgot(email: string, hostUrl: string): Promise<ForgotPasswordResponse> {
    return User.findOne({where: {email: email.toLowerCase()}}).then(user => {
        if (user) {
            const token = crypto.randomBytes(20).toString('hex');
            user.passwordResetToken = token;
            user.passwordResetExpires = new Date(Date.now()+3600000);
            return user.save().then(user => {
                //send email
                sgMail.setApiKey(config.sendgridAPIKey);
                var message = {
                    to: user.email,
                    from: 'noreply@jatgam.com',
                    subject: 'Wishlist Password Reset',
                    text: 'You are receiving this because you (or someone else) requested the reset of the password for your account.\n\n'+
                        'Please click the following link, or paste into your browser to complete the process:\n\n'+
                        'https://'+hostUrl+'/password_reset/'+token+'\n\n'+
                        'Username: '+user.username+'\n\n'+
                        'If you did not request this, please ignore this email and your password will remain unchanged.'
                };
                return sgMail.send(message).then(() => {
                    console.log('Sending Password Reset Email to: '+user.email);
                    return {success: true, message: 'Sending an Email to the provided address.'};
                }).catch(error => {
                    console.error(error.toString());
                    return {success: false, message: 'Error Sending Password Reset Email'};
                });
            }).catch(err => {
                console.error(err);
                return {success: false, message: 'Error Trying to Reset Password'};
            });
        } else {
            return {success: false, message: 'No account with that Email.'};
        }
    });
}

export async function PasswordReset(email: string, password: string, token: string): Promise<ResetPasswordResponse> {
    if (await PasswordResetTokenValid(token)) {
        return User.scope('auth').findOne({where: {email: email.toLowerCase(), passwordResetToken: token}})
            .then(async user => {
                if (user) {
                    user.hash = await hashPassword(password);
                    user.passwordreset = false;
                    user.passwordResetToken = null;
                    user.passwordResetExpires = null;

                    return user.save().then(user => {
                        console.debug("Reset password for: %s", user.username);
                        return {success: true, message: 'Password Reset'};
                    }).catch(err => {
                        console.error(err);
                        return {success: false, message: 'Failed to rest the password'};
                    });
                } else {
                    return {success: false, message: 'Password reset token is invalid or expired.'};
                }
            });
    } else {
        return {success: false, message: 'Password reset token is invalid or expired.'};
    }
}

export async function PasswordResetTokenValid(token: string): Promise<boolean> {
    return User.scope('passReset').findOne({where: {passwordResetToken: token}}).then(user => {
        if (user && user.passwordResetToken && user.passwordResetExpires) {
            if (user.passwordResetToken.length < 40) {
                return false;
            }
            if ((token === user.passwordResetToken) && (user.passwordResetExpires.getTime() > new Date(Date.now()).getTime())) {
                return true;
            } else {
                return false;
            }
        } else {
            return false;
        }
    }).catch(err => {
        console.error(err);
        return false;
    });
}

export async function getById(id: number) {
    return await User.findByPk(id);
}

async function usernameTaken(username: string): Promise<boolean> {
    return User.findOne({where: {username: username.toLowerCase()}}).then(user => {
        if (user) {
            return true;
        } else {
            return false;
        }
    }).catch(err => {
        console.error(err.message);
        throw new Error(err)
    })
}

async function emailTaken(email: string): Promise<boolean> {
    return User.findOne({where: {email: email.toLowerCase()}}).then(email => {
        if (email) {
            return true;
        } else {
            return false;
        }
    }).catch(err => {
        console.error(err.message);
        throw new Error(err);
    })
}

async function hashPassword(password: string): Promise<string> {
    return bcrypt.hash(password, 13);
}
