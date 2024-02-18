import { HTMLProps } from "react";
import cx from "classnames";

export type ListIconProps = HTMLProps<HTMLDivElement>;

export function ListIcon(props: Readonly<ListIconProps>) {
    const { className, ...others } = props;
    return <div className={cx("list-icon", className)} {...others} />;
}
