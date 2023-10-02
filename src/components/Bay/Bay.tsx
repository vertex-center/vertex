import styles from "./Bay.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import { Horizontal, Vertical } from "../Layouts/Layouts";
import { SiDocker } from "@icons-pack/react-simple-icons";
import { InstallMethod, InstanceUpdate } from "../../models/instance";
import { Link } from "react-router-dom";
import { Fragment, MouseEventHandler } from "react";

type ButtonProps = {
    symbol: string;
    onClick: MouseEventHandler<HTMLSpanElement>;
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
                [styles.ledGreen]: status === "running",
                [styles.ledYellow]: status === "building",
                [styles.ledOrange]:
                    status === "starting" || status === "stopping",
                [styles.ledRed]: status === "error",
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
    count?: number;
    dockerized?: boolean;
};

function LCD(props: LCDProps) {
    const { name, status, count, dockerized } = props;

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
                <Horizontal gap={8}>
                    {name}
                    {count && <div className={styles.lcdCount}>{count}</div>}
                </Horizontal>
                {dockerized && (
                    <SiDocker
                        size={15}
                        style={{ marginTop: -1, opacity: 0.5 }}
                    />
                )}
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
};

export default function Bay(props: Props) {
    const { instances } = props;

    const onPower = (e: any, instance: any) => {
        instance?.onPower();
        e.preventDefault();
    };

    return (
        <div className={styles.group}>
            {instances.map((instance) => {
                const content = (
                    <Fragment>
                        <LED status={instance.status} />
                        <LCD
                            name={instance.name}
                            count={instance.count}
                            status={instance.status}
                            dockerized={instance.method === "docker"}
                        />
                        {instance?.update && (
                            <div>
                                Update available:{" "}
                                {instance.update.current_version} {"->"}{" "}
                                {instance.update.latest_version}
                            </div>
                        )}
                        {instance?.onPower && (
                            <Button
                                symbol="power_rounded"
                                onClick={(e: any) => onPower(e, instance)}
                            />
                        )}
                    </Fragment>
                );

                const classnames = classNames({
                    [styles.bay]: true,
                    [styles.bayClickable]: instance.to,
                });

                if (instance.to)
                    return (
                        <Link to={instance.to} className={classnames}>
                            {content}
                        </Link>
                    );

                return <div className={classnames}>{content}</div>;
            })}
        </div>
    );
}
