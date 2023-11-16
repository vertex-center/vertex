import { HTMLProps } from "react";
import cx from "classnames";
import "./List.sass";

export type ListProps = HTMLProps<HTMLDivElement>;

export function List(props: Readonly<ListProps>) {
    const { className, ...others } = props;
    return <div className={cx("list", className)} {...others} />;
}
