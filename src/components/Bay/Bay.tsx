import styles from "./Bay.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import { Vertical } from "../Layouts/Layouts";

type ButtonProps = {
    symbol: string;
};

function Button({ symbol }: ButtonProps) {
    return <Symbol name={symbol} />;
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

type Status = "running" | "error" | "downloading";

type LCDProps = {
    name: string;
    status: Status | string;
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
        case "downloading":
            message = "Downloading...";
            break;
        default:
            message = status;
    }

    return (
        <div className={styles.lcd}>
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
        </div>
    );
}

type Props = {
    name: string;
    status: Status | string;
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
