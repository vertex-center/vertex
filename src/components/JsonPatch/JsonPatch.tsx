import { Operation } from "../../models/instance";

import styles from "./JsonPatch.module.sass";
import classNames from "classnames";
import Symbol from "../Symbol/Symbol";

type LineProps = {
    operation: Operation;
};

function Line(props: LineProps) {
    const { operation } = props;

    let { op, from, path, value } = operation;

    path = path.replace("/", "");
    path = path.replaceAll("/", " → ");
    path = path.replaceAll("~1", "/");

    let symbol = "";
    switch (op) {
        case "add":
            symbol = "add";
            break;
        case "remove":
            symbol = "remove";
            break;
        case "replace":
            symbol = "update";
            break;
        case "copy":
            symbol = "content_copy";
            break;
        case "move":
            symbol = "drag_indicator";
            break;
    }

    return (
        <div className={styles.line}>
            <div
                className={classNames({
                    [styles.operation]: true,
                })}
            >
                <Symbol name={symbol} />
            </div>
            {from && <div className={styles.path}>{from}</div>}
            <div className={styles.path}>{path}</div>
            {value && <Symbol className={styles.equal} name="equal" />}
            {value && <div>{value}</div>}
        </div>
    );
}

type Props = {
    operations: Operation[];
};

export default function JsonPatch(props: Props) {
    const { operations } = props;
    return (
        <div className={styles.patch}>
            {operations?.map((operation: Operation, i: number) => (
                <Line key={i} operation={operation} />
            ))}
        </div>
    );
}
