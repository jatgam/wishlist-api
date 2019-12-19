import {createServer} from "http";
import {app} from "./app";
import {sequelize} from "./sequelize";

// start server
const port = 4000;

(async () => {
    await sequelize.sync();

    createServer(app)
        .listen(
            port,
            () => console.info(`Server running on port ${port}`)
        );
})();
