import cx from "classnames";
import { HTMLProps } from "react";
import "./Box.sass";
import { MaterialIcon } from "../MaterialIcon/MaterialIcon";

export type BoxType = "info" | "tip" | "warning" | "error";

export type BoxProps = HTMLProps<HTMLDivElement> & {
    type: BoxType;
};

export function Box(props: Readonly<BoxProps>) {
    const { className, type, children, ...others } = props;

    let label = "",
        icon = "";

    switch (type) {
        case "info":
            label = "Info";
            icon = "info";
            break;
        case "tip":
            label = "Tip";
            icon = "lightbulb";
            break;
        case "warning":
            label = "Warning";
            icon = "warning";
            break;
        case "error":
            label = "Error";
            icon = "error";
            break;
    }

    return (
        <div className={cx("box", `box-${type}`, className)} {...others}>
            <div className="box-header">
                <MaterialIcon icon={icon} />
                {label}
            </div>
            <div className="box-content">{children}</div>
        </div>
    );
}
