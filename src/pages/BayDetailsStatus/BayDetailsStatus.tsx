import { Fragment, useEffect, useState } from "react";
import UptimeGraphs from "../../components/UptimeGraph/UptimeGraph";
import { Title } from "../../components/Text/Text";
import { useParams } from "react-router-dom";
import { getInstanceStatus, route, Uptime } from "../../backend/backend";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";

type Props = {};

export default function BayDetailsStatus(props: Props) {
    const { uuid } = useParams();

    const [uptimes, setUptimes] = useState<Uptime[]>();

    const reload = () => {
        getInstanceStatus(uuid).then((uptime) => {
            console.log(uptime);
            setUptimes(uptime);
        });
    };

    useEffect(() => {
        if (uuid === undefined) return;

        reload();

        const sse = registerSSE(route(`/instance/${uuid}/events`));

        const onStatusChange = (e) => {
            reload();
        };

        registerSSEEvent(sse, "uptime_status_change", onStatusChange);

        return () => {
            unregisterSSEEvent(sse, "uptime_status_change", onStatusChange);
            unregisterSSE(sse);
        };
    }, [uuid]);

    return (
        <Fragment>
            <Title>Status</Title>
            <UptimeGraphs uptimes={uptimes} />
        </Fragment>
    );
}
