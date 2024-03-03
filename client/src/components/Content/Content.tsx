import { HTMLProps } from "react";
import cx from "classnames";
import { Vertical } from "../Layouts/Layouts";

type PageContentProps = HTMLProps<HTMLDivElement>;

export default function Content(props: Readonly<PageContentProps>) {
    const { className, ...others } = props;
    return (
        <Vertical gap={24} className={cx("content", className)} {...others} />
    );
}
