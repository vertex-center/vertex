import cx from "classnames";
import { HTMLProps } from "react";
import "./Box.sass";
import { MaterialIcon } from "../MaterialIcon/MaterialIcon";

type BoxType = "info" | "tip" | "warning";

type BoxProps = HTMLProps<HTMLDivElement> & {
    type: BoxType;
};

export default function Box(props: Readonly<BoxProps>) {
    const { className, type, children, ...others } = props;

    let label = "";
    switch (type) {
        case "info":
            label = "Info";
            break;
        case "tip":
            label = "Tip";
            break;
        case "warning":
            label = "Warning";
            break;
    }

    return (
        <div className={cx("box", `box-${type}`, className)} {...others}>
            <div className="box-header">
                <MaterialIcon icon="error" />
                <h1>{label}</h1>
            </div>
            <div className="box-content">{children}</div>
        </div>
    );
}
