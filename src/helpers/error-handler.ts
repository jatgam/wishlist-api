


export function errorHandler(err: any, req: any, res: any, next: any) {
    if (typeof (err) === "string") {
        // custom application error
        return res.status(400).json({ message: err });
    }

    if (err.name === "ValidationError") {
        // mongoose validation error
        return res.status(400).json({ message: err.message });
    }

    if (err.name === "UnauthorizedError") {
        // jwt authentication error
        return res.status(401).json({ message: "Invalid Token", path: req.path });
    }

    // default to 500 server error
    return res.status(500).json({ message: err.message });
}

// module.exports = errorHandler;