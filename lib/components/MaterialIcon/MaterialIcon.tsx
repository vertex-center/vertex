import { HTMLProps } from "react";

type Props = HTMLProps<HTMLSpanElement>;

export function MaterialIcon(props: Readonly<Props>) {
    return <span className="material-symbols-rounded">{props.name}</span>;
}
