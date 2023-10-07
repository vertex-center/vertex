import classNames from "classnames";
import styles from "./InstanceLed.module.sass";

export type Status =
    | "off"
    | "building"
    | "starting"
    | "running"
    | "stopping"
    | "error"
    | "downloading"
    | "not-installed";

type LEDProps = {
    small?: boolean;
    status: Status | string;
};

export function InstanceLed(props: Readonly<LEDProps>) {
    const { small, status } = props;
    return (
        <div
            className={classNames({
                [styles.led]: true,
                [styles.ledSmall]: small,
                [styles.ledRed]: status === "error" || status === "off",
                [styles.ledGreen]: status === "running",
                [styles.ledYellow]: status === "building",
                [styles.ledOrange]:
                    status === "starting" || status === "stopping",
            })}
        ></div>
    );
}
