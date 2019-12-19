import {Request, Response} from "express";
import {Item} from "../models/Item";

declare global {
    namespace Express {
        interface Request {
            auth: ExpressJwtPayloadResponse;
        }
    }
}



export type AuthResponse = {
    success: boolean;
    message: string;
    token: string;
}

export type RegisterResponse = {
    success: boolean;
    message: string;
}

export type ForgotPasswordResponse = {
    success: boolean;
    message: string;
}

export type ResetPasswordResponse = {
    success: boolean;
    message: string;
}

export type JwtPayload = {
    id: number;
    level: number;
    passwordreset: boolean;
}

export interface ExpressJwtPayloadResponse extends JwtPayload {
    iat: number;
    exp: number;
}

export type GetItemsResponse = {
    success: boolean;
    items: Item[];
}

export type GenericResponse = {
    success: boolean;
    message: string;
}

export type Nullable<T> = T | null;
