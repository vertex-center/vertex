import { Horizontal, Vertical } from "../Layouts/Layouts";

import styles from "./UptimeGraph.module.sass";
import classNames from "classnames";
import { Text } from "../Text/Text";
import { Uptime } from "../../backend/backend";

type Props = {
    uptimes?: Uptime[];
};

export default function UptimeGraph(props: Props) {
    const { uptimes } = props;

    return (
        <Vertical gap={12} alignItems="flex-start">
            {uptimes?.map((uptime) => (
                <Vertical gap={16}>
                    <Text>{uptime.name}</Text>
                    <Horizontal gap={6} className={styles.graph}>
                        {uptime.history.map((point, i) => (
                            <div
                                key={i}
                                className={classNames({
                                    [styles.bar]: true,
                                    [styles.barOff]: point.status === "off",
                                    [styles.barOn]: point.status === "on",
                                })}
                            />
                        ))}
                        <div
                            className={classNames({
                                [styles.bar]: true,
                                [styles.barCurrent]: true,
                                [styles.barOff]: uptime.current === "off",
                                [styles.barOn]: uptime.current === "on",
                            })}
                        />
                    </Horizontal>
                </Vertical>
            ))}
        </Vertical>
    );
}
