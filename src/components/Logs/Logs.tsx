import { HTMLProps, useEffect, useRef, useState } from "react";

import styles from "./Logs.module.sass";
import classNames from "classnames";
import useScrollPercentage from "react-scroll-percentage-hook";

type LineProps = {
    type: string;
    text: string;
};

function Line(props: LineProps) {
    const { type, text } = props;

    if (text === "") return null;

    return (
        <div
            className={classNames({
                [styles.line]: true,
                [styles.lineError]: type === "error",
            })}
        >
            {text}
        </div>
    );
}

type Props = HTMLProps<HTMLDivElement> & {
    lines: any[];
};

export default function Logs(props: Props) {
    const { lines } = props;

    const { ref } = useScrollPercentage<HTMLDivElement>({
        onProgress: (percentage) => {
            setAutoScroll(percentage.vertical === 100);
        },
    });

    const [autoScroll, setAutoScroll] = useState<boolean>(true);

    const scroll = useRef();

    useEffect(() => {
        if (!autoScroll) return;
        let s: any = scroll;
        s.current.scrollIntoView({ behavior: "smooth" });
    }, [autoScroll, lines]);

    return (
        <div className={styles.logs} ref={ref}>
            {lines.map((line, i) => (
                <Line key={i} type={line.type} text={line.message} />
            ))}
            <div ref={scroll} />
        </div>
    );
}
