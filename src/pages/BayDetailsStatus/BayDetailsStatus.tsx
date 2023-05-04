import { Fragment, useEffect, useState } from "react";
import UptimeGraph from "../../components/UptimeGraph/UptimeGraph";
import { Title } from "../../components/Text/Text";
import { useParams } from "react-router-dom";
import { getInstanceStatus, Uptime } from "../../backend/backend";

type Props = {};

export default function BayDetailsStatus(props: Props) {
    const { uuid } = useParams();

    const [uptimes, setUptimes] = useState<Uptime[]>();

    useEffect(() => {
        getInstanceStatus(uuid).then((uptime) => {
            console.log(uptime);
            setUptimes(uptime);
        });
    }, [uuid]);

    return (
        <Fragment>
            <Title>Status</Title>
            <UptimeGraph uptimes={uptimes} />
        </Fragment>
    );
}
