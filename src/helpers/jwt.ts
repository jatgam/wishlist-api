import expressJwt from 'express-jwt';
const config = require('config.json');
import * as userService from '../service/userService';

export function jwt() {
    const secret = config.secret;
    return expressJwt({ secret, isRevoked, getToken, requestProperty: 'auth' }).unless({
        path: [
            // public routes that don't require authentication
            '/user/auth',
            '/user/register',
            /^\/user\/password_reset\/.*/,
            '/user/password_forgot',
            {url: '/item', methods: ['GET']}
        ]
    });
}

function getToken(req: any) {
    if (req.headers['x-access-token']) {
        return req.headers['x-access-token'];
    }
    return null;
};

async function isRevoked(req: any, payload: any, done: any) {
    const user = await userService.getById(payload.id);

    // revoke token if user no longer exists
    if (!user) {
        return done(null, true);
    }

    done();
};
