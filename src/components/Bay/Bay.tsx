import styles from "./Bay.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";

type ButtonProps = {
    symbol: string;
};

function Button({ symbol }: ButtonProps) {
    return <Symbol name={symbol} />;
}

type LEDProps = {
    status: Status;
};

function LED({ status }: LEDProps) {
    return (
        <div
            className={classNames({
                [styles.led]: true,
                [styles.ledGreen]: status === "running",
                [styles.ledRed]: status === "error",
            })}
        ></div>
    );
}

type Status = "running" | "error";

type LCDProps = {
    name: string;
    status: Status;
};

function LCD(props: LCDProps) {
    const { name, status } = props;

    let message;
    switch (status) {
        case "running":
            message = "Running";
            break;
        case "error":
            message = "Fatal error";
            break;
    }

    return (
        <div className={styles.lcd}>
            <div>{name}</div>
            <div
                className={classNames({
                    [styles.lcdGreen]: status === "running",
                    [styles.lcdRed]: status === "error",
                })}
            >
                {message}
            </div>
        </div>
    );
}

type Props = {
    name: string;
    status: Status;
};

export default function Bay(props: Props) {
    const { name, status } = props;
    return (
        <div className={styles.bay}>
            <LED status={status} />
            <LCD name={name} status={status} />
            <Button symbol="power_rounded" />
        </div>
    );
}
