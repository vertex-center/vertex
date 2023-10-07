import styles from "./Bay.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import { Horizontal, Vertical } from "../Layouts/Layouts";
import { Link } from "react-router-dom";
import { Fragment, MouseEventHandler } from "react";
import LoadingValue from "../LoadingValue/LoadingValue";
import { Instance } from "../../models/instance";
import LogoIcon from "../Logo/LogoIcon";

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
                [styles.ledRed]: status === "error" || status === "off",
                [styles.ledGreen]: status === "running",
                [styles.ledYellow]: status === "building",
                [styles.ledOrange]:
                    status === "starting" || status === "stopping",
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
    | "downloading"
    | "not-installed";

type LCDProps = {
    instance: Partial<Instance>;
    count?: number;
};

function LCD(props: LCDProps) {
    const { instance, count } = props;
    const { display_name, service, status } = instance ?? {};
    const { name } = service ?? {};

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
        case "not-installed":
            message = "Not installed";
            break;
        default:
            message = status;
    }

    let content = (
        <Vertical gap={10}>
            <Horizontal gap={8}>
                <Horizontal gap={8}>
                    {display_name ?? name ?? <LoadingValue />}
                    {count && <div className={styles.lcdCount}>{count}</div>}
                </Horizontal>
            </Horizontal>
            <div
                className={classNames({
                    [styles.lcdGray]: true,
                    [styles.lcdRed]: status === "error",
                    [styles.lcdGreen]: status === "running",
                    [styles.lcdYellow]: status === "building",
                    [styles.lcdOrange]:
                        status === "starting" || status === "stopping",
                    [styles.lcdDownloading]: status === "downloading",
                })}
            >
                {message ?? <LoadingValue />}
            </div>
        </Vertical>
    );

    return <div className={styles.lcd}>{content}</div>;
}

type Props = {
    instances: {
        value: Partial<Instance>;
        count?: number;
        to?: string;
        onPower?: () => void;
        onInstall?: () => void;
    }[];
};

export default function Bay(props: Props) {
    const { instances } = props;

    const onPower = (e: any, instance: any) => {
        instance?.onPower();
        e.preventDefault();
    };

    const tags = {
        "vertex-prometheus-collector": "Vertex Monitoring",
        "vertex-grafana-visualizer": "Vertex Monitoring",
    };

    return (
        <div className={styles.group}>
            {instances.map((instance) => {
                const inst = instance.value;
                const tag = inst?.tags?.find((tag) => tag in tags);
                const count = instance.count;
                const content = (
                    <Fragment>
                        <LED status={inst?.status} />
                        <LCD instance={inst} count={count} />

                        {tag && (
                            <div className={styles.lcdTag}>
                                <LogoIcon />
                                <div>{tags[tag]}</div>
                            </div>
                        )}

                        {instance?.onPower &&
                            inst?.status !== "not-installed" && (
                                <Button
                                    symbol="power_rounded"
                                    onClick={(e: any) => onPower(e, instance)}
                                />
                            )}
                        {instance?.onInstall &&
                            inst?.status === "not-installed" && (
                                <Button
                                    symbol="download"
                                    onClick={instance.onInstall}
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
