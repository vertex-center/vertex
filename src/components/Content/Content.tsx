import { HTMLProps } from "react";
import "./Content.sass";
import cx from "classnames";

type PageContentProps = HTMLProps<HTMLDivElement>;

export default function Content(props: Readonly<PageContentProps>) {
    const { className, ...others } = props;
    return <div className={cx("content", className)} {...others} />;
}
