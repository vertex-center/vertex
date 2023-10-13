import styles from "./Container.module.sass";
import Icon from "../Icon/Icon";
import classNames from "classnames";
import { Horizontal, Vertical } from "../Layouts/Layouts";
import { Link } from "react-router-dom";
import { Fragment, HTMLProps, MouseEventHandler } from "react";
import LoadingValue from "../LoadingValue/LoadingValue";
import { Container as ContainerModel } from "../../models/container";
import LogoIcon from "../Logo/LogoIcon";
import { ContainerLed } from "../ContainerLed/ContainerLed";
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
    container: Partial<ContainerModel>;
};

function LCD(props: Readonly<LCDProps>) {
    const { container } = props;
    const { display_name, service, status } = container ?? {};
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
    container: {
        value: Partial<ContainerModel>;
        to?: string;
        onPower?: () => Promise<void>;
        onInstall?: () => void;
    };
};

export default function Container(props: Readonly<Props>) {
    const { container } = props;

    const onPower = (e: any, container: any) => {
        container?.onPower();
        e.preventDefault();
    };

    const tags = [
        "Vertex Internal",
        "Vertex Monitoring",
        "Vertex SQL",
        "Vertex Tunnels",
    ];

    const inst = container.value;
    const tag = tags?.find((t) => inst?.tags?.includes(t));
    // The uuidv4() is used to generate a unique key for containers that are not yet loaded.
    const key = inst?.uuid ?? uuidv4();

    const content = (
        <Fragment>
            <ContainerLed status={inst?.status} />
            <LCD container={inst} />

            {tag && (
                <div className={styles.lcdTag}>
                    <LogoIcon />
                    <div>{tag}</div>
                </div>
            )}

            {container?.onPower && inst?.status !== "not-installed" && (
                <Button
                    icon="power_rounded"
                    onClick={(e: any) => onPower(e, container)}
                    disabled={
                        inst?.status === "building" ||
                        inst?.status === "starting" ||
                        inst?.status === "stopping"
                    }
                />
            )}
            {container?.onInstall && inst?.status === "not-installed" && (
                <Button icon="download" onClick={container.onInstall} />
            )}
        </Fragment>
    );

    const classnames = classNames({
        [styles.container]: true,
        [styles.containerClickable]: container.to,
    });

    if (container.to)
        return (
            <Link key={key} to={container.to} className={classnames}>
                {content}
            </Link>
        );

    return (
        <div key={key} className={classnames}>
            {content}
        </div>
    );
}

type ContainersProps = HTMLProps<HTMLDivElement>;

export function Containers(props: Readonly<ContainersProps>) {
    const { className, children, ...others } = props;
    return (
        <div className={classNames(styles.containers, className)} {...others}>
            {children}
        </div>
    );
}
