import { HTMLProps, useEffect, useRef, useState } from "react";

import styles from "./Logs.module.sass";
import classNames from "classnames";
import useScrollPercentage from "react-scroll-percentage-hook";
import Downloads from "../Downloads/Downloads";

export type LogLine = {
    kind: "out" | "err" | "downloads";
    message: any;
};

function Line(props: LogLine) {
    const { kind, message } = props;

    if (kind === "downloads") {
        return <Downloads downloads={message} />;
    }

    return (
        <div
            className={classNames({
                [styles.line]: true,
                [styles.lineError]: kind === "err",
            })}
        >
            <div>{message}</div>
        </div>
    );
}

type Props = HTMLProps<HTMLDivElement> & {
    lines: LogLine[];
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
            {lines.map((line, i) => {
                return <Line key={i} kind={line.kind} message={line.message} />;
            })}
            <div ref={scroll} />
        </div>
    );
}
