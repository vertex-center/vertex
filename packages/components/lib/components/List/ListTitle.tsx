import { HTMLProps } from "react";

export type ListTitleProps = HTMLProps<HTMLDivElement>;

export function ListTitle(props: Readonly<ListTitleProps>) {
    return <div {...props} />;
}
