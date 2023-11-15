import { HTMLProps } from "react";
import cx from "classnames";
import "./Paragraph.sass";

export type ParagraphProps = HTMLProps<HTMLParagraphElement>;

export function Paragraph(props: Readonly<ParagraphProps>) {
    const { className, ...others } = props;
    return <p className={cx("paragraph", className)} {...others} />;
}
