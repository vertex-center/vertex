import styles from "./Instance.module.sass";
import Icon from "../Icon/Icon";
import classNames from "classnames";
import { Horizontal, Vertical } from "../Layouts/Layouts";
import { Link } from "react-router-dom";
import { Fragment, HTMLProps, MouseEventHandler } from "react";
import LoadingValue from "../LoadingValue/LoadingValue";
import { Instance as InstanceModel } from "../../models/instance";
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
    instance: Partial<InstanceModel>;
};

function LCD(props: Readonly<LCDProps>) {
    const { instance } = props;
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
    instance: {
        value: Partial<InstanceModel>;
        to?: string;
        onPower?: () => Promise<void>;
        onInstall?: () => void;
    };
};

export default function Instance(props: Readonly<Props>) {
    const { instance } = props;

    const onPower = (e: any, instance: any) => {
        instance?.onPower();
        e.preventDefault();
    };

    const tags = {
        // Vertex Monitoring
        "vertex-prometheus-collector": "Vertex Monitoring",
        "vertex-grafana-visualizer": "Vertex Monitoring",

        // Vertex Tunnels
        "vertex-cloudflare-tunnel": "Vertex Tunnels",

        // Vertex SQL
        "vertex-postgres-sql": "Vertex SQL",
    };

    const inst = instance.value;
    const tag = inst?.tags?.find((tag) => tag in tags);
    // The uuidv4() is used to generate a unique key for instances that are not yet loaded.
    const key = inst?.uuid ?? uuidv4();

    const content = (
        <Fragment>
            <InstanceLed status={inst?.status} />
            <LCD instance={inst} />

            {tag && (
                <div className={styles.lcdTag}>
                    <LogoIcon />
                    <div>{tags[tag]}</div>
                </div>
            )}

            {instance?.onPower && inst?.status !== "not-installed" && (
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
            {instance?.onInstall && inst?.status === "not-installed" && (
                <Button icon="download" onClick={instance.onInstall} />
            )}
        </Fragment>
    );

    const classnames = classNames({
        [styles.instance]: true,
        [styles.instanceClickable]: instance.to,
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
}

type InstancesProps = HTMLProps<HTMLDivElement> & {
    toolbar?: JSX.Element;
};

export function Instances(props: Readonly<InstancesProps>) {
    const { className, children, toolbar, ...others } = props;
    return (
        <div className={classNames(styles.instances, className)} {...others}>
            {toolbar}
            {children}
        </div>
    );
}
