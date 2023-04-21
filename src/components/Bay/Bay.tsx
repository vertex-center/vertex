import styles from "./Bay.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import { Horizontal, Vertical } from "../Layouts/Layouts";
import { Link } from "react-router-dom";
import Spacer from "../Spacer/Spacer";

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

type Status =
    | "off"
    | "building"
    | "starting"
    | "running"
    | "error"
    | "downloading";

type LCDProps = {
    name: string;
    status: Status | string;
    to?: string;
    count?: number;
};

function LCD(props: LCDProps) {
    const { name, status, to, count } = props;

    let message;
    switch (status) {
        case "off":
            message = "Off";
            break;
        case "building":
            message = "Building...";
            break;
        case "starting":
            message = "Starting...";
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
            <Horizontal>
                <div>{name}</div>
                <Spacer />
                <div className={styles.lcdCount}>{count}</div>
            </Horizontal>
            <div
                className={classNames({
                    [styles.lcdGray]: true,
                    [styles.lcdGreen]: status === "running",
                    [styles.lcdYellow]: status === "building",
                    [styles.lcdOrange]: status === "starting",
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
    instances: {
        name: string;
        status: Status | string;
        count?: number;

        to?: string;

        onPower?: () => void;
    }[];
    showCables?: boolean;
};

export default function Bay(props: Props) {
    const { instances, showCables } = props;

    return (
        <div className={styles.group}>
            {instances.map((instance) => (
                <div
                    className={classNames({
                        [styles.bay]: true,
                        [styles.bayWithCable]:
                            showCables && instance.status !== "off",
                    })}
                >
                    <LED status={instance.status} />
                    <LCD
                        name={instance.name}
                        count={instance.count}
                        status={instance.status}
                        to={instance.to}
                    />
                    {instance?.onPower && (
                        <Button
                            symbol="power_rounded"
                            onClick={instance.onPower}
                        />
                    )}
                </div>
            ))}
        </div>
    );
}
