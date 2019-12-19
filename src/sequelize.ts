import {Sequelize} from "sequelize-typescript";

const env = process.env.NODE_ENV || 'development';
const config = require("./config/config.js")[env];

export const sequelize = new Sequelize({
    username: config.username,
    password: config.password,
    database: config.database,
    host: config.host,
    dialect: config.dialect,
    models: [__dirname + "/models"]
});
