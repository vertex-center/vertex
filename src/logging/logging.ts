type Line = {
    message: string;
    args: any[];
    tag: string;
    color: string;
    background: string;
};

const display = (line: Line) => {
    const { message, tag, color, background, args } = line;
    console.log(
        `%c${tag}%c ${message}`,
        `color: ${color}; background-color: ${background}; padding: 0 6px;`,
        "",
        ...args
    );
};

const event = (message: string, ...args: any[]) => {
    display({
        message: message,
        args: args,
        tag: "EVENT",
        color: "white",
        background: "#d97e44",
    });
};

const request = (message: string, ...args: any[]) => {
    display({
        message: message,
        args: args,
        tag: "REQUEST",
        color: "white",
        background: "#0e805a",
    });
};

export const Console = {
    event: event,
    request: request,
};
