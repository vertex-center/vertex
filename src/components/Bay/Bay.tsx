import styles from "./Bay.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import { Vertical } from "../Layouts/Layouts";
import { Link } from "react-router-dom";

type ButtonProps = {
    symbol: string;
    onClick: () => void;
};

function Button({ symbol, onClick }: ButtonProps) {
    return <Symbol className={styles.button} name={symbol} onClick={onClick} />;
}

type LEDProps = {
    status: Status | string;
};

function LED({ status }: LEDProps) {
    return (
        <div
            className={classNames({
                [styles.led]: true,
                [styles.ledRunning]: status === "running",
                [styles.ledError]: status === "error",
                [styles.ledDownloading]: status === "downloading",
            })}
        ></div>
    );
}

type Status = "off" | "running" | "error" | "downloading";

type LCDProps = {
    name: string;
    status: Status | string;
    to?: string;
};

function LCD(props: LCDProps) {
    const { name, status, to } = props;

    let message;
    switch (status) {
        case "off":
            message = "Off";
            break;
        case "running":
            message = "Running";
            break;
        case "error":
            message = "Fatal error";
            break;
        case "downloading":
            message = "Downloading...";
            break;
        default:
            message = status;
    }

    let content = (
        <Vertical gap={10}>
            <div>{name}</div>
            <div
                className={classNames({
                    [styles.lcdGray]: true,
                    [styles.lcdGreen]: status === "running",
                    [styles.lcdRed]: status === "error",
                    [styles.lcdDownloading]: status === "downloading",
                })}
            >
                {message}
            </div>
        </Vertical>
    );

    if (to)
        return (
            <Link
                to={to}
                className={classNames({
                    [styles.lcd]: true,
                    [styles.lcdClickable]: to,
                })}
            >
                {content}
            </Link>
        );

    return <div className={styles.lcd}>{content}</div>;
}

type Props = {
    name: string;
    status: Status | string;

    to?: string;

    onPower?: () => void;
};

export default function Bay(props: Props) {
    const { name, status, to, onPower } = props;

    return (
        <div className={styles.bay}>
            <LED status={status} />
            <LCD name={name} status={status} to={to} />
            {onPower && <Button symbol="power_rounded" onClick={onPower} />}
        </div>
    );
}
