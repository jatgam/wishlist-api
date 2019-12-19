require("rootpath")();
import express from "express";

import cors from "cors";
import bodyParser from "body-parser";

import {jwt} from "./helpers/jwt";
import {errorHandler} from "./helpers/error-handler";

import {users} from "./routes/users";
import {items} from "./routes/items";

export const app = express();

app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.use(cors());

// use JWT auth to secure the api
app.use(jwt());

app.use("/user", users);
app.use("/item", items);

// global error handler
app.use(errorHandler);

app.disable("x-powered-by")
