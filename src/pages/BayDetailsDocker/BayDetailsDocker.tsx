import { Fragment } from "react";
import { Text, Title } from "../../components/Text/Text";

import styles from "./BayDetailsDocker.module.sass";
import { useParams } from "react-router-dom";
import Loading from "../../components/Loading/Loading";
import useInstance from "../../hooks/useInstance";

export default function BayDetailsDocker() {
    const { uuid } = useParams();

    const { instance } = useInstance(uuid);

    let label;

    if (instance?.use_docker) {
        label = <span className={styles.enabled}>enabled</span>;
    } else {
        label = <span className={styles.disabled}>disabled</span>;
    }

    return (
        <Fragment>
            <Title>Docker</Title>
            {!instance && <Loading />}
            {instance && <Text>Docker is {label} for this instance.</Text>}
        </Fragment>
    );
}
