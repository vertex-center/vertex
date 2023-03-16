import { HTMLProps } from "react";

import styles from "./Logs.module.sass";

type LineProps = {
    text: string;
};

function Line(props: LineProps) {
    const { text } = props;

    return <div>{text}</div>;
}

type Props = HTMLProps<HTMLDivElement> & {
    lines: any[];
};

export default function Logs(props: Props) {
    const { lines } = props;

    return (
        <div className={styles.logs}>
            {lines.map((line, i) => (
                <Line key={i} text={line.message} />
            ))}
        </div>
    );
}
