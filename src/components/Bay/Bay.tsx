import styles from "./Bay.module.sass";
import Icon from "../Icon/Icon";
import classNames from "classnames";
import { Horizontal, Vertical } from "../Layouts/Layouts";
import { Link } from "react-router-dom";
import { Fragment, MouseEventHandler } from "react";
import LoadingValue from "../LoadingValue/LoadingValue";
import { Instance } from "../../models/instance";
import LogoIcon from "../Logo/LogoIcon";
import { InstanceLed } from "../InstanceLed/InstanceLed";
import { v4 as uuidv4 } from "uuid";

type ButtonProps = {
    icon: string;
    onClick: MouseEventHandler<HTMLSpanElement>;
    disabled?: boolean;
};

function Button({ icon, onClick, disabled }: Readonly<ButtonProps>) {
    return (
        <Icon
            disabled={disabled}
            className={classNames({
                [styles.button]: true,
                [styles.buttonDisabled]: disabled,
            })}
            name={icon}
            onClick={onClick}
        />
    );
}

type LCDProps = {
    instance: Partial<Instance>;
    count?: number;
};

function LCD(props: Readonly<LCDProps>) {
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
                    {count !== undefined && (
                        <div className={styles.lcdCount}>{count}</div>
                    )}
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
        onPower?: () => Promise<void>;
        onInstall?: () => void;
    }[];
};

export default function Bay(props: Readonly<Props>) {
    const { instances } = props;

    const onPower = (e: any, instance: any) => {
        instance?.onPower();
        e.preventDefault();
    };

    const tags = {
        "vertex-prometheus-collector": "Vertex Monitoring",
        "vertex-grafana-visualizer": "Vertex Monitoring",
        "vertex-cloudflare-tunnel": "Vertex Tunnels",
    };

    return (
        <div className={styles.group}>
            {instances.map((instance) => {
                const inst = instance.value;
                const tag = inst?.tags?.find((tag) => tag in tags);
                const count = instance.count;
                // The uuidv4() is used to generate a unique key for instances that are not yet loaded.
                const key = inst?.uuid ?? uuidv4();
                const content = (
                    <Fragment>
                        <InstanceLed status={inst?.status} />
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
                                    icon="power_rounded"
                                    onClick={(e: any) => onPower(e, instance)}
                                    disabled={
                                        inst?.status === "building" ||
                                        inst?.status === "starting" ||
                                        inst?.status === "stopping"
                                    }
                                />
                            )}
                        {instance?.onInstall &&
                            inst?.status === "not-installed" && (
                                <Button
                                    icon="download"
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
                        <Link key={key} to={instance.to} className={classnames}>
                            {content}
                        </Link>
                    );

                return (
                    <div key={key} className={classnames}>
                        {content}
                    </div>
                );
            })}
        </div>
    );
}
