import { HTMLProps } from "react";

export type ListTitleProps = HTMLProps<HTMLDivElement>;

export default function ListTitle(props: ListTitleProps) {
    return <div {...props} />;
}
