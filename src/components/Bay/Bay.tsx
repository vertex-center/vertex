import styles from "./Bay.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import { Horizontal, Vertical } from "../Layouts/Layouts";
import { Link } from "react-router-dom";
import Spacer from "../Spacer/Spacer";
import { SiDocker } from "@icons-pack/react-simple-icons";
import { InstallMethod, InstanceUpdate } from "../../models/instance";

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
    | "stopping"
    | "error"
    | "downloading";

type LCDProps = {
    name: string;
    status: Status | string;
    to?: string;
    count?: number;
    dockerized?: boolean;
};

function LCD(props: LCDProps) {
    const { name, status, to, count, dockerized } = props;

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
        case "stopping":
            message = "Stopping...";
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
            <Horizontal gap={8}>
                <div>{name}</div>
                {dockerized && (
                    <SiDocker
                        size={15}
                        style={{ marginTop: -1, opacity: 0.5 }}
                    />
                )}
                <Spacer />
                <div className={styles.lcdCount}>{count}</div>
            </Horizontal>
            <div
                className={classNames({
                    [styles.lcdGray]: true,
                    [styles.lcdGreen]: status === "running",
                    [styles.lcdYellow]: status === "building",
                    [styles.lcdOrange]:
                        status === "starting" || status === "stopping",
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
        method?: InstallMethod;
        update?: InstanceUpdate;

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
                        dockerized={instance.method === "docker"}
                    />
                    {instance?.update && (
                        <div>
                            Update available: {instance.update.current_version}{" "}
                            {"->"} {instance.update.latest_version}
                        </div>
                    )}
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
