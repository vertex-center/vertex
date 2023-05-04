import { Fragment } from "react";
import UptimeGraph from "../../components/UptimeGraph/UptimeGraph";
import { Title } from "../../components/Text/Text";

type Props = {};

export default function BayDetailsStatus(props: Props) {
    return (
        <Fragment>
            <Title>Status</Title>
            <UptimeGraph title="Service" />
        </Fragment>
    );
}
