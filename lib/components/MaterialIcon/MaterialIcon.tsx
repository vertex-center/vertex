import { HTMLProps } from "react";
import classNames from "classnames";

type Props = HTMLProps<HTMLSpanElement>;

export function MaterialIcon(props: Readonly<Props>) {
    const { className, ...others } = props;
    return (
        <span
            className={classNames("material-symbols-rounded", className)}
            {...others}
        >
            {props.name}
        </span>
    );
}
