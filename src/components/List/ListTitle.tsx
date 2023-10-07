import { HTMLProps } from "react";

export type ListTitleProps = HTMLProps<HTMLDivElement>;

export default function ListTitle(props: Readonly<ListTitleProps>) {
    return <div {...props} />;
}
