import { Horizontal, Vertical } from "../Layouts/Layouts";

import styles from "./UptimeGraph.module.sass";
import classNames from "classnames";
import { Text } from "../Text/Text";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import { Uptime } from "../../models/uptime";

function UptimeGraph({ uptime }: { uptime: Uptime }) {
    const [tick, setTick] = useState<boolean>();

    const { uuid } = useParams();

    useEffect(() => {
        if (uuid === undefined) return;

        const sse = registerSSE(`/instance/${uuid}/events`);

        const onStatusChange = (e) => {
            if (e.data == "on") {
                setTick(true);
                setTimeout(() => setTick(false), 200);
            }
        };

        registerSSEEvent(sse, "uptime_status_change", onStatusChange);

        return () => {
            unregisterSSEEvent(sse, "uptime_status_change", onStatusChange);
            unregisterSSE(sse);
        };
    }, [uuid]);

    return (
        <Vertical gap={16}>
            <Horizontal gap={12} alignItems="center">
                <Text>{uptime.name}</Text>
                <div
                    className={classNames({
                        [styles.livedot]: true,
                        [styles.livedotTick]: tick,
                    })}
                />
            </Horizontal>
            <Horizontal gap={6} className={styles.graph} alignItems="flex-end">
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
    );
}

type Props = {
    uptimes?: Uptime[];
};

export default function UptimeGraphs(props: Props) {
    const { uptimes } = props;

    return (
        <Vertical gap={12} alignItems="flex-start">
            {uptimes?.map((uptime) => (
                <UptimeGraph key={uptime.name} uptime={uptime} />
            ))}
        </Vertical>
    );
}
