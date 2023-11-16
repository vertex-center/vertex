import { HTMLProps } from "react";
import cx from "classnames";

export type ListInfoProps = HTMLProps<HTMLDivElement>;

export function ListInfo(props: Readonly<ListInfoProps>) {
    const { className, ...others } = props;
    return <div className={cx("list-info", className)} {...others} />;
}
