import "./InlineCode.sass";
import { HTMLProps } from "react";
import cx from "classnames";

export type InlineCodeProps = HTMLProps<HTMLDivElement>;

export function InlineCode(props: Readonly<InlineCodeProps>) {
    const { children, className, ...others } = props;

    return (
        <code className={cx("inline-code", className)} {...others}>
            {children}
        </code>
    );
}
