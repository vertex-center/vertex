import cx from "classnames";
import { HTMLProps } from "react";
import "./Box.sass";
import { Info, Lightbulb, Warning, WarningCircle } from "@phosphor-icons/react";

export type BoxType = "info" | "tip" | "warning" | "error";

export type BoxProps = HTMLProps<HTMLDivElement> & {
    type: BoxType;
};

export function Box(props: Readonly<BoxProps>) {
    const { className, type, children, ...others } = props;

    let label = "",
        icon = null;

    switch (type) {
        case "info":
            label = "Info";
            icon = <Info />;
            break;
        case "tip":
            label = "Tip";
            icon = <Lightbulb />;
            break;
        case "warning":
            label = "Warning";
            icon = <Warning />;
            break;
        case "error":
            label = "Error";
            icon = <WarningCircle />;
            break;
    }

    return (
        <div className={cx("box", `box-${type}`, className)} {...others}>
            <div className="box-header">
                {icon}
                {label}
            </div>
            <div className="box-content">{children}</div>
        </div>
    );
}
